package vagrant

func (v *Vagrant) Destroy(name string) error {
	params := make([]string, 2, 3)
	params[0] = "destroy"
	params[1] = "-f"
	if name != "" {
		params = append(params, name)
	}

	_, err := v.runCommand(params)
	if err != nil {
		return err
	}

	return nil
}
