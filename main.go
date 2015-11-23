package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands/destroy"
	"github.com/dereulenspiegel/anvil/commands/provision"
	"github.com/dereulenspiegel/anvil/commands/setup"
	"github.com/dereulenspiegel/anvil/commands/status"
	"github.com/dereulenspiegel/anvil/commands/verify"
	"github.com/dereulenspiegel/anvil/config"
)

var (
	App *cli.App
)

func createFlags() []cli.Flag {
	configFlag := cli.StringFlag{
		Name:   "config, c",
		Value:  config.DefaultConfigPath,
		Usage:  "Specify an alternative config file",
		EnvVar: "ANVIL_YAML",
	}

	updateFlag := cli.BoolFlag{
		Name:   "no-refresh,n",
		Usage:  "Trust state in local cache and do not update it from driver each run",
		EnvVar: "ANVIL_NO_REFRESH",
	}

	flags := make([]cli.Flag, 0, 5)
	flags = append(flags, configFlag, updateFlag)
	return flags
}

func createSubCommands(app *cli.App) {
	setup.BuildCommand(app)
	destroy.BuildCommand(app)
	status.BuildCommand(app)
	provision.BuildCommand(app)
	verify.BuildCommand(app)
}

func main() {
	App = cli.NewApp()
	App.Flags = createFlags()
	App.Name = "Anvil"
	App.Author = "Till Klocke"
	App.Copyright = "MIT License"
	App.Version = "0.0.1-alpha"
	App.Usage = "Forge your infrastructure"
	createSubCommands(App)

	App.Run(os.Args)
}
