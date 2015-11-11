package vagrant

import "fmt"

type Status struct {
	Name     string
	Provider string
	State    string
}

func (v *Vagrant) Status(instanceName string) (*Status, error) {
	params := make([]string, 1, 2)
	params[0] = "status"
	if instanceName != "" {
		params = append(params, instanceName)
	}
	out, err := v.runCommand(params)
	if err != nil {
		return nil, err
	}
	if len(out.Targets()) > 1 {
		return nil, fmt.Errorf("Trying to get status of multiple targtes")
	}
	provider := out.GetData(instanceName, PROVIDER_NAME)
	state := out.GetData(instanceName, STATE)
	name := out.Targets()[0]
	status := &Status{
		Provider: provider,
		State:    state,
		Name:     name,
	}
	return status, nil
}
