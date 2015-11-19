package apis

import (
	"fmt"
	"os"
)

type RpcParams map[string]interface{}

var (
	DefaultAnvilFolder = ".anvil"
	DefaultTestFolder  = "tests"
)

type DriverPluginResults struct {
	DriverInstance Instance
}

func Log(message string) {
	os.Stderr.WriteString(message + "\n")
}

func Logf(message string, params ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(message+"\n", params))
}
