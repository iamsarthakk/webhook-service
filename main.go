package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
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
	loadConfig()

	logger, _ = zap.NewProduction()
	defer logger.Sync()
}

func main() {
	r := chi.NewRouter()

	setupMiddleware(r)

	r.Get("/healthz", healthzHandler)
	r.Post("/log", logHandler)
	r.Post("/send-data", sendBatchHandler)

	logger.Info("Webhook receiver started")

	startBatchSender()

	http.ListenAndServe(":8080", r)
}

func startBatchSender() {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Info("Step2: Interval Reached", zap.Int("BatchSize", len(payloads)))
				sendBatch(payloads)
				payloads = make([]Payload, 0, batchSize)

			case <-time.After(0): // Non-blocking check
				if len(payloads) >= batchSize {
					logger.Info("Step2: Batch Size Reached", zap.Int("BatchSize", len(payloads)))
					sendBatch(payloads)
					payloads = make([]Payload, 0, batchSize)
				}
			}
		}
	}()
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
