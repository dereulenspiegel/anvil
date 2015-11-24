package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path"

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
	containers, err := docker.ListContainers(true, false, "")
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		return "", fmt.Errorf("Can't find container with name %s, because no containers are there", name)
	}
	for _, container := range containers {
		for _, containerName := range container.Names {
			containerName = containerName[1:]
			if containerName == name {
				return container.Id, nil
			}
		}
	}
	return "", fmt.Errorf("Can't find container with name %s", name)
}

func createConnection(containerId string) apis.Connection {
	conn := apis.Connection{
		Type:   apis.Docker,
		Config: make(map[string]interface{}),
	}
	conn.Config["containerId"] = containerId
	return conn
}

func ensureImageIsAvailable(imageName string) error {
	images, err := docker.ListImages(true)
	if err != nil {
		return err
	}
	found := false
search:
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				found = true
				break search
			}
		}
	}
	if !found {
		err = docker.PullImage(imageName, &dockerclient.AuthConfig{})
		if err != nil {
			return err
		}
	}
	return nil
}

func createTlsConfig(options map[string]interface{}) (tlsConfig *tls.Config, err error) {
	var clientCertFile string
	var clientKeyFile string
	var serverCertFile string
	if cert, exists := options["certPath"]; exists {
		certPath := cert.(string)
		clientCertFile = path.Join(certPath, "cert.pem")
		clientKeyFile = path.Join(certPath, "key.pem")
		serverCertFile = path.Join(certPath, "ca.pem")
	} else {
		clientCert, certExists := options["clientCert"]
		clientKey, keyExists := options["clientKey"]
		if !certExists || keyExists {
			return
		}
		clientCertFile = clientCert.(string)
		clientKeyFile = clientKey.(string)
	}

	clientCert, errTls := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		err = errTls
		return
	}
	tlsConfig = &tls.Config{
		Certificates: make([]tls.Certificate, 1, 1),
	}
	tlsConfig.Certificates[0] = clientCert

	if serverCert, exists := options["serverCert"]; exists {
		serverCertFile = serverCert.(string)
	}

	if serverCertFile != "" {
		pemData, errIo := ioutil.ReadFile(serverCertFile)
		if errIo != nil {
			err = errIo
			return
		}
		if tlsConfig.RootCAs == nil {
			certs := x509.NewCertPool()
			tlsConfig.RootCAs = certs
		}
		tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	}
	tlsConfig.BuildNameToCertificate()
	return
}

func (d *DockerDriver) Init(options map[string]interface{}) error {
	url := "unix:///var/run/docker.sock"
	if urlParam, exists := options["url"]; exists {
		url = urlParam.(string)
	}

	var err error
	tlsConfig, err := createTlsConfig(options)
	if err != nil {
		return err
	}
	docker, err = dockerclient.NewDockerClient(url, tlsConfig)
	return err
}

func (d *DockerDriver) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	image := options["image"].(string)
	command := "/bin/bash"
	if cmd, exists := options["command"]; exists {
		command = cmd.(string)
	}
	ensureImageIsAvailable(image)
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
	return apis.Instance{Name: name, State: apis.STARTED, Connection: createConnection(containerId)}, err
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
	return apis.Instance{Name: name, State: apiState, Connection: createConnection(containerId)}, nil
}

func main() {
	err := apis.RegisterDriverPlugin(&DockerDriver{})
	if err != nil {
		log.Panicf("Can't register Docker Driver plugin")
	}
}
