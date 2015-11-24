package apis

import (
	"net/rpc/jsonrpc"

	"github.com/dereulenspiegel/anvil/config"
	"github.com/natefinch/pie"
)

type VerifyResult struct {
	Verifier    string
	CaseResults []VerifyCaseResult
}

type VerifyCaseResult struct {
	Success  bool
	Name     string
	Message  string
	Output   string
	ErrorMsg string
}

type Verifier interface {
	Verify(inst Instance, suite *config.SuiteConfig) (VerifyResult, error)
}

type VerifierWrapper struct {
	impl Verifier
}

type VerifyParams struct {
	Inst  Instance
	Suite *config.SuiteConfig
}

func (v *VerifierWrapper) Verify(params VerifyParams, result *VerifyResult) error {
	verifyResult, err := v.impl.Verify(params.Inst, params.Suite)
	result.Verifier = verifyResult.Verifier
	result.CaseResults = verifyResult.CaseResults
	return err
}

func RegisterVerifierPlugin(verifier Verifier) error {
	p := pie.NewProvider()
	if err := p.RegisterName("Verifier", &VerifierWrapper{verifier}); err != nil {
		return err
	}
	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
