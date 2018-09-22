package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config is a config type
type Config struct {
	Version  string `yaml:"version"`
	Telegram struct {
		Token string `yaml:"token"`
	} `yaml:"telegram"`
	Trello struct {
		Token string `yaml:"token"`
	} `yaml:"trello"`
}

// Load a config from path
func Load(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
