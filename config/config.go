package config

import (
	"io/ioutil"

	"log"

	"text/template"

	"bytes"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	DefaultConfigPath = ".anvil.yml"
	Cfg               *Config
)

type DriverConfig struct {
	Name    string
	Options map[string]interface{}
}

type ProvisionerConfig struct {
	Name    string
	Options map[string]interface{}
}

type PlatformConfig struct {
	Name   string
	Driver map[string]interface{}
}

type SuiteConfig struct {
	Name        string
	Provisioner map[string]interface{}
}

type Config struct {
	Driver      *DriverConfig
	Provisioner *ProvisionerConfig
	Platforms   []*PlatformConfig
	Suites      []*SuiteConfig
}

type configVars struct {
	Env map[string]string
}

func buildConfigVars() configVars {
	vars := configVars{
		Env: make(map[string]string),
	}
	for _, envVar := range os.Environ() {
		parts := strings.Split(envVar, "=")
		if len(parts) == 2 {
			vars.Env[parts[0]] = parts[1]
		}
	}
	return vars
}

func readConfig(configPath string) *Config {
	configTemplate := template.New("config")

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Can't open config file %s", configPath)
	}
	configTemplate, err = configTemplate.Parse(string(data))
	if err != nil {
		log.Fatalf("Can't parse config template: %v")
	}
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))
	err = configTemplate.Execute(buffer, buildConfigVars())
	if err != nil {
		log.Fatalf("Can't template config file: %v", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		log.Fatalf("Invalid config file %s", configPath)
	}
	return config
}

func LoadConfig(configPath string) *Config {
	if Cfg == nil {
		Cfg = readConfig(configPath)
	}
	return Cfg
}
