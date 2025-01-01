package trains

import (
	"testing"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
)

func TestProcessPredictions(t *testing.T) {
	predictions := []tfl.Prediction{
		{
			VehicleID:       "123",
			StationName:     "Victoria Underground Station",
			CurrentLocation: "At Victoria",
			Towards:         "Walthamstow Central",
			TimeToStation:   120,
		},
		{
			VehicleID:       "456",
			StationName:     "Oxford Circus Underground Station",
			CurrentLocation: "Between Warren Street and Oxford Circus",
			Towards:         "Brixton",
			TimeToStation:   240,
		},
	}
	trainMap := ProcessPredictions(predictions)

	if len(trainMap) != 2 {
		t.Errorf("Expected 2 trains, got %d", len(trainMap))
	}

	if train, ok := trainMap["123"]; ok {
		// Test direction (raw towards value)
		if train.Direction != "Walthamstow Central" {
			t.Errorf("Expected direction Walthamstow Central, got %s", train.Direction)
		}
		// Test station
		if train.Location.StationID != "Victoria Underground Station" {
			t.Errorf("Expected station Victoria Underground Station, got %s", train.Location.StationID)
		}
		// Test state
		if train.Location.IsBetween {
			t.Error("Train 123 should not be between stations")
		}
		// Test time
		if train.TimeToNext != 120 {
			t.Errorf("Expected time to next 120, got %d", train.TimeToNext)
		}
	} else {
		t.Error("Train 123 not found in map")
	}

	if train, ok := trainMap["456"]; ok {
		// Test direction (raw towards value)
		if train.Direction != "Brixton" {
			t.Errorf("Expected direction Brixton, got %s", train.Direction)
		}
		// Test station
		if train.Location.StationID != "Oxford Circus" {
			t.Errorf("Expected next station Oxford Circus, got %s", train.Location.StationID)
		}
		// Test between state
		if !train.Location.IsBetween {
			t.Error("Train 456 should be between stations")
		}
		// Test previous station
		if train.Location.PrevStationID != "Warren Street" {
			t.Errorf("Expected previous station Warren Street, got %s", train.Location.PrevStationID)
		}
		// Test time
		if train.TimeToNext != 240 {
			t.Errorf("Expected time to next 240, got %d", train.TimeToNext)
		}
	} else {
		t.Error("Train 456 not found in map")
	}
}
