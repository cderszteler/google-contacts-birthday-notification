package main

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Mail      Mail      `yaml:"mail"`
	PeopleApi PeopleApi `yaml:"peopleApi"`
}

type Mail struct {
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
	Receiver string `yaml:"receiver"`
	Tls      bool   `yaml:"tls"`
	Secure   bool   `yaml:"secure"`
}

type PeopleApi struct {
	AccessToken  string `yaml:"accessToken"`
	TokenType    string `yaml:"tokenType"`
	RefreshToken string `yaml:"refreshToken"`
	Expiry       string `yaml:"expiry"`
}

const path = "config.yml"

func ReadConfig(cfg *Config) error {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		if err := createDefaultConfig(); err != nil {
			return err
		}
		file, err = os.Open(path)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(cfg)
}

var defaultConfig = Config{
	Mail: Mail{
		Port:     465,
		Host:     "smpt.mail.com",
		User:     "form",
		Password: "HNJ38a",
		Sender:   "form@mail.com",
		Receiver: "admin@mail.com",
		Tls:      false,
		Secure:   false,
	},
	PeopleApi: PeopleApi{
		AccessToken:  "",
		TokenType:    "Bearer",
		RefreshToken: "",
		Expiry:       "2024-07-31T15:07:33.6502681+02:00",
	},
}

func createDefaultConfig() error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(defaultConfig)
}
