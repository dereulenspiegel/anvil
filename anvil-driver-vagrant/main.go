package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/dereulenspiegel/anvil/anvil-driver-vagrant/vagrant"
	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/dereulenspiegel/anvil/util"
)

var (
	DefaultAnvilFolder      = ".anvil"
	DefaultVargantSubfolder = "vagrant"
	Vagrant                 *vagrant.Vagrant

	// TODO map more states
	statusMap = map[string]apis.InstanceState{
		"not_created": apis.DESTROYED,
		"running":     apis.STARTED,
	}
)

func getVagrant(instanceName string) *vagrant.Vagrant {
	workdir := path.Join(DefaultAnvilFolder, DefaultVargantSubfolder, instanceName)
	util.CreateDirectoryIfNotExists(workdir)
	return vagrant.NewVagrant(workdir)
}

type VagrantDriver struct {
}

func (v *VagrantDriver) Init(options map[string]interface{}) error {
	os.Stderr.WriteString("[Vagrant Driver] Init called\n")
	return nil
}

func (v *VagrantDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	vagrantOpts := vagrant.InitOptions{}
	boxName, exists := options["box"].(string)
	if exists {
		vagrantOpts.Name = boxName
	}
	boxUrl, exists := options["url"].(string)
	if exists {
		vagrantOpts.Url = boxUrl
	}
	var err error
	err = vagrantClient.Init(vagrantOpts)
	os.Stderr.WriteString(fmt.Sprintf("[Vagrant Driver] Create instance %s called with options %v\n", name, options))
	return apis.Instance{Name: name, State: apis.CREATED}, err
}

func (v *VagrantDriver) StartInstance(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	err := vagrantClient.Up("")
	os.Stderr.WriteString(fmt.Sprintf("[Vagrant Driver] Start instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STARTED}, err
}

func (v *VagrantDriver) StopInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Stop instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STOPPED}, nil
}

func (v *VagrantDriver) DestroyInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Destroy instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.DESTROYED}, nil
}

func (v *VagrantDriver) RebootInstance(name string) (apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Reboot instance %s called\n", name))
	return apis.Instance{Name: name, State: apis.STARTED}, nil
}

func (v *VagrantDriver) ListInstances() ([]apis.Instance, error) {
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] List instances called\n"))
	return make([]apis.Instance, 0, 10), nil
}

func (v *VagrantDriver) UpdateState(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	status, err := vagrantClient.Status("")
	os.Stderr.WriteString(fmt.Sprintf("[Dummy Driver] Updating state of %s\n", name))
	return apis.Instance{Name: status.Name, State: statusMap[status.State]}, err
}

func main() {
	err := apis.RegisterDriverPlugin(&VagrantDriver{})
	if err != nil {
		log.Panicf("Can't register Vagrant Driver plugin: %v", err)
	}
}
