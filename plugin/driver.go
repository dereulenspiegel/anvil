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

type DriverPlugin struct {
	rpcClient *rpc.Client
}

func LoadDriver(name string) *DriverPlugin {
	driverPath := fmt.Sprintf("anvil-driver-%s", name)
	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, driverPath)
	if err != nil {
		log.Fatalf("Can't load driver %s", name)
	}
	return &DriverPlugin{client}
}

func (d *DriverPlugin) mustCall(method string, args interface{}, reply interface{}) {
	mustCall(d.rpcClient, fmt.Sprintf("Driver.%s", method), args, reply)
}

func (d *DriverPlugin) call(method string, args interface{}, reply interface{}) error {
	return d.rpcClient.Call(fmt.Sprintf("Driver.%s", method), args, reply)
}

func (d *DriverPlugin) Close() {
	d.rpcClient.Close()
}

func (d *DriverPlugin) Init(options map[string]interface{}) {
	d.mustCall("Init", apis.RpcParams{
		"options": options,
	}, nil)
}

func (d *DriverPlugin) CreateInstance(name string, options map[string]interface{}) (apis.Instance, error) {
	result := apis.DriverPluginResults{}
	err := d.call("CreateInstance", apis.RpcParams{
		"name":    name,
		"options": options,
	}, &result)
	return result.DriverInstance, err
}

func (d *DriverPlugin) StartInstance(name string) (apis.Instance, error) {
	result := apis.DriverPluginResults{}
	err := d.call("StartInstance", apis.RpcParams{
		"name": name,
	}, &result)
	return result.DriverInstance, err
}

func (d *DriverPlugin) DestroyInstance(name string) (apis.Instance, error) {
	result := apis.DriverPluginResults{}
	err := d.call("DestroyInstance", apis.RpcParams{
		"name": name,
	}, &result)
	return result.DriverInstance, err
}

func (d *DriverPlugin) RebootInstance(name string) (apis.Instance, error) {
	result := apis.DriverPluginResults{}
	err := d.call("RebootInstance", apis.RpcParams{
		"name": name,
	}, &result)
	return result.DriverInstance, err
}

func (d *DriverPlugin) ListInstances() ([]apis.Instance, error) {
	instances := make([]apis.Instance, 0, 10)
	err := d.call("ListInstances", nil, &instances)
	return instances, err
}

func (d *DriverPlugin) UpdateState(name string) (apis.Instance, error) {
	result := apis.DriverPluginResults{}
	err := d.call("UpdateState", apis.RpcParams{
		"name": name,
	}, &result)
	return result.DriverInstance, err
}
