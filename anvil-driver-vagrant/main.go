package main

import (
	"io/ioutil"
	"log"
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

type VagrantDriver struct{}

func convertSshConfig(cfg *vagrant.SshConfig) apis.Connection {
	connection := apis.Connection{}
	connection.Type = apis.SSH
	options := make(map[string]interface{})
	options["Host"] = cfg.Host

	for key, value := range cfg.Options {
		options[key] = value
	}
	connection.Config = options
	return connection
}

func addSshConfig(v *vagrant.Vagrant, inst *apis.Instance) (*apis.Instance, error) {
	sshCfg, err := v.SshConfig("")
	if err != nil {
		return nil, err
	}
	anvilSshCfg := convertSshConfig(sshCfg)
	inst.Connection = anvilSshCfg
	return inst, nil
}

func (v *VagrantDriver) Init(options map[string]interface{}) error {
	return nil
}

func (v *VagrantDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	if util.FileExists(path.Join(DefaultAnvilFolder, DefaultVargantSubfolder, name, "Vagrantfile")) {
		status, err := vagrantClient.Status("")
		if err != nil {
			return apis.Instance{}, err
		}
		instance, err := addSshConfig(vagrantClient, &apis.Instance{Name: name, State: statusMap[status[0].State]})
		if err != nil {
			return apis.Instance{}, err
		}
		return *instance, err
	}
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
	return apis.Instance{Name: name, State: apis.CREATED}, err
}

func (v *VagrantDriver) StartInstance(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	err := vagrantClient.Up("")
	for status, err := vagrantClient.Status(""); len(status) == 0 && status[0].State != "running" && err == nil; {
		if err != nil {
			return apis.Instance{}, err
		}
	}
	instance, err := addSshConfig(vagrantClient, &apis.Instance{Name: name, State: apis.STARTED})
	if err != nil {
		return apis.Instance{}, err
	}
	return *instance, err
}

func (v *VagrantDriver) StopInstance(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	err := vagrantClient.Halt("")
	return apis.Instance{Name: name, State: apis.STOPPED}, err
}

func (v *VagrantDriver) DestroyInstance(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	err := vagrantClient.Destroy("")
	return apis.Instance{Name: name, State: apis.DESTROYED}, err
}

func (v *VagrantDriver) RebootInstance(name string) (apis.Instance, error) {
	vagrantClient := getVagrant(name)
	instance, err := addSshConfig(vagrantClient, &apis.Instance{Name: name, State: apis.STARTED})
	if err != nil {
		return apis.Instance{}, err
	}
	return *instance, err
}

func (v *VagrantDriver) ListInstances() ([]apis.Instance, error) {
	instancesPath := path.Join(DefaultAnvilFolder, DefaultVargantSubfolder)
	fileinfos, err := ioutil.ReadDir(instancesPath)
	if err != nil {
		return nil, err
	}
	instances := make([]apis.Instance, 0, len(fileinfos))
	for _, fi := range fileinfos {
		if fi.IsDir() {
			v := getVagrant(fi.Name())
			status, err := v.Status("")
			if err != nil {
				return nil, err
			}
			instance := &apis.Instance{
				Name:  fi.Name(),
				State: statusMap[status[0].State],
			}
			instance, err = addSshConfig(v, instance)
			if err != nil {
				return nil, err
			}
			instances = append(instances, *instance)
		}
	}
	return instances, nil
}

func (v *VagrantDriver) UpdateState(name string) (apis.Instance, error) {
	if !util.FileExists(path.Join(DefaultAnvilFolder, DefaultVargantSubfolder, name, "Vagrantfile")) {
		return apis.Instance{Name: name, State: apis.DESTROYED}, nil
	}
	vagrantClient := getVagrant(name)
	status, err := vagrantClient.Status("")
	if len(status) == 0 {
		return apis.Instance{Name: name, State: apis.DESTROYED}, nil
	}
	instance, err := addSshConfig(vagrantClient, &apis.Instance{Name: name, State: statusMap[status[0].State]})
	if err != nil {
		return apis.Instance{}, err
	}
	return *instance, err
}

func main() {
	err := apis.RegisterDriverPlugin(&VagrantDriver{})
	if err != nil {
		log.Panicf("Can't register Vagrant Driver plugin: %v", err)
	}
}
