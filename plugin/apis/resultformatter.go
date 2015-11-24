package apis

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type VerifyResultFormatter interface {
	Format(results VerifyResult) ([]byte, error)
}

type VerifyResultFormatterWrapper struct {
	impl VerifyResultFormatter
}

type VerifyFormatterResult struct {
	Data  []byte
	Error error
}

func (v *VerifyResultFormatterWrapper) Format(results VerifyResult, out *VerifyFormatterResult) error {
	data, err := v.impl.Format(results)
	out.Data = data
	out.Error = err
	return err
}

func RegisterVerifyResultFormatterPlugin(formatter VerifyResultFormatter) error {
	p := pie.NewProvider()
	if err := p.RegisterName("VerifyResultFormatter", &VerifyResultFormatterWrapper{formatter}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
