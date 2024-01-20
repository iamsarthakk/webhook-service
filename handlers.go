package main

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	var payload Payload
	body, err := io.ReadAll(r.Body)

	err = json.Unmarshal(body, &payload)
	if err != nil {
		logger.Error("Failed to decode JSON payload", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payloads = append(payloads, payload)

	w.Write([]byte("OK"))
	logger.Info("Step1", zap.Int("BatchSize", len(payloads)))
}

func sendBatchHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Data Received"))
	logger.Info("Batch Received Successfully")
}
