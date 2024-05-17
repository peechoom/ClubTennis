package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`

	Database struct {
		User     string `json:"user"`
		Password string `json:"pass"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
	} `json:"database"`

	Redis struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Pass string `json:"pass"`
	} `json:"redis"`

	GoogleOauth struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"google_oauth"`
}

func LoadConfig(filename string) (*Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
