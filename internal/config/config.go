package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Configure struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	PlusKey  string `yaml:"PlusKey"`
}

type Config struct {
	Mode    string `yaml:"Mode"`
	Name    string `yaml:"Name"`
	Host    string `yaml:"Host"`
	Port    int    `yaml:"Port"`
	Timeout int64  `yaml:"Timeout"`
	Configure
}

func MustLoad(path string, c *Config) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(data, c); err != nil {
		panic(err)
	}
}
