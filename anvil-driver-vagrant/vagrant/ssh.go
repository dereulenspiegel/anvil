package vagrant

import (
	"regexp"
	"strings"
)

type SshConfig struct {
	Host    string
	Options map[string]string
}

var extractHostRegexp = regexp.MustCompile(`^Host\s+(.*)$`)
var extractOptionsRegexp = regexp.MustCompile(`^\s+(\w+)\s+(.*)$`)

func parseSshConfig(in string) (*SshConfig, error) {
	configLines := strings.Split(in, "\n")
	sshConfig := &SshConfig{
		Options: make(map[string]string),
	}
	for _, line := range configLines {
		if extractHostRegexp.MatchString(line) {
			submatches := extractHostRegexp.FindAllStringSubmatch(line, -1)
			sshConfig.Host = submatches[0][1]
		}
		if extractOptionsRegexp.MatchString(line) {
			submatches := extractOptionsRegexp.FindAllStringSubmatch(line, -1)
			key := submatches[0][1]
			value := submatches[0][2]
			sshConfig.Options[key] = value
		}
	}
	return sshConfig, nil
}

func (v *Vagrant) SshConfig(machineName string) (*SshConfig, error) {
	params := make([]string, 1, 2)
	params[0] = "ssh-config"
	if machineName != "" {
		params = append(params, machineName)
	}
	out, err := v.runCommand(params)
	if err != nil {
		return nil, err
	}
	rawConfig := out[0].Data
	sshConfig, err := parseSshConfig(rawConfig)
	return sshConfig, err
}
