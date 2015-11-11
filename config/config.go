package config

import (
	"io/ioutil"

	"log"

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

func readConfig(configPath string) *Config {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Can't open config file %s", configPath)
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
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
