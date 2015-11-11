package destroy

import (
	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
)

type DestroyCommand struct {
}

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "destroy",
		Usage:  "destroy [regexp]",
		Action: destroyAction,
	}
}

func destroyAction(ctx *cli.Context) {
	filteredCases := commands.GetTestCases(ctx)
	for _, testCase := range filteredCases {
		err := testCase.Transition(test.DESTROYED)
		if err != nil {
			commands.Error(err)
		}
	}
}
