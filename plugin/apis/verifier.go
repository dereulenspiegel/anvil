package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type VerifyResult struct {
}

type Verifier interface {
	Verify(inst Instance) (VerifyResult, error)
}

type VerifierWrapper struct {
	impl Verifier
}

func (v *VerifierWrapper) Verify(params RpcParams, result *VerifyResult) error {
	_, err := v.impl.Verify(params["instance"].(Instance))
	// TODO copy meaningful results to object in result pointer
	return err
}

func RegisterVerifierPlugin(verifier Verifier) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Driver", &VerifierWrapper{verifier}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
