package vagrant

import (
	"fmt"
	"os"
)

func (v *Vagrant) Up(machineName string) error {
	params := make([]string, 1, 2)
	params[0] = "up"
	if machineName != "" {
		params = append(params, machineName)
	}
	_, err := v.runCommand(params)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Up received error: %v", err))
		return err
	}
	return nil
}
