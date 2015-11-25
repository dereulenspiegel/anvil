package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/dereulenspiegel/anvil/plugin/apis"
)

type DummyDriver struct{}

func (d *DummyDriver) Init(options map[string]interface{}) error {
	// Nothing to do
	return nil
}

func (d *DummyDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.CREATED}, nil
}

func (d *DummyDriver) StartInstance(name string) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (d *DummyDriver) StopInstance(name string) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.STOPPED}, nil
}

func (d *DummyDriver) DestroyInstance(name string) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.DESTROYED}, nil
}

func (d *DummyDriver) RebootInstance(name string) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (d *DummyDriver) ListInstances() ([]apis.Instance, error) {
	// Nothing to do
	return make([]apis.Instance, 0, 10), fmt.Errorf("Not implemented")
}

func (d *DummyDriver) UpdateState(name string) (apis.Instance, error) {
	// Nothing to do
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (d *DummyDriver) EnableSmartConnection(name string) (string, error) {
	smartSock, path, err := apis.StartSmartConnection(name)
	if err != nil {
		return "", nil
	}
	localShell := exec.Command("/bin/bash")
	localShell.Stdin = smartSock
	localShell.Stdout = smartSock
	err = localShell.Start()
	// TODO verify we can do this
	return path, err
}

func main() {
	err := apis.RegisterDriverPlugin(&DummyDriver{})
	if err != nil {
		log.Panicf("Can't register Local Driver plugin")
	}
}
