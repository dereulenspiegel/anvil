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
		Action: commands.AnvilAction(destroyAction),
	}
}

func destroyAction(testCases []*test.TestCase, ctx *cli.Context) {
	for _, testCase := range testCases {
		err := testCase.Transition(test.DESTROYED)
		if err != nil {
			commands.Error(err)
		}
	}
}
