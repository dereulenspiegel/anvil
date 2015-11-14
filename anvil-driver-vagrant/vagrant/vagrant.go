package vagrant

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Vagrant struct {
	Workdir string
}

type InitOptions struct {
	Name string
	Url  string
}

func NewVagrant(workdir string) *Vagrant {
	return &Vagrant{workdir}
}

func (v *Vagrant) runCommand(params []string) (Output, error) {
	params = append(params, "--machine-readable")
	cmd := exec.Command("vagrant", params...)
	if v.Workdir != "" {
		cmd.Dir = v.Workdir
	}
	out, err := cmd.CombinedOutput()
	// Dirty hack since ssh-config doesn't support machine readable format
	if params[0] == "ssh-config" {
		targetName := "default"
		if len(params) > 1 {
			targetName = params[1]
		}
		output := make(Output, 1, 1)
		output[0].Target = targetName
		output[0].Data = string(out)
		output[0].Timestamp = time.Now()
		output[0].Type = SSH_CONFIG

		return output, nil
	}
	lines := strings.Split(string(out), "\n")
	outputLines := make(Output, 0, len(lines))

	for _, line := range lines {
		if len(line) > 0 {
			msg, err := ParseOutputMessage(line)
			if err != nil {
				return nil, err
			}
			outputLines = append(outputLines, msg)
		}
	}
	if outErr := outputLines.Error(); outErr != nil {
		return outputLines, fmt.Errorf("Vagrant error: %s", outErr.Data)
	}
	return outputLines, err
}
