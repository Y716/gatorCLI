package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFileData, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer configFileData.Close()

	encoder := json.NewEncoder(configFileData)
	return encoder.Encode(cfg)
}

func Read() (Config, error) {
	cfgFilePath, _ := getConfigFilePath()

	cfgFileData, err := os.Open(cfgFilePath)
	if err != nil {
		return Config{}, err
	}
	defer cfgFileData.Close()

	var cfg Config
	decoder := json.NewDecoder(cfgFileData)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) SetUser(newName string) error {
	c.CurrentUserName = newName
	return write(*c)
}
