package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type Provisioner interface {
	Provision(inst Instance, opts map[string]interface{}) error
}

type ProvisonerWrapper struct {
	impl Provisioner
}

type ProvisionerParams struct {
	Inst Instance
	Opts map[string]interface{}
}

func (p *ProvisonerWrapper) Provision(params ProvisionerParams, result *string) error {
	result = nil
	return p.impl.Provision(params.Inst, params.Opts)
}

func RegisterProvisionerPlugin(provisioner Provisioner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Provisioner", &ProvisonerWrapper{provisioner}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
