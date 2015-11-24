package main

import (
	"fmt"
	"log"

	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/samalba/dockerclient"
)

type DockerDriver struct{}

var (
	docker         *dockerclient.DockerClient
	containerIdTag = "containerId"
	defaultTimeout = 30
)

func findContainerByName(name string) (string, error) {
	containers, err := docker.ListContainers(true, false, fmt.Sprintf("name=%s", name))
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("Can't find container with name %s", name)
	}
	return containers[0].Id, nil
}

func (d *DockerDriver) Init(options map[string]interface{}) error {
	url := "unix:///var/run/docker.sock"
	if urlParam, exists := options["url"]; exists {
		url = urlParam.(string)
	}
	var err error
	docker, err = dockerclient.NewDockerClient(url, nil)
	docker.Info()
	return err
}

func (d *DockerDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	image := options["image"].(string)
	command := "/bin/bash"
	if cmd, exists := options["command"]; exists {
		command = cmd.(string)
	}
	containerConfig := &dockerclient.ContainerConfig{
		Image:       image,
		Cmd:         []string{command},
		AttachStdin: true,
		Tty:         true,
	}
	_, err := docker.CreateContainer(containerConfig, name)
	if err != nil {
		return apis.Instance{}, err
	}
	return apis.Instance{Name: name, State: apis.CREATED}, nil
}

func (d *DockerDriver) StartInstance(name string) (apis.Instance, error) {
	containerId, err := findContainerByName(name)
	if err != nil {
		return apis.Instance{}, err
	}
	if containerId == "" {
		return apis.Instance{}, fmt.Errorf("Container with name %s not found", name)
	}
	hostConfig := &dockerclient.HostConfig{}
	// TODO probaly make the host config configurable
	err = docker.StartContainer(containerId, hostConfig)
	return apis.Instance{Name: name, State: apis.STARTED}, err
}

func (d *DockerDriver) StopInstance(name string) (apis.Instance, error) {
	containerId, err := findContainerByName(name)
	if err != nil {
		return apis.Instance{}, err
	}
	err = docker.StopContainer(containerId, defaultTimeout)
	return apis.Instance{Name: name, State: apis.STOPPED}, err
}

func (d *DockerDriver) DestroyInstance(name string) (apis.Instance, error) {
	containerId, err := findContainerByName(name)
	if err != nil {
		return apis.Instance{}, err
	}
	err = docker.RemoveContainer(containerId, true, true)
	return apis.Instance{Name: name, State: apis.DESTROYED}, err
}

func (d *DockerDriver) RebootInstance(name string) (apis.Instance, error) {
	containerId, err := findContainerByName(name)
	if err != nil {
		return apis.Instance{}, err
	}
	err = docker.RestartContainer(containerId, defaultTimeout)
	return apis.Instance{Name: name, State: apis.STARTED}, err
}

func (d *DockerDriver) ListInstances() ([]apis.Instance, error) {
	apis.Log("[Docker Driver] List instances is not implemented\n")
	return make([]apis.Instance, 0, 10), fmt.Errorf("Not implemented")
}

func (d *DockerDriver) UpdateState(name string) (apis.Instance, error) {
	containerId, err := findContainerByName(name)
	if err != nil {
		return apis.Instance{Name: name, State: apis.DESTROYED}, nil
	}
	containerInfo, err := docker.InspectContainer(containerId)
	if err != nil {
		return apis.Instance{}, err
	}
	dockerState := containerInfo.State
	apiState := apis.DESTROYED
	if dockerState.Running {
		apiState = apis.STARTED
	} else if dockerState.Paused {
		apiState = apis.CREATED
	} else {
		apiState = apis.STOPPED
	}
	return apis.Instance{Name: name, State: apiState}, nil
}

func main() {
	err := apis.RegisterDriverPlugin(&DockerDriver{})
	if err != nil {
		log.Panicf("Can't register Docker Driver plugin")
	}
}
