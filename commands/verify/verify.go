package verify

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
	"github.com/ttacon/chalk"
)

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "verify",
		Usage:  "verify [regexp]",
		Action: commands.AnvilAction(verifyAction),
	}
}

func verifyAction(testCases []*test.TestCase, ctx *cli.Context) {
	for _, testCase := range testCases {
		fmt.Printf("%sVerifying %s%s\n", chalk.Bold, testCase.Name, chalk.Reset)
		err := testCase.Transition(test.VERIFIED)
		if err != nil {
			commands.Error(err)
		}
	}
}
