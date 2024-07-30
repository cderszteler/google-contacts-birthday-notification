package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConfig(cfg *Config) error {
	file, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	return yaml.NewDecoder(file).Decode(cfg)
}

type Config struct {
	Mail struct {
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Receiver string `yaml:"receiver"`
		Secure   bool   `yaml:"secure"`
		Tls      bool   `yaml:"tls"`
	} `yaml:"mail"`
}
