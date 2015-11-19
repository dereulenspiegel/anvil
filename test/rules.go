package test

import (
	"github.com/ryanfaerman/fsm"
)

func createTestCaseRules() fsm.Ruleset {
	rules := fsm.Ruleset{}

	rules.AddTransition(fsm.T{DESTROYED, SETUP})
	rules.AddTransition(fsm.T{SETUP, PROVISIONED})
	rules.AddTransition(fsm.T{PROVISIONED, VERIFIED})

	rules.AddTransition(fsm.T{SETUP, FAILED})
	rules.AddTransition(fsm.T{PROVISIONED, FAILED})
	rules.AddTransition(fsm.T{VERIFIED, FAILED})

	rules.AddTransition(fsm.T{SETUP, DESTROYED})
	rules.AddTransition(fsm.T{PROVISIONED, DESTROYED})
	rules.AddTransition(fsm.T{VERIFIED, DESTROYED})

	rules.AddTransition(fsm.T{PROVISIONED, SETUP})
	rules.AddTransition(fsm.T{PROVISIONED, PROVISIONED})
	rules.AddTransition(fsm.T{VERIFIED, VERIFIED})

	rules.AddTransition(fsm.T{FAILED, VERIFIED})
	rules.AddTransition(fsm.T{FAILED, PROVISIONED})
	return rules
}
