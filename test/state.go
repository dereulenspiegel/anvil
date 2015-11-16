package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/dereulenspiegel/anvil/util"
	"gopkg.in/yaml.v2"
)

var (
	DefaultStateFolder = ".anvil"
	DefaultStatePrefix = "state"
)

func (t *TestCase) PersistState() {
	out, err := yaml.Marshal(t)
	if err != nil {
		log.Fatalf("Can't marshal state: %v", err)
	}
	err = util.CreateDirectoryIfNotExists(DefaultStateFolder)
	if err != nil {
		log.Panicf("Can't create state directory %s: %v", DefaultStateFolder, err)
	}
	statePath := path.Join(DefaultStateFolder, fmt.Sprintf("%s_%s.yml", DefaultStatePrefix, t.Name))
	err = ioutil.WriteFile(statePath, out, 0655)
	if err != nil {
		log.Fatalf("Can' write state to %s: %v", statePath, err)
	}
}

func (t *TestCase) LoadState() {
	stateFilePath := path.Join(DefaultStateFolder, fmt.Sprintf("%s_%s.yml", DefaultStatePrefix, t.Name))
	if !util.FileExists(stateFilePath) {
		return
	}
	state := &TestCase{}
	data, err := ioutil.ReadFile(stateFilePath)
	if err != nil {
		log.Fatalf("Can't read state from %s: %v", stateFilePath, err)
	}
	err = yaml.Unmarshal(data, state)
	if err != nil {
		log.Fatalf("Can't unmarshall sate: %v", err)
	}
	t.State = state.State
}
