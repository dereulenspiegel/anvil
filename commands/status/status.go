package status

import (
	"fmt"
	"log"

	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
	"github.com/ryanfaerman/fsm"
	"github.com/ttacon/chalk"
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

var (
	colorStatusMap = map[fsm.State]chalk.Color{
		test.DESTROYED:   chalk.Magenta,
		test.VERIFIED:    chalk.Green,
		test.PROVISIONED: chalk.Blue,
		test.SETUP:       chalk.Blue,
		test.FAILED:      chalk.Red,
	}
)

func statusAction(testCases []*test.TestCase, ctx *cli.Context) {
	writer := tabwriter.NewWriter(os.Stdout, 20, 4, 2, '\t', 0)
	for _, testCase := range testCases {
		fmt.Fprintf(writer, "%s%s\t%s%s%s\n", chalk.Bold, testCase.Name,
			colorStatusMap[testCase.State], testCase.State, chalk.Reset)
	}
	err := writer.Flush()
	if err != nil {
		log.Printf("Error writing table: %v", err)
	}
}
