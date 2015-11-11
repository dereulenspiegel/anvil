package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type Provisioner interface {
	Provision(inst Instance) error
}

type ProvisonerWrapper struct {
	impl Provisioner
}

func (p *ProvisonerWrapper) Provision(params RpcParams, result *string) error {
	inst := params["instance"].(Instance)
	result = nil
	return p.impl.Provision(inst)
}

func RegisterProvisionerPlugin(provisioner Provisioner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Provisioner", &ProvisonerWrapper{provisioner}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
