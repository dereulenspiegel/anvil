# anvil - Forge your Infrastructure [![Build Status](https://travis-ci.org/dereulenspiegel/anvil.svg)](https://travis-ci.org/dereulenspiegel/anvil)

Many of you may already know [Test Kitchen](http://kitchen.ci/). This very handy
tool allows you to create automated tests for your infrastructure automation. I
use it every chance I get to automatically test my ansible roles. While Test Kitchen
works well enough for most people I had some issues with it. The biggest one being
the implementation language.

So I created something new in Go called anvil. Anvil is extensible via plugins.
But you are not limited to Go when writing plugins. Anvil comminicates with its
plugins via JSON-RPC over Stdin and Stdout, so you can write plugins in any
language which has access to Stdin and Stdout and has a JSON-RPC library (or
you write your own JSON-RPC library).

## Features

* [x] Flexible driver plugin interface to create, start, stop and destroy instances
* [x] Flexible provisioner plugin interface to run provisioners like ansible
* [x] verifier plugin interface to verify instances
* [x] Driver plugin for vagrant
* [x] Provisioner plugin for ansible
* [x] Verifier plugin for testinfra
* [] Plugin interface for test report formatters (if this makes sense, tbd)
* [] Plugin interface for notifier (for these long running tests)
* [] Let multiple test suites run in parallel
* [] More plugins for more drivers, provisioners and verifiers
* [] Easier plugin installation and management

## Installation

To install anvil you need to you Go setup on your machine. When Go is setup you
can simply `go install github.com/dereulenspiegel/anvil`.

Plugins simply live as binaries in your path. Currently the following plugins are
available:

* Driver: vagrant
* Provisioner: ansible
* Verifier: testinfra

To install one of these plugins simply execute
`go install github.com/dereulenspiegel/anvil/anvil-[driver|provisioner|verifier]-{name}`.

## Usage

Anvil requires a configuration file in your infrastructure project which is very
similar to the configuration file required by kitchen. A simple example looks
like this:

```yaml
driver:
  name: vagrant   # Name of the driver plugin to use
  options:        # Optional parameters, driver specific

provisioner:
  name: ansible   # Name of the provisioner plugin to use
  options:        # Optional parameters, provisioner specific

platforms:                  # List of platforms to test on
  - name: platform1         # Platform name
    driver:                 # Driver and instance specific configuration
      box: ubuntu/trusty64
    #  url: http://somewhere

suites:                     # List of test suites to execute on all instances
  - name: Suite1            # Name of the test suite
    provisioner:            # Provisioner and test suite specific options
      playbook: playbook.yml
```

Your project also requires a subfolder called `tests`. In this subfolder each test
suite has its own subfolder and in each suite folder is a folder with the name
of the verifier plugin you want to use. You can use multiple verifiers per test suite.

If you have setup your configuration and tests you can simply execute
`anvil verify` to setup, provision and verify all platforms with all test suites.
You can limit the scope of the command by appending a regex describing all suites or
platforms you want to verify. `anvil help` should give you an overview over all
available commands.
