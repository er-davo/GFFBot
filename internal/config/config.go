package config

import (
	"log"
	"os"
	"sync"

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

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config.TelegramBotApiToken = loadEnvStr("BOT_API_TOKEN")	
	})
	return &config
}

func loadEnvStr(key string) string {
	return os.Getenv(key)
}
