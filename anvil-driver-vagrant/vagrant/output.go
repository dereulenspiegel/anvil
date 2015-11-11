package vagrant

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type MessageType string

const (
	BOX_NAME          MessageType = "box-name"
	BOX_PROVIDER      MessageType = "box-provider"
	CLI_COMMAND       MessageType = "cli-command"
	ERROR_EXIT        MessageType = "error-exit"
	PROVIDER_NAME     MessageType = "provider-name"
	STATE             MessageType = "state"
	STATE_HUMAN_LONG  MessageType = "state-human-long"
	STATE_HUMAN_SHORT MessageType = "state-human-short"
	SSH_CONFIG        MessageType = "ssh-config"
)

type OutputLine struct {
	Timestamp time.Time
	Target    string
	Type      MessageType
	Data      string
}

type Output []OutputLine

func (o Output) GetData(target string, msgType MessageType) string {
	for _, out := range o {
		if out.Type == msgType {
			if out.Target == target || target == "" {
				return out.Data
			}
		}
	}
	return ""
}

func (o Output) HasError() bool {
	return o.Error() != nil
}

// TODO is more than one error-exit line possible?
func (o Output) Error() *OutputLine {
	for _, out := range o {
		if out.Type == ERROR_EXIT {
			return &out
		}
	}
	return nil
}

func (o Output) Targets() []string {
	targets := make([]string, 0, 10)
	for _, out := range o {
		var found bool = false
		for _, name := range targets {
			if name == out.Target {
				found = true
			}
		}
		if !found {
			targets = append(targets, out.Target)
		}
	}
	return targets
}

func ParseOutputMessage(line string) (OutputLine, error) {
	parts := strings.Split(line, ",")
	if len(parts) != 4 {
		return OutputLine{}, fmt.Errorf("Unparseable line %s", line)
	}
	intTime, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return OutputLine{}, fmt.Errorf("Unparseable timestamp %s: %v", parts[0], err)
	}
	data := strings.Replace(parts[3], "%!(VAGRANT_COMMA)", ",", -1)
	outputLine := OutputLine{
		Timestamp: time.Unix(intTime, 0),
		Target:    parts[1],
		Type:      MessageType(parts[2]),
		Data:      data,
	}
	return outputLine, nil
}
