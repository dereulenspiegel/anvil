package login

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
	"github.com/dereulenspiegel/anvil/commands"
	"github.com/dereulenspiegel/anvil/plugin/apis"
	"github.com/dereulenspiegel/anvil/test"
	"github.com/dereulenspiegel/anvil/util"
	"github.com/ttacon/chalk"
)

func BuildCommand(app *cli.App) {
	app.Commands = append(app.Commands, SubCommand())
}

func SubCommand() cli.Command {
	return cli.Command{
		Name:   "login",
		Usage:  "login [regexp]",
		Action: commands.AnvilAction(loginAction),
	}
}

func loginAction(testCases []*test.TestCase, ctx *cli.Context) {
	if len(testCases) != 1 {
		fmt.Printf("%s[ERROR] You have to specify a single instance to login%s\n", chalk.Red, chalk.Reset)
		os.Exit(1)
	}
	if testCases[0].State == test.DESTROYED {
		fmt.Printf("%s%s is not running%s\n", chalk.Red, testCases[0].Name, chalk.Reset)
		os.Exit(1)
	}
	connection := testCases[0].Instance.Connection

	switch connection.Type {
	case apis.SSH:
		loginViaSsh(testCases[0].Instance)
	}
}

func loginViaSsh(inst apis.Instance) {
	args := make([]string, 0, 10)
	sshConfigPath, err := util.GenerateTempSshConfig(inst.Connection, path.Join(apis.DefaultAnvilFolder, "ssh"))
	if err != nil {
		log.Fatalf("%sCan't generate SSH config: %v%s", chalk.Red, err, chalk.Reset)
	}
	defer os.Remove(sshConfigPath)
	args = append(args, "-F")
	args = append(args, sshConfigPath)
	args = append(args, inst.Connection.Config["Host"].(string))

	sshCmd := exec.Command("ssh", args...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	err = sshCmd.Start()
	if err != nil {
		log.Fatalf("%sStarting SSH command failed: %v%s", chalk.Red, err, chalk.Reset)
	}
	err = sshCmd.Wait()
	if err != nil {
		log.Fatalf("%sSSH connection terminated abnormally: %v%s", chalk.Red, err, chalk.Reset)
	}
}
