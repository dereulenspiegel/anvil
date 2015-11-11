package state

import (
	"log"
	"time"

	"io/ioutil"
	"path"

	"github.com/dereulenspiegel/anvil/test"
	"github.com/dereulenspiegel/anvil/util"
	"gopkg.in/yaml.v2"
)

var (
	DefaultStateFolder   = ".anvil"
	DefaultStateFileName = "state.yml"
)

type GlobalState struct {
	TestCases []*test.TestCase
	Timestamp int64
}

func PersistGlobalState(testCases []*test.TestCase) {
	state := GlobalState{
		TestCases: testCases,
		Timestamp: time.Now().Unix(),
	}
	out, err := yaml.Marshal(state)
	if err != nil {
		log.Fatalf("Can't marshal state: %v", err)
	}
	err = util.CreateDirectoryIfNotExists(DefaultStateFolder)
	if err != nil {
		log.Panicf("Can't create state directory %s: %v", DefaultStateFolder, err)
	}
	statePath := path.Join(DefaultStateFolder, DefaultStateFileName)
	err = ioutil.WriteFile(statePath, out, 0655)
	if err != nil {
		log.Fatalf("Can' write state to %s: %v", statePath, err)
	}
}

func LoadGlobalState() *GlobalState {
	state := &GlobalState{}
	statePath := path.Join(DefaultStateFolder, DefaultStateFileName)
	data, err := ioutil.ReadFile(statePath)
	if err != nil {
		log.Fatalf("Can't read state from %s: %v", statePath, err)
	}
	err = yaml.Unmarshal(data, state)
	if err != nil {
		log.Fatalf("Can't unmarshall sate: %v", err)
	}
	return state
}
