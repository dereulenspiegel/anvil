package commands

import (
	"log"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/config"
	"github.com/dereulenspiegel/anvil/test"
)

type CliCommand interface {
	SubCommand() cli.Command
}

type BuildCliCommand func() *CliCommand

func FilterTestCases(in []*test.TestCase, expression string) []*test.TestCase {
	caseMatcher := regexp.MustCompile(expression)
	filteredCases := make([]*test.TestCase, 0, len(in))
	for _, testCase := range in {
		if caseMatcher.MatchString(testCase.Name) {
			filteredCases = append(filteredCases, testCase)
		}
	}
	return filteredCases
}

func GetTestCases(ctx *cli.Context) []*test.TestCase {
	caseRegexp := ".*"
	if len(ctx.Args()) > 0 {
		caseRegexp = ctx.Args()[0]
	}
	testCases := test.CompileTestCasesFromConfig(config.Cfg)
	return FilterTestCases(testCases, caseRegexp)
}

func Error(err error) {
	log.Printf("[ERROR] %s", err)
}
