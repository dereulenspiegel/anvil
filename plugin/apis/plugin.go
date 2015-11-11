package apis

import (
	"fmt"
	"os"
)

type RpcParams map[string]interface{}

type DriverPluginResults struct {
	DriverInstance Instance
}

type Log struct{}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) PrintString(in string) {
	os.Stderr.WriteString(in)
}

func (l *Log) Printf(in string, args ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(in, args))
}
