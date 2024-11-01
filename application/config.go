package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress      string
	PostgreSQLAddress string
	ServerPort        int
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress:      "localhost:6379",
		PostgreSQLAddress: "postgres://postgres:postgres@db:5432/bolt?sslmode=disable",
		ServerPort:        3000,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if postgresAddr, exists := os.LookupEnv("POSTGRES_ADDR"); exists {
		cfg.PostgreSQLAddress = postgresAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.Atoi(serverPort); err == nil {
			cfg.ServerPort = port
		}
	}
	return cfg
}
