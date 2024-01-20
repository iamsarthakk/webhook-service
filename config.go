package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"time"
)

func loadConfig() {
	err := godotenv.Load("config.env")
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}

	batchSize = getEnvAsInt("BATCH_SIZE", 10)
	batchInterval = time.Duration(getEnvAsInt("BATCH_INTERVAL", 60)) * time.Second
	postEndpoint = getEnv("POST_ENDPOINT", "http://localhost:8080/log_data")
}
