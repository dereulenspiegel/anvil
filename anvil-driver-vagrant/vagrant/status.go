package vagrant

import "fmt"

type Status struct {
	Name     string
	Provider string
	State    string
}

func (v *Vagrant) Status(instanceName string) ([]*Status, error) {
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

	instanceNames := out.Targets()
	statusSlice := make([]*Status, 0, len(instanceNames))
	for _, name := range instanceNames {
		provider := out.GetData(name, PROVIDER_NAME)
		state := out.GetData(name, STATE)
		status := &Status{
			Name:     name,
			State:    state,
			Provider: provider,
		}
		statusSlice = append(statusSlice, status)
	}
	return statusSlice, nil
}
