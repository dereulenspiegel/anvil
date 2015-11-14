package vagrant

// TODO Probably creating the Vagrantfile with an template would be more flexibel
func (v *Vagrant) Init(opts InitOptions) error {
	params := make([]string, 2, 3)
	params[0] = "init"
	params[1] = opts.Name
	if opts.Url != "" {
		params = append(params, opts.Url)
	}
	_, err := v.runCommand(params)
	if err != nil {
		return err
	}
	return nil
}
