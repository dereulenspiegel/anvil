package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/dereulenspiegel/anvil/util"
)

var (
	DefaultAnvilFolder      = ".anvil"
	DefaultAnsibleSubfolder = "ansible"
)

type AnsibleProvisioner struct{}

func (a *AnsibleProvisioner) Provision(inst apis.Instance, opts map[string]interface{}) error {

	extraVars, _ := opts["extra_vars"].(map[string]string)
	ansibleCfg, exists := opts["ansible_cfg"].(string)
	if !exists {
		ansibleCfg = ""
	}
	playbook, exists := opts["playbook"]
	if !exists || playbook == "" {
		return fmt.Errorf("No playbook specified")
	}
	err := runAnsible(inst, playbook.(string), extraVars, ansibleCfg)
	return err
}

var (
	mappedSshParams = map[string]string{
		"HostName":     "ansible_ssh_host",
		"Port":         "ansible_ssh_port",
		"User":         "ansible_ssh_user",
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
	buffer.WriteString("    ")
	for key, value := range inst.Connection.Config {
		ansibleKey, mapped := mappedSshParams[key]
		if mapped {
			buffer.WriteString(fmt.Sprintf("%s=%s    ", ansibleKey, value))
		}
	}
	buffer.WriteString("\n")

	err := util.CreateDirectoryIfNotExists(path.Join(DefaultAnvilFolder, DefaultAnsibleSubfolder))
	if err != nil {
		return "", err
	}
	tempFile, err := ioutil.TempFile(path.Join(DefaultAnvilFolder, DefaultAnsibleSubfolder), "ansible_inventory")
	if err != nil {
		return "", err
	}
	tempFile.Write(buffer.Bytes())
	if err := tempFile.Sync(); err != nil {
		return "", err
	}
	if err := tempFile.Close(); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func generateExtraVars(extraVars map[string]string) string {
	if len(extraVars) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	buffer.WriteString("\"")
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
	defer os.Remove(inventory)

	params := make([]string, 0, 10)
	params = append(params, "-i")
	params = append(params, inventory)
	params = append(params, "-vvvv")
	if len(extraVars) > 0 {
		params = append(params, "-e")
		params = append(params, extraVarsParam)
	}
	params = append(params, playbook)
	ansibleCmd := exec.Command("ansible-playbook", params...)
	ansibleCmd.Env = append(ansibleCmd.Env, "ANSIBLE_HOST_KEY_CHECKING=False")
	ansibleCmd.Env = append(ansibleCmd.Env, "ANSIBLE_FORCE_COLOR=True")
	if ansibleCfgPath != "" {
		ansibleCmd.Env = append(ansibleCmd.Env, fmt.Sprintf("ANSIBLE_CONFIG=%s", ansibleCfgPath))
	}
	ansibleCmd.Stderr = os.Stderr
	ansibleCmd.Stdout = os.Stderr
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
