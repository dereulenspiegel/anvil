package test

import (
	"fmt"
	"log"

	"github.com/dereulenspiegel/anvil/config"
	"github.com/dereulenspiegel/anvil/plugin"
	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/ryanfaerman/fsm"
)

const (
	DESTROYED   fsm.State = "DESTROYED"
	SETUP       fsm.State = "SETUP"
	PROVISIONED fsm.State = "PROVISIONED"
	VERIFIED    fsm.State = "VERIFIED"
	FAILED      fsm.State = "FAILED"
)

type TestCase struct {
	State    fsm.State
	platform *config.PlatformConfig
	suite    *config.SuiteConfig
	Name     string
	machine  fsm.Machine
	driver   *plugin.DriverPlugin
	Instance apis.Instance
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
	instance, err := testCase.driver.UpdateState(testCase.Name)
	if err == nil {
		testCase.Instance = instance
		if instance.State == apis.STARTED {
			testCase.State = SETUP
		}
	} else {
		log.Fatalf("Got error when udpating state: %v", err)
	}
	return testCase
}

func (t *TestCase) Transition(s fsm.State) error {
	log.Printf("Transitioning from state %s to state %s", t.State, s)
	var err error
	err = t.machine.Transition(s)
	return err
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
	log.Printf("Creating instance....")
	instance, err := t.driver.CreateInstance(t.Name, t.platform.Driver)
	if err != nil {
		log.Printf("[ERROR]: Creating instance failed: %v", err)
		return
	}
	t.Instance = instance
	log.Printf("Starting instance...")
	instance, err = t.driver.StartInstance(instance.Name)
	if err != nil {
		log.Printf("[ERROR]: Starting instance failed: %v", err)
		return
	}
	t.Instance = instance
	t.State = SETUP
}

func (t *TestCase) provision() {
	t.State = PROVISIONED
}

func (t *TestCase) verify() {
	t.State = VERIFIED
}

func (t *TestCase) destroy() {
	log.Printf("Destroying instance %s", t.Instance.Name)
	instance, err := t.driver.DestroyInstance(t.Instance.Name)
	if err != nil {
		log.Panicf("Can't destroy instance %s", t.Instance.Name)
	}
	t.Instance = instance
	t.State = DESTROYED
}
