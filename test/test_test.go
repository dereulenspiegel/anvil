package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestCaseTransition(t *testing.T) {
	assert := assert.New(t)
	testCase := NewTestCase(nil, nil)
	err := testCase.Transition(SETUP)
	assert.Nil(err)
}
