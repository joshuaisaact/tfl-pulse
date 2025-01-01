package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/poller"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/trains"
)

func addRoutes(
	mux *http.ServeMux,
	client *tfl.Client,
) {
	// Create hub and poller
	poller := poller.New(client)

	hub := NewHub(poller)

	poller.SetUpdateCallback(hub.broadcastTrains)

	// Start polling in background
	go poller.Start(context.Background())

	// Websocket endpoint
	mux.HandleFunc("/ws", hub.handleWebSocket)

	mux.HandleFunc("/api/victoria", handleVictoria(client))
	mux.HandleFunc("/api/trains", handleTrains(client))
	mux.Handle("/", http.NotFoundHandler())
}

func handleVictoria(client *tfl.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		predictions, err := client.GetVictoriaPredictions(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get Victoria Line predictions: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(predictions)
	}
}

func handleTrains(client *tfl.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		predictions, err := client.GetVictoriaPredictions(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get Victoria Line predictions: %v", err), http.StatusInternalServerError)
			return
		}

		trainMap := trains.ProcessPredictions(predictions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(trainMap)
	}
}
