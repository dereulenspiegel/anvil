package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type Provisioner interface {
	Init(opts map[string]interface{}) error
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

func (p *ProvisonerWrapper) Init(opts map[string]interface{}, result *string) error {
	result = nil
	return p.impl.Init(opts)
}

func RegisterProvisionerPlugin(provisioner Provisioner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Provisioner", &ProvisonerWrapper{provisioner}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
