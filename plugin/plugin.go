package plugin

import (
	"log"
	"net/rpc"
	"os/user"
	"path"
)

const (
	pluginDirName = "plugins"
)

var (
	PluginsDirectory string
)

type Plugin interface {
	Close()
}

func init() {
	user, err := user.Current()
	if err != nil {
		log.Panicf("Can't determine current user")
	}
	PluginsDirectory = path.Join(user.HomeDir, ".anvil", pluginDirName)
}

func mustCall(client *rpc.Client, method string, args interface{}, reply interface{}) {
	err := client.Call(method, args, reply)
	if err != nil {
		log.Panicf("Error calling plugin method %s: %s", method, err)
	}
}
