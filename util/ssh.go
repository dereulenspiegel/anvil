package util

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/dereulenspiegel/anvil/plugin/apis"
)

func GenerateTempSshConfig(conn apis.Connection, targetPath string) (string, error) {
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
	CreateDirectoryIfNotExists(targetPath)
	tempSshFile, err := ioutil.TempFile(targetPath, "ssh_config_")
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
