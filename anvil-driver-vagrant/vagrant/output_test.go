package vagrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validLines = []string{
		"1447091727,node01,provider-name,virtualbox",
		"1447091727,,state,not_created",
		"1447091727,node01,state-human-short,not created",
		"1447091727,node01,state-human-long,The environment has not yet been created. Run `vagrant up` to\ncreate the environment. If a machine is not created%!(VAGRANT_COMMA) only the\ndefault provider will be shown. So if a provider is not listed%!(VAGRANT_COMMA)\nthen the machine is not created for that environment.",
		"1447416672,,error-exit,Vagrant::Errors::VagrantfileExistsError,`Vagrantfile` already exists in this directory. Remove it before\nrunning `vagrant init`.",
	}
	invalidLines = []string{
		"node01,provider-name,virtualbox",
		"144ab7091727,node01,provider-name,virtualbox",
	}
)

func TestValidOutputMessageParsing(t *testing.T) {
	assert := assert.New(t)

	out, err := ParseOutputMessage(validLines[0])
	assert.Nil(err)
	assert.Equal("node01", out.Target)
	assert.Equal(PROVIDER_NAME, out.Type)
	assert.Equal("virtualbox", out.Data)

	out, err = ParseOutputMessage(validLines[3])
	assert.Nil(err)
	assert.Equal("node01", out.Target)
	assert.Equal(STATE_HUMAN_LONG, out.Type)
	expectedData := "The environment has not yet been created. Run `vagrant up` to\ncreate the environment. If a machine is not created, only the\ndefault provider will be shown. So if a provider is not listed,\nthen the machine is not created for that environment."
	assert.Equal(expectedData, out.Data)

	out, err = ParseOutputMessage(validLines[1])
	assert.Nil(err)
	assert.Equal("", out.Target)
	assert.Equal(STATE, out.Type)
	assert.Equal("not_created", out.Data)

	out, err = ParseOutputMessage(validLines[4])
	assert.Nil(err)
	assert.Equal("", out.Target)
	assert.Equal(ERROR_EXIT, out.Type)
}

func TestInvalidOutputMessageParsing(t *testing.T) {
	assert := assert.New(t)

	for _, line := range invalidLines {
		_, err := ParseOutputMessage(line)
		assert.NotNil(err)
	}
}
