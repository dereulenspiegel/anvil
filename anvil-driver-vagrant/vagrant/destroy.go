package vagrant

import "fmt"

func (v *Vagrant) Destroy(name string) error {
	params := make([]string, 1, 2)
	params[0] = "destroy"
	if name != "" {
		params = append(params, name)
	}

	out, err := v.runCommand(params)
	if err != nil {
		return err
	}

	if outErr := out.Error(); outErr != nil {
		return fmt.Errorf("Error destroying machine: %s", outErr.Data)
	}

	return nil
}
