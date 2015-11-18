package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dereulenspiegel/anvil/plugin/apis"
)

type AnsibleProvisioner struct{}

func (a *AnsibleProvisioner) Provision(inst apis.Instance, opts map[string]interface{}) error {

	extraVars, _ := opts["extra_vars"].(map[string]string)
	ansibleCfg, exists := opts["ansible_cfg"].(string)
	if !exists {
		ansibleCfg = ""
	}
	playbook, exists := opts["playbook"]
	if !exists {
		return fmt.Errorf("No playbook specified")
	}
	err := runAnsible(inst, playbook.(string), extraVars, ansibleCfg)
	return err
}

var (
	mappedSshParams = map[string]string{
		"HostName":     "ansible_host",
		"Port":         "ansible_port",
		"User":         "ansible_user",
		"IdentityFile": "ansible_ssh_private_key_file",
	}
)

func generateInventory(inst apis.Instance) (string, error) {
	if inst.Connection.Type != apis.SSH {
		return "", fmt.Errorf("%s is not a supported connection type for the ansible provisioner", inst.Connection.Type)
	}
	var buffer bytes.Buffer
	hostname, exists := inst.Connection.Config["Host"].(string)
	if !exists {
		hostname = "default"
	}
	buffer.WriteString(hostname)
	buffer.WriteString(" ")
	for key, value := range inst.Connection.Config {
		ansible_key, mapped := mappedSshParams[key]
		if mapped {
			buffer.WriteString(fmt.Sprintf("%s=%s ", ansible_key, value))
		}
	}
	buffer.WriteString(",")
	return buffer.String(), nil
}

func generateExtraVars(extraVars map[string]string) string {
	if len(extraVars) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("-e=\"")
	for key, value := range extraVars {
		buffer.WriteString(fmt.Sprintf("%s=%s ", key, value))
	}
	buffer.WriteString("\"")
	return buffer.String()
}

func runAnsible(inst apis.Instance, playbook string, extraVars map[string]string, ansibleCfgPath string) error {
	extraVarsParam := generateExtraVars(extraVars)
	inventory, err := generateInventory(inst)
	if err != nil {
		return err
	}
	ansibleCmd := exec.Command("ansible-playbook", extraVarsParam, fmt.Sprintf("-i \"%s\"", inventory), playbook)
	if ansibleCfgPath != "" {
		ansibleCmd.Env = append(ansibleCmd.Env, fmt.Sprintf("ANSIBLE_CONFIG=%s", ansibleCfgPath))
	}
	stderr, err := ansibleCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Can't connect to ansibles stderr")
	}
	stdout, err := ansibleCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Can't connect to ansibles stdout")
	}
	go func() {
		buffer := make([]byte, 0, 1024)
		for n, err := stderr.Read(buffer); err != nil; {
			os.Stderr.Write(buffer[:n])
		}
	}()

	go func() {
		buffer := make([]byte, 0, 1024)
		for n, err := stdout.Read(buffer); err != nil; {
			os.Stderr.Write(buffer[:n])
		}
	}()
	err = ansibleCmd.Start()
	if err != nil {
		return err
	}
	err = ansibleCmd.Wait()
	return err
}

func main() {
	err := apis.RegisterProvisionerPlugin(&AnsibleProvisioner{})
	if err != nil {
		log.Panicf("Can't register ansible Provisioner plugin: %v", err)
	}
}
