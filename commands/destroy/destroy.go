package destroy

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
	"github.com/ttacon/chalk"
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
		fmt.Printf("%sDestroying instance %s%s\n", chalk.Bold, testCase.Name, chalk.Reset)
		err := testCase.Transition(test.DESTROYED)
		if err != nil {
			commands.Error(err)
		}
	}
}
