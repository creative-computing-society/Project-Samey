package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, loading environment variables from system")
	}
}

func GetSudoPassword() string {
	return os.Getenv("SUDO_PASSWORD")
}
