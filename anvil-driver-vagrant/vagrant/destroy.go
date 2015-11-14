package vagrant

func (v *Vagrant) Destroy(name string) error {
	params := make([]string, 1, 2)
	params[0] = "destroy"
	if name != "" {
		params = append(params, name)
	}

	_, err := v.runCommand(params)
	if err != nil {
		return err
	}

	return nil
}
