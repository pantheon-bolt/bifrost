package application

import (
	"os"
	"strconv"
)

type Config struct {
	RedisAddress string
	ServerPort   int
}

func LoadConfig() Config {
	cfg := Config{
		RedisAddress: "localhost:6379",
		ServerPort:   3000,
	}

	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr
	}

	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		if port, err := strconv.Atoi(serverPort); err == nil {
			cfg.ServerPort = port
		}
	}
	return cfg
}
