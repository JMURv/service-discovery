package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ServiceName string        `yaml:"serviceName" env-required:"true"`
	Server      *ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port   int    `yaml:"port" env-required:"true"`
	Mode   string `yaml:"mode" env-default:"dev"`
	Scheme string `yaml:"scheme" env-default:"http"`
	Domain string `yaml:"domain" env-default:"localhost"`
}

func MustLoad(configPath string) *Config {
	var conf Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		panic("failed to unmarshal config: " + err.Error())
	}

	return &conf
}
