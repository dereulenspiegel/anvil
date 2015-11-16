package status

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
)

type StatusCommand struct {
}

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "status",
		Usage:  "status [regexp]",
		Action: commands.AnvilAction(statusAction),
	}
}

func statusAction(testCases []*test.TestCase, ctx *cli.Context) {
	for _, testCase := range testCases {
		fmt.Printf("%s \t %s", testCase.Name, testCase.CurrentState())
	}
}
