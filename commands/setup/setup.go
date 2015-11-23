package setup

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
	"github.com/ttacon/chalk"
)

type SetupCommand struct {
}

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "setup",
		Usage:  "setup [regexp]",
		Action: commands.AnvilAction(setupAction),
	}
}

func setupAction(testCases []*test.TestCase, ctx *cli.Context) {
	for _, testCase := range testCases {
		fmt.Printf("%sProvisioning %s%s\n", chalk.Bold, testCase.Name, chalk.Reset)
		err := testCase.Transition(test.SETUP)
		if err != nil {
			commands.Error(err)
		}
	}
}
