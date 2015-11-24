package test

import (
	"bytes"
	"fmt"

	"github.com/dereulenspiegel/anvil/plugin/apis"
)

type DefaultConsoleResultFormatter struct{}

func (d *DefaultConsoleResultFormatter) Format(results apis.VerifyResult) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s Results:\n", results.Verifier))
	for _, caseResult := range results.CaseResults {
		resultString := "FAILED"
		if caseResult.Success {
			resultString = "SUCCESS"
		}
		buf.WriteString(fmt.Sprintf("\t[%s] %s: %s", caseResult.Name, resultString, caseResult.Message))
	}
	return buf.Bytes(), nil
}
