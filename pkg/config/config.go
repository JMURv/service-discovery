package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type DB string

const (
	InMem  DB = "in-mem"
	SQLite DB = "sqlite"
)

type AcceptReq string

const (
	GRPC AcceptReq = "grpc"
	HTTP AcceptReq = "http"
)

type Config struct {
	DB        DB             `yaml:"db" env-default:"in-mem"`
	AcceptReq AcceptReq      `yaml:"accept-req" env-default:"grpc"`
	Server    *ServerConfig  `yaml:"server"`
	Checker   *CheckerConfig `yaml:"checker"`
}

type ServerConfig struct {
	Port   int    `yaml:"port" env-required:"true"`
	Mode   string `yaml:"mode" env-default:"dev"`
	Scheme string `yaml:"scheme" env-default:"http"`
	Domain string `yaml:"domain" env-default:"localhost"`
}

type CheckerConfig struct {
	Req           AcceptReq `yaml:"req" env-default:"grpc"`
	MaxRetriesReq int       `yaml:"max_retries_req" env-default:"3"`
	CooldownReq   int       `yaml:"cooldown_req" env-default:"5"`
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
