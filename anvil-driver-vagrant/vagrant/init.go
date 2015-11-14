package vagrant

// TODO Probably creating the Vagrantfile with an template would be more flexibel
func (v *Vagrant) Init(opts InitOptions) error {
	_, err := v.runCommand([]string{"init", opts.Name, opts.Url})
	if err != nil {
		return err
	}
	return nil
}
