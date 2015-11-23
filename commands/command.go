package commands

import (
	"log"
	"regexp"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/config"
	"github.com/dereulenspiegel/anvil/test"
)

type AnvilCommand func(cases []*test.TestCase, ctx *cli.Context)

func AnvilAction(command AnvilCommand) func(*cli.Context) {
	return func(ctx *cli.Context) {
		testCases := GetTestCases(ctx)
		for _, tc := range testCases {
			tc.UpdateState()
		}
		command(testCases, ctx)
		for _, tc := range testCases {
			tc.PersistState()
		}
	}
}

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
