package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dereulenspiegel/anvil/config"
	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/dereulenspiegel/anvil/util"
)

var (
	DefaultTestinfraSubfolder = "testinfra"
	sshConfigPrexix           = "ssh_config_"
)

type TestinfraVerifier struct{}

func (t *TestinfraVerifier) Verify(inst apis.Instance, suite *config.SuiteConfig) (apis.VerifyResult, error) {
	testinfraTestDir := path.Join(apis.DefaultTestFolder, suite.Name, "testinfra")
	testFiles, err := ioutil.ReadDir(testinfraTestDir)
	if err != nil {
		return apis.VerifyResult{}, err
	}
	resultSlice := make([]apis.VerifyCaseResult, 0, 10)
	for _, file := range testFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".py") {
			testFile := path.Join(testinfraTestDir, file.Name())
			var out []byte
			out, err = executeTestinfraFile(testFile, inst)
			result := apis.VerifyCaseResult{
				Name:   file.Name(),
				Error:  err,
				Output: string(out),
			}
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("Testinfra failed to verify %s", file.Name())
			} else {
				result.Success = true
				result.Message = fmt.Sprintf("Testinfra successfully verified %s", file.Name())
			}
			resultSlice = append(resultSlice, result)
		}
	}

	return apis.VerifyResult{
		Verifier:    "testinfra",
		CaseResults: resultSlice,
	}, err
}

func generateConnectionParams(inst apis.Instance) ([]string, error) {
	switch inst.Connection.Type {
	case apis.SSH:
		return generateSshConnectionParams(inst)
	}
	return make([]string, 0, 0), fmt.Errorf("Unknown connection type %s", inst.Connection.Type)
}

func generateSshConnectionParams(inst apis.Instance) ([]string, error) {
	sshParams := make([]string, 0, 10)
	configPath, err := generateTempSshConfig(inst.Connection)
	if err != nil {
		return sshParams, err
	}
	sshParams = append(sshParams, "--connection=paramiko")
	sshParams = append(sshParams, fmt.Sprintf("--hosts=%s", inst.Connection.Config["Host"]))
	sshParams = append(sshParams, fmt.Sprintf("--ssh-config=%s", configPath))
	return sshParams, nil
}

func generateTempSshConfig(conn apis.Connection) (string, error) {
	var buffer bytes.Buffer
	host, exists := conn.Config["Host"].(string)
	if !exists {
		return "", fmt.Errorf("Host name property does not exist in SSH config")
	}
	buffer.WriteString("Host ")
	buffer.WriteString(host)
	buffer.WriteString("\n")
	for key, value := range conn.Config {
		if key != "Host" {
			buffer.WriteString(fmt.Sprintf("  %s    %s\n", key, value))
		}
	}
	util.CreateDirectoryIfNotExists(path.Join(apis.DefaultAnvilFolder, DefaultTestinfraSubfolder))
	tempSshFile, err := ioutil.TempFile(path.Join(apis.DefaultAnvilFolder, DefaultTestinfraSubfolder), sshConfigPrexix)
	if err != nil {
		return "", err
	}
	_, err = tempSshFile.Write(buffer.Bytes())
	if err != nil {
		return "", err
	}
	tempSshFile.Close()
	return tempSshFile.Name(), nil
}

func executeTestinfraFile(testFileName string, inst apis.Instance) ([]byte, error) {
	params := make([]string, 0, 10)
	params = append(params, "--color=yes")
	params = append(params, "--sudo")
	connParams, err := generateConnectionParams(inst)
	if err != nil {
		return nil, err
	}
	params = append(params, connParams...)
	params = append(params, testFileName)
	testinfraCmd := exec.Command("testinfra", params...)
	outBuffer := bytes.NewBuffer(make([]byte, 16032))
	multiWriter := io.MultiWriter(outBuffer)
	testinfraCmd.Stdout = multiWriter
	testinfraCmd.Stderr = multiWriter
	err = testinfraCmd.Start()
	if err != nil {
		return nil, err
	}
	err = testinfraCmd.Wait()
	removeLeftovers()
	return outBuffer.Bytes(), err
}

func removeLeftovers() {
	testInfraDir := path.Join(apis.DefaultAnvilFolder, DefaultTestinfraSubfolder)
	files, err := ioutil.ReadDir(testInfraDir)
	if err != nil {
		apis.Logf("Error while cleaning left over files: %v", err)
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), sshConfigPrexix) {
			path := path.Join(apis.DefaultAnvilFolder, DefaultTestinfraSubfolder, file.Name())
			err := os.Remove(path)
			if err != nil {
				apis.Logf("Can't remove %s: %v", path, err)
			}
		}
	}
}

func main() {
	err := apis.RegisterVerifierPlugin(&TestinfraVerifier{})
	if err != nil {
		log.Panicf("Can't register testinfra Verifier plugin: %v", err)
	}
}
