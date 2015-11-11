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

type VerifierPlugin struct {
	rpcClient *rpc.Client
}

func (d *VerifierPlugin) mustCall(method string, args interface{}, reply interface{}) {
	mustCall(d.rpcClient, fmt.Sprintf("Verifier.%s", method), args, reply)
}

func (d *VerifierPlugin) call(method string, args interface{}, reply interface{}) error {
	return d.rpcClient.Call(fmt.Sprintf("Verifier.%s", method), args, reply)
}

func (v *VerifierPlugin) Verify(inst apis.Instance) (apis.VerifyResult, error) {
	result := &apis.VerifyResult{}
	err := v.call("Verify", apis.RpcParams{
		"instance": inst,
	}, result)
	return *result, err
}

func LoadVerifier(name string) *VerifierPlugin {
	driverPath := fmt.Sprintf("anvil-verifier-%s", name)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {
		log.Fatalf("Can't load verifier %s", name)
	}
	return &VerifierPlugin{client}
}
