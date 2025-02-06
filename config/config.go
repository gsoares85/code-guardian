package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment variables")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
