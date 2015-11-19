package plugin

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/natefinch/pie"
)

type ProvisionerPlugin struct {
	rpcClient *rpc.Client
}

func (p *ProvisionerPlugin) mustCall(method string, args interface{}, reply interface{}) {
	mustCall(p.rpcClient, fmt.Sprintf("Provisioner.%s", method), args, reply)
}

func (p *ProvisionerPlugin) call(method string, args interface{}, reply interface{}) error {
	return p.rpcClient.Call(fmt.Sprintf("Provisioner.%s", method), args, reply)
}

func (p *ProvisionerPlugin) Provision(inst apis.Instance, opts map[string]interface{}) error {
	return p.call("Provision", apis.ProvisionerParams{
		Inst: inst,
		Opts: opts,
	}, nil)
}

func LoadProvisioner(name string) *ProvisionerPlugin {
	driverPath := fmt.Sprintf("anvil-provisioner-%s", name)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {
		log.Fatalf("Can't load provisioner %s", name)
	}
	return &ProvisionerPlugin{client}
}
