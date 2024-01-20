package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

type Payload struct {
	UserID    int     `json:"user_id"`
	Total     float64 `json:"total"`
	Title     string  `json:"title"`
	Meta      Meta    `json:"meta"`
	Completed bool    `json:"completed"`
}

type Meta struct {
	Logins       []Login `json:"logins"`
	PhoneNumbers struct {
		Home   string `json:"home"`
		Mobile string `json:"mobile"`
	} `json:"phone_numbers"`
}

type Login struct {
	Time time.Time `json:"time"`
	IP   string    `json:"ip"`
}

var (
	batchSize     int
	batchInterval time.Duration
	postEndpoint  string
	logger        *zap.Logger
	payloads      []Payload
)

func init() {
	batchSize = getEnvAsInt("BATCH_SIZE", 10)
	batchInterval = time.Duration(getEnvAsInt("BATCH_INTERVAL", 60)) * time.Second
	postEndpoint = getEnv("POST_ENDPOINT", "http://localhost:8080/log")

	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", healthzHandler)
	r.Post("/log", logHandler)

	logger.Info("Webhook receiver started")

	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	payloads := make([]Payload, 0, batchSize)

	go func() {
		for {
			select {
			case <-ticker.C:
				sendBatch(payloads)
				payloads = make([]Payload, 0, batchSize)
			}
		}
	}()

	http.ListenAndServe(":8080", r)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		logger.Error("Failed to decode JSON payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Store the payload in-memory
	payloads = append(payloads, payload)

	w.Write([]byte("OK"))
}

func sendBatch(payloads []Payload) {
	if len(payloads) == 0 {
		return
	}

	startTime := time.Now()
	payloadsJSON, _ := json.Marshal(payloads)

	resp, err := http.Post(postEndpoint, "application/json", bytes.NewBuffer(payloadsJSON))
	if err != nil {
		logger.Error("Failed to send batch", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	logger.Info("Batch sent",
		zap.Int("BatchSize", len(payloads)),
		zap.Int("StatusCode", resp.StatusCode),
		zap.Duration("Duration", time.Since(startTime)),
	)

	if resp.StatusCode != http.StatusOK {
		// Retry logic
		retryBatch(payloads)
	}
}

func retryBatch(payloads []Payload) {
	for i := 0; i < 3; i++ {
		time.Sleep(2 * time.Second)
		logger.Info("Retrying batch", zap.Int("RetryAttempt", i+1))
		sendBatch(payloads)
	}
	logger.Error("Failed to send batch after 3 attempts, exiting application")
	os.Exit(1)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}