package test

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/test"
)

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "test",
		Usage:  "test [regexp]",
		Action: commands.AnvilAction(testAction),
		Flags:  buildFlags(),
	}
}

func buildFlags() []cli.Flag {
	flags := make([]cli.Flag, 0, 2)
	doNotDestroyFlag := cli.BoolFlag{
		Name:  "dont-destroy,d",
		Usage: "Don't destroy instances in case of failure, so you can inspect them",
	}
	flags = append(flags, doNotDestroyFlag)
	return flags
}

func testAction(testCases []*test.TestCase, ctx *cli.Context) {
	failed := false
	for _, testCase := range testCases {
		err := testCase.Transition(test.VERIFIED)
		if err != nil {
			commands.Error(err)
		}
		if testCase.State == test.FAILED {
			failed = true
		}
		if !ctx.Bool("dont-destroy") {
			err = testCase.Transition(test.DESTROYED)
			if err != nil {
				commands.Error(err)
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}
