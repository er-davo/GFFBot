package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../env/.env")
	if err != nil {
		log.Panic(err)
	}
}

type Config struct {
	TelegramBotApiToken string
	PathToLogFile       string
}

func Get() *Config {
	return &Config{
		TelegramBotApiToken: getEnvStr("BOT_API_TOKEN"),
	}
}

func getEnvStr(key string) string {
	return os.Getenv(key)
}