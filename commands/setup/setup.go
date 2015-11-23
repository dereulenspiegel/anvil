package setup

import (
	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
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
		err := testCase.Transition(test.SETUP)
		if err != nil {
			commands.Error(err)
		}
	}
}
