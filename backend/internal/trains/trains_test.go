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
		if !train.Direction {
			t.Errorf("Train 123 should be Northbound")
		}
	}
}
