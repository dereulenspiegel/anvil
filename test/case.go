package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/dereulenspiegel/anvil/config"
	"github.com/dereulenspiegel/anvil/plugin"
	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/dereulenspiegel/anvil/util"
	"github.com/ryanfaerman/fsm"
	"github.com/ttacon/chalk"
)

const (
	DESTROYED   fsm.State = "DESTROYED"
	SETUP       fsm.State = "SETUP"
	PROVISIONED fsm.State = "PROVISIONED"
	VERIFIED    fsm.State = "VERIFIED"
	FAILED      fsm.State = "FAILED"
)

var (
	orderedStateList = []fsm.State{DESTROYED, SETUP, PROVISIONED, VERIFIED}
)

func stateIndex(state fsm.State) int {
	for i, s := range orderedStateList {
		if s == state {
			return i
		}
	}
	return -1
}

type TestCase struct {
	State       fsm.State
	platform    *config.PlatformConfig
	suite       *config.SuiteConfig
	Name        string
	machine     fsm.Machine
	driver      *plugin.DriverPlugin
	provisioner *plugin.ProvisionerPlugin
	Instance    apis.Instance
	lastError   error
}

func CompileTestCasesFromConfig(cfg *config.Config) []*TestCase {
	testCases := make([]*TestCase, 0, len(cfg.Platforms)*len(cfg.Suites))
	for _, platform := range cfg.Platforms {
		for _, suite := range cfg.Suites {
			testCases = append(testCases, NewTestCase(platform, suite))
		}
	}
	return testCases
}

func NewTestCase(platform *config.PlatformConfig, suite *config.SuiteConfig) *TestCase {
	testCase := &TestCase{
		platform: platform,
		suite:    suite,
		State:    DESTROYED,
		Name:     fmt.Sprintf("anvil-%s-%s", platform.Name, suite.Name),
	}
	machine := fsm.New(fsm.WithRules(createTestCaseRules()), fsm.WithSubject(testCase))
	testCase.machine = machine
	testCase.LoadState()
	testCase.driver = plugin.LoadDriver(config.Cfg.Driver.Name)
	testCase.provisioner = plugin.LoadProvisioner(config.Cfg.Provisioner.Name)

	return testCase
}

func (t *TestCase) UpdateState() {
	instance, err := t.driver.UpdateState(t.Name)
	if err == nil {
		t.Instance = instance
		if instance.State == apis.STARTED && t.State == DESTROYED {
			t.State = SETUP
		}
		if instance.State == apis.DESTROYED && t.State != DESTROYED {
			t.State = DESTROYED
		}
	} else {
		log.Fatalf("Got error when udpating state: %v", err)
	}
}

func (t *TestCase) Transition(s fsm.State) error {
	// Allow destruction in all states
	t.lastError = nil
	if s == DESTROYED && t.State != DESTROYED {
		err := t.machine.Transition(DESTROYED)
		if err != nil {
			return err
		}
		if t.lastError != nil {
			return t.lastError
		}
		return nil
	}
	/*// Allow reprovisioning
	if s == PROVISIONED && t.State == PROVISIONED {
		err := t.machine.Transition(PROVISIONED)
		if err != nil {
			return err
		}
		if t.lastError != nil {
			return t.lastError
		}
	}

	if s == VERIFIED && t.State == VERIFIED {
		err := t.machine.Transition(VERIFIED)
		if err != nil {
			return err
		}
		if t.lastError != nil {
			return t.lastError
		}
	}*/
	currStateIndex := stateIndex(t.State)
	nextStateIndex := stateIndex(s)
	if t.State != FAILED && (currStateIndex == -1 || nextStateIndex == -1) {
		log.Fatalf("Either %s or %s are unknown states", t.State, s)
	}
	if t.State == FAILED {
		currStateIndex = nextStateIndex
	}
	if currStateIndex == nextStateIndex {
		t.lastError = nil
		err := t.machine.Transition(s)
		if err != nil {
			return err
		}
		if t.lastError != nil {
			return t.lastError
		}
	}
	for i := currStateIndex + 1; i <= nextStateIndex; i++ {
		t.lastError = nil
		err := t.machine.Transition(orderedStateList[i])
		if err != nil {
			return err
		}
		if t.lastError != nil {
			return t.lastError
		}
	}
	/*var err error
	err = t.machine.Transition(s)
	*/
	return nil
}

func (t *TestCase) CurrentState() fsm.State {
	return t.State
}

func (t *TestCase) SetState(s fsm.State) {
	switch s {
	case SETUP:
		t.setup()
	case PROVISIONED:
		t.provision()
	case VERIFIED:
		t.verify()
	case DESTROYED:
		t.destroy()
	}
}

func (t *TestCase) setup() {
	fmt.Printf("%sSetting up %s%s\n", chalk.Bold, t.Name, chalk.Reset)
	instance, err := t.driver.CreateInstance(t.Name, t.platform.Driver)
	if err != nil {
		t.lastError = err
		t.State = FAILED
		return
	}
	t.Instance = instance
	instance, err = t.driver.StartInstance(instance.Name)
	if err != nil {
		t.State = FAILED
		t.lastError = err
		return
	}
	t.Instance = instance
	t.State = SETUP
}

func (t *TestCase) provision() {
	fmt.Printf("%sProvisioning %s%s\n", chalk.Bold, t.Name, chalk.Reset)
	err := t.provisioner.Provision(t.Instance, t.suite.Provisioner)
	if err != nil {
		t.State = FAILED
		t.lastError = err
		return
	}
	t.State = PROVISIONED
}

func (t *TestCase) verify() {
	testPath := path.Join(apis.DefaultTestFolder, t.suite.Name)
	if !util.FileExists(testPath) {
		// TODO check if this makes sense. Probably state PROVSIONED is better
		t.State = VERIFIED
		return
	}
	files, err := ioutil.ReadDir(testPath)
	if err != nil {
		t.State = FAILED
		t.lastError = err
		return
	}
	fmt.Printf("%sVerifying %s%s\n", chalk.Bold, t.Name, chalk.Reset)
	for _, file := range files {
		if file.IsDir() {
			result, err := t.verifyUsingVerifier(file.Name())
			failed := false
			for _, r := range result.CaseResults {
				if !r.Success || r.ErrorMsg != "" {
					failed = true
				}
			}
			if err != nil || failed {
				t.State = FAILED
				t.lastError = err
			} else {
				t.State = VERIFIED
			}
			t.printVerifyResult(result)
		}
	}
}

func (t *TestCase) loadFormatter() apis.VerifyResultFormatter {
	if config.Cfg.Formatter != nil {
		return plugin.LoadVerifyResultFormatter(config.Cfg.Formatter.Name)
	}
	return &DefaultConsoleResultFormatter{}
}

func (t *TestCase) printVerifyResult(result apis.VerifyResult) {
	formatter := t.loadFormatter()
	data, err := formatter.Format(result)
	if err != nil {
		log.Fatalf("Can't format verify results: %v", err)
		return
	}
	os.Stdout.Write(data)
}

func (t *TestCase) verifyUsingVerifier(name string) (apis.VerifyResult, error) {
	verifier := plugin.LoadVerifier(name)
	return verifier.Verify(t.Instance, t.suite)
}

func (t *TestCase) destroy() {
	fmt.Printf("%sDestroying %s%s\n", chalk.Bold, t.Name, chalk.Reset)
	instance, err := t.driver.DestroyInstance(t.Instance.Name)
	if err != nil {
		t.lastError = err
		// TODO Fail or continue with other instances?
		log.Fatalf("Can't destroy instance %s: %v", t.Instance.Name, err)
	}
	t.Instance = instance
	t.State = DESTROYED
}
