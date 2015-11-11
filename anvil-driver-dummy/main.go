package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dereulenspiegel/anvil/plugin/apis"
)

type DummyDriver struct{}

func (d *DummyDriver) Init(options map[string]interface{}) error {
	os.Stderr.WriteString("[Dummy Driver] Init called\n")
	return nil
}

func (d *DummyDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Create instance %s called with options %v\n", name, options))
	return apis.Instance{Name: name, State: apis.CREATED}, nil
}

func (d *DummyDriver) StartInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Start instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (d *DummyDriver) StopInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Stop instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STOPPED}, nil
}

func (d *DummyDriver) DestroyInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Destroy instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.DESTROYED}, nil
}

func (d *DummyDriver) RebootInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Reboot instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (d *DummyDriver) ListInstances() ([]apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] List instances called\n"))
	return make([]apis.Instance, 0, 10), nil
}

func (d *DummyDriver) UpdateState(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Updating state of %s\n", name))
	return apis.Instance{Name: name, State: apis.DESTROYED}, nil
}

func main() {
	err := apis.RegisterDriverPlugin(&DummyDriver{})
	if err != nil {
		log.Panicf("Can't register Dummy Driver plugin")
	}
}
