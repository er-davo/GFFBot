package config

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotApiToken    string
	DatabaseURL            string
	DatabaseUser           string
	DatabasePassword       string
	DatabaseName           string
	DatabasePort           string
	DatabaseHost           string
	InactiveDaysDuration   int64  `json:"inactive_days_duration"`
	InactiveMinutsDuration int64  `json:"inactive_minuts_duration"`
	InactiveHoursDuration  int64  `json:"inactive_hours_duration"`
	CheckTimeDatabase      int64  `json:"check_time_database"`
	CheckTimeMemory        int64  `json:"check_time_memory"`
	SemaphoreTickets       int    `json:"semaphore_tickets"`
	LogFilePath            string `json:"log_file_path"`
}

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load("configs/.env")
		if err != nil {
			err = godotenv.Load("../../configs/.env")
			if err != nil {
				panic(err)
			}
		}

		config.TelegramBotApiToken = loadEnvStr("BOT_API_TOKEN")
		config.DatabaseURL = loadEnvStr("DATABASE_URL")
		config.DatabaseUser = loadEnvStr("POSTGRES_USER")
		config.DatabasePassword = loadEnvStr("POSTGRES_PASSWORD")
		config.DatabaseName = loadEnvStr("POSTGRES_DB")
		config.DatabasePort = loadEnvStr("POSTGRES_PORT")
		config.DatabaseHost = loadEnvStr("DB_HOST")

		jsonFile, err := os.ReadFile("configs/config.json")
		if err != nil {
			jsonFile, err = os.ReadFile("../../configs/config.json")
			if err != nil {
				panic(err)
			}
		}

		err = json.Unmarshal(jsonFile, &config)
		if err != nil {
			panic(err)
		}
	})
	return &config
}

func loadEnvStr(key string) string {
	return os.Getenv(key)
}
