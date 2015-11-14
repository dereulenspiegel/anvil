package vagrant

func (v *Vagrant) Halt(machineName string) error {
	params := make([]string, 1, 2)
	params[0] = "up"
	if machineName != "" {
		params = append(params, machineName)
	}
	_, err := v.runCommand(params)
	if err != nil {
		return err
	}
	return nil
}
