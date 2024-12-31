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
	Station     string
	IsBetween   bool
	PrevStation string
}

type TrainInfo struct {
	Location     Location
	Direction    Direction
	DirectionStr string
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
		// stations := strings.Split(loc, "At ")
		return Location{
			Station:     loc,
			IsBetween:   false,
			PrevStation: "",
		}
	}
}

func ProcessPredictions(predictions []tfl.Prediction) TrainMap {
	// TODO: Transform predictions into simple TrainMap MVP
	trainMap := map[string]TrainInfo{}
	for _, p := range predictions {
		if _, ok := trainMap[p.VehicleID]; !ok {
			direction := determineDirection(p.Towards)
			trainMap[p.VehicleID] = TrainInfo{
				Location:     parseLocation(p.CurrentLocation),
				Direction:    direction,
				DirectionStr: direction.String(),
			}
		}
	}
	return trainMap
}
