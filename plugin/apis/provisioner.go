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

func (p *ProvisonerWrapper) Provision(params RpcParams, result *string) error {
	inst := params["instance"].(Instance)
	opts := params["opts"].(map[string]interface{})
	result = nil
	return p.impl.Provision(inst, opts)
}

func RegisterProvisionerPlugin(provisioner Provisioner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Provisioner", &ProvisonerWrapper{provisioner}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
