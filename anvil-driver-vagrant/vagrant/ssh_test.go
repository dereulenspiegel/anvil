package vagrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sshConfigString = `Host default
  HostName 127.0.0.1
  User vagrant
  Port 2222
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /Volumes/Extras/openwrt-vagrant-buildroot/.vagrant/machines/default/virtualbox/private_key
  IdentitiesOnly yes
  LogLevel FATAL`

func TestParseValidSshConfig(t *testing.T) {
	assert := assert.New(t)

	sshConfig, err := parseSshConfig(sshConfigString)
	assert.Nil(err)
	assert.NotNil(sshConfig)
	assert.Equal("default", sshConfig.Host)
	assert.Equal("2222", sshConfig.Options["Port"])
}
