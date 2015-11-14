package vagrant

import (
	"fmt"
	"regexp"
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

var (
	lineRegexp = regexp.MustCompile(`^([\d]+),([\w\d]*),([\w\d\-]+),(?s:(.*))$`)
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
	if !lineRegexp.MatchString(line) {
		return OutputLine{}, fmt.Errorf("Regexp: Unparseable line: %s", line)
	}
	submatches := lineRegexp.FindAllStringSubmatch(line, -1)
	intTime, err := strconv.ParseInt(submatches[0][1], 10, 64)
	if err != nil {
		return OutputLine{}, fmt.Errorf("Unparseable timestamp %s: %v", submatches[0][1], err)
	}
	data := strings.Replace(submatches[0][4], "%!(VAGRANT_COMMA)", ",", -1)
	outputLine := OutputLine{
		Timestamp: time.Unix(intTime, 0),
		Target:    submatches[0][2],
		Type:      MessageType(submatches[0][3]),
		Data:      data,
	}
	return outputLine, nil
}
