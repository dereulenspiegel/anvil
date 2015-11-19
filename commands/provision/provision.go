package provision

import (
	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
)

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "provision",
		Usage:  "provision [regexp]",
		Action: commands.AnvilAction(provisionAction),
	}
}

func provisionAction(testCases []*test.TestCase, ctx *cli.Context) {
	for _, testCase := range testCases {
		err := testCase.Transition(test.PROVISIONED)
		if err != nil {
			commands.Error(err)
		}
	}
}
