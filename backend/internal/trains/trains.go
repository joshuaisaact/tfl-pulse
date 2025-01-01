package trains

import (
	"strings"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
)

type Direction bool

const (
	Northbound Direction = true
	Southbound Direction = false
)

func (d Direction) String() string {
	if d == Northbound {
		return "Northbound"
	} else {
		return "Southbound"
	}
}

type Location struct {
	StationID     string
	IsBetween     bool
	PrevStationID string
	State         TrainState
}

type TrainInfo struct {
	Location   Location
	Direction  string
	TimeToNext int
}

type TrainMap map[string]TrainInfo

func extractStationsFromLocation(text string) (string, string) {
	switch {
	case strings.HasPrefix(text, "Between "):
		parts := strings.Split(strings.TrimPrefix(text, "Between "), " and ")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		}
	case strings.HasPrefix(text, "Left "), strings.HasPrefix(text, "Departed "):
		// For Left/Departed, the station mentioned is the previous station
		station := ""
		if strings.HasPrefix(text, "Left ") {
			station = strings.TrimPrefix(text, "Left ")
		} else {
			station = strings.TrimPrefix(text, "Departed ")
		}
		station = strings.TrimSpace(station)
		return station, ""
	case strings.HasPrefix(text, "Approaching "):
		// For Approaching, the station mentioned is the next station
		station := strings.TrimSpace(strings.TrimPrefix(text, "Approaching "))
		return "", station
	}
	return "", ""
}

func parseLocation(p tfl.Prediction) Location {
	state := DetectState(p.CurrentLocation)
	loc := Location{
		StationID: p.StationName,
		State:     state,
		IsBetween: state == Between || state == Left || state == Departed,
	}

	// Get previous and next stations based on state
	prev, next := extractStationsFromLocation(p.CurrentLocation)

	switch state {
	case Between:
		loc.PrevStationID = prev
		loc.StationID = next // Update current to next station
	case Left, Departed:
		loc.PrevStationID = prev // The station we just left
		loc.StationID = prev     // Current position is still closest to prev
	case Approaching:
		loc.StationID = next // Update to station we're approaching
	case AtStation, AtPlatform:
		// Keep stationID from prediction
	}

	return loc
}

func ProcessPredictions(predictions []tfl.Prediction) TrainMap {
	trains := make(TrainMap)

	// First, group predictions by vehicle ID to find minimum time
	predictionsByVehicle := make(map[string][]tfl.Prediction)
	for _, p := range predictions {
		predictionsByVehicle[p.VehicleID] = append(predictionsByVehicle[p.VehicleID], p)
	}

	// Then process each vehicle's predictions
	for vehicleID, vehiclePredictions := range predictionsByVehicle {
		// Find prediction with minimum TimeToStation
		minTimeIdx := 0
		for i, p := range vehiclePredictions {
			if p.TimeToStation < vehiclePredictions[minTimeIdx].TimeToStation {
				minTimeIdx = i
			}
		}

		// Use the prediction with minimum time
		p := vehiclePredictions[minTimeIdx]
		trains[vehicleID] = TrainInfo{
			Location:   parseLocation(p),
			Direction:  p.Towards,
			TimeToNext: p.TimeToStation,
		}
	}

	return trains
}
