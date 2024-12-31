package trains

import (
	"testing"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
)

func TestProcessPredictions(t *testing.T) {
	predictions := []tfl.Prediction{
		{
			VehicleID:       "123",
			CurrentLocation: "At Victoria",
			Towards:         "Walthamstow Central",
		},
		{
			VehicleID:       "456",
			CurrentLocation: "Between Warren Street and Oxford Circus",
			Towards:         "Brixton",
		},
	}
	trainMap := ProcessPredictions(predictions)

	if len(trainMap) != 2 {
		t.Errorf("Expected 2 trains, got %d", len(trainMap))
	}

	if train, ok := trainMap["123"]; ok {
		if train.DirectionStr != "Northbound" {
			t.Errorf("Train 123 should be Northbound")
		}
		if train.Location.Station != "Victoria" {
			t.Errorf("Train 123 should be at Victoria, got %s", train.Location.Station)
		}
		if train.Location.IsBetween {
			t.Errorf("Train 123 should not be between stations")
		}
	} else {
		t.Error("Train 123 not found in map")
	}

	if train, ok := trainMap["456"]; ok {
		if train.DirectionStr != "Southbound" {
			t.Errorf("Train 456 should be Southbound")
		}
		if train.Location.Station != "Oxford Circus" {
			t.Errorf("Train 456 should be at Oxford Circus, got %s", train.Location.Station)
		}
		if !train.Location.IsBetween {
			t.Error("Train 456 should be between stations")
		}
		if train.Location.PrevStation != "Warren Street" {
			t.Errorf("Train previous station should be Warren Street, got %s", train.Location.PrevStation)
		}
	} else {
		t.Error("Train 456 not found in map")
	}
}
