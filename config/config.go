package config

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
	Credentials Credentials `yaml:"credentials"`
	Token       Token       `yaml:"token"`
}

type Credentials struct {
	ClientId            string   `yaml:"clientId"`
	ProjectId           string   `yaml:"projectId"`
	AuthUri             string   `yaml:"authUri"`
	TokenUri            string   `yaml:"tokenUri"`
	AuthProviderCertUrl string   `yaml:"authProviderCertUrl"`
	ClientSecret        string   `yaml:"clientSecret"`
	RedirectUris        []string `yaml:"redirectUris"`
}

type Token struct {
	AccessToken  string `yaml:"accessToken"`
	TokenType    string `yaml:"tokenType"`
	RefreshToken string `yaml:"refreshToken"`
	Expiry       string `yaml:"expiry"`
}

func ReadConfig() (*Config, error) {
	var cfg *Config
	path := configPath()
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		if err := createDefaultConfig(path); err != nil {
			return cfg, err
		}
		file, err = os.Open(path)
	}
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	return cfg, yaml.NewDecoder(file).Decode(cfg)
}

func configPath() string {
	path, exists := os.LookupEnv("CONFIG_PATH")
	if exists {
		return path
	}
	return "config.yml"
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
		Credentials: Credentials{
			ClientId:            "",
			ProjectId:           "birthday-notification",
			AuthUri:             "https://accounts.google.com/o/oauth2/auth",
			TokenUri:            "https://oauth2.googleapis.com/token",
			AuthProviderCertUrl: "https://www.googleapis.com/oauth2/v1/certs",
			ClientSecret:        "",
			RedirectUris: []string{
				"http://localhost",
			},
		},
		Token: Token{
			AccessToken:  "",
			TokenType:    "Bearer",
			RefreshToken: "",
			Expiry:       "2024-07-31T15:07:33.6502681+02:00",
		},
	},
}

func createDefaultConfig(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	return encoder.Encode(defaultConfig)
}
