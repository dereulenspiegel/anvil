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

type VerifyResultFormatterPlugin struct {
	rpcClient *rpc.Client
}

var formatterPlugin *VerifyResultFormatterPlugin

func (p *VerifyResultFormatterPlugin) mustCall(method string, args interface{}, reply interface{}) {
	mustCall(p.rpcClient, fmt.Sprintf("VerifyResultFormatter.%s", method), args, reply)
}

func (p *VerifyResultFormatterPlugin) call(method string, args interface{}, reply interface{}) error {
	return p.rpcClient.Call(fmt.Sprintf("VerifyResultFormatter.%s", method), args, reply)
}

func (v *VerifyResultFormatterPlugin) Format(results apis.VerifyResult) ([]byte, error) {
	out := &apis.VerifyFormatterResult{}
	err := v.call("Format", results, out)
	if err != nil {
		return nil, err
	}
	return out.Data, out.Error
}

func LoadVerifyResultFormatter(name string) *VerifyResultFormatterPlugin {
	if formatterPlugin == nil {
		driverPath := fmt.Sprintf("anvil-formatter-%s", name)
		client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
		if err != nil {
			log.Fatalf("Can't load verify result formatter %s", name)
		}
		formatterPlugin = &VerifyResultFormatterPlugin{client}
	}
	return formatterPlugin
}
