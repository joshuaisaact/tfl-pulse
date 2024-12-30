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

type Location struct {
	Station     string
	IsBetween   bool
	PrevStation string
}

type TrainInfo struct {
	Location  Location
	Direction Direction
}

// This will map vehicle id to Train Info
type TrainMap map[string]TrainInfo

func determineDirection(towards string) Direction {
	return towards == "Walthamstow Central"
}

func parseLocation(loc string) Location {
	if strings.Contains(loc, "Between") {
		parts := strings.Split(loc, "Between ")
		stations := strings.Split(parts[1], " and ")
		return Location{
			Station:     stations[1],
			IsBetween:   true,
			PrevStation: stations[0],
		}
	} else {
		stations := strings.Split(loc, "At ")
		return Location{
			Station:     stations[0],
			IsBetween:   false,
			PrevStation: "",
		}
	}
}

func determineBetween(between string) bool {
	return strings.Contains(between, "Between")
}

func ProcessPredictions(predictions []tfl.Prediction) TrainMap {
	// TODO: Transform predictions into simple TrainMap MVP
	trainMap := map[string]TrainInfo{}
	for _, p := range predictions {
		if _, ok := trainMap[p.VehicleID]; !ok {
			trainMap[p.VehicleID] = TrainInfo{
				Location:  parseLocation(p.CurrentLocation),
				Direction: determineDirection(p.Towards),
			}
		}
	}
	return trainMap
}
