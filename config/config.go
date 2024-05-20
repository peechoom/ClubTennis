package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig(filename string) error {
	cloud := os.Getenv("CLOUD_INSTANCE")
	if len(cloud) == 0 || cloud == "false" {
		return godotenv.Load(filename)
	}
	return nil
}
