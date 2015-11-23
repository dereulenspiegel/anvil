package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type Driver interface {
	Init(options map[string]interface{}) error
	CreateInstance(name string, options map[string]interface{}) (Instance, error)
	StartInstance(name string) (Instance, error)
	StopInstance(name string) (Instance, error)
	DestroyInstance(name string) (Instance, error)
	RebootInstance(name string) (Instance, error)
	ListInstances() ([]Instance, error)
	UpdateState(name string) (Instance, error)
}

type DriverWrapper struct {
	impl Driver
}

// Result is not used here
func (d *DriverWrapper) Init(params RpcParams, result *string) error {
	opts, exist := params["options"]
	if exist && opts != nil {
		return d.impl.Init(opts.(map[string]interface{}))
	}
	return nil
}

func (d *DriverWrapper) CreateInstance(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.CreateInstance(params["name"].(string), params["options"].(map[string]interface{}))
	result.DriverInstance = machine
	return err
}

func (d *DriverWrapper) StartInstance(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.StartInstance(params["name"].(string))
	result.DriverInstance = machine
	return err
}

func (d *DriverWrapper) StopInstance(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.StopInstance(params["name"].(string))
	result.DriverInstance = machine
	return err
}

func (d *DriverWrapper) DestroyInstance(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.DestroyInstance(params["name"].(string))
	result.DriverInstance = machine
	return err
}

func (d *DriverWrapper) RebootInstance(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.RebootInstance(params["name"].(string))
	result.DriverInstance = machine
	return err
}

func (d *DriverWrapper) ListInstances(params RpcParams, result *[]Instance) error {
	instances, err := d.impl.ListInstances()
	result = &instances
	return err
}

func (d *DriverWrapper) UpdateState(params RpcParams, result *DriverPluginResults) error {
	machine, err := d.impl.UpdateState(params["name"].(string))
	result.DriverInstance = machine
	return err
}

func RegisterDriverPlugin(driver Driver) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Driver", &DriverWrapper{driver}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
