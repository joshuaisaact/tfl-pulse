package tfl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testPrediction() Prediction {
	return Prediction{
		VehicleID:       "123",
		StationName:     "Victoria",
		TimeToStation:   120,
		CurrentLocation: "At Platform",
		Towards:         "Walthamstow Central",
	}
}

func TestClient_GetVictoriaPredictions(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		errContains string
	}{
		{
			name: "happy path",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("User-Agent") == "" {
					t.Error("missing User-Agent header")
				}

				json.NewEncoder(w).Encode([]Prediction{testPrediction()})
			},
			wantErr: false,
		},
		{
			name: "server error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "something broke", http.StatusInternalServerError)
			},
			wantErr:     true,
			errContains: "500",
		},
		{
			name: "bad json",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{bad json]]`))
			},
			wantErr:     true,
			errContains: "parsing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			client := NewClient("test-key")
			client.baseURL = srv.URL + "/"

			preds, err := client.GetVictoriaPredictions(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got none")
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %v", tt.errContains, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(preds) == 0 {
				t.Error("expected predictions but got none")
			}
		})
	}
}
