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
	} else if strings.Contains(loc, "At") {
		stations := strings.Split(loc, "At ")
		return Location{
			Station:     stations[1],
			IsBetween:   false,
			PrevStation: "",
		}
	} else if strings.Contains(loc, "Approaching") {
		stations := strings.Split(loc, "Approaching ")
		return Location{
			Station:     stations[1],
			IsBetween:   false,
			PrevStation: "",
		}
	} else {
		return Location{
			Station:     loc,
			IsBetween:   false,
			PrevStation: "",
		}
	}
}

// func parseLocation(loc string):
//     // First detect the "state" of the train
//     state := detectState(loc)  // Could return enum/string like AT, BETWEEN, APPROACHING, DEPARTED, etc

//     switch state:
//         case "BETWEEN":
//             // Current approach works fine
//             parts := split "Between "
//             stations := split " and "
//             return Location{stations[1], true, stations[0]}

//         case "AT_PLATFORM":
//             // Need to get station from prediction.stationName instead
//             // Because "At Platform" doesn't tell us which station
//             return Location{stationName, false, ""}

//         case "AT":
//             if strings.Contains(loc, "Platform"):
//                 // Handle "At Platform 5" case
//                 return Location{stationName, false, ""}
//             else:
//                 // Handle "At Victoria" case
//                 station := split "At "
//                 return Location{station, false, ""}

//         case "DEPARTED", "LEFT":
//             // Just left a station, could track this state
//             station := split "Departed/Left "
//             return Location{station, true, station}  // Maybe mark as between?

//         case "APPROACHING":
//             // Almost at next station
//             station := split "Approaching "
//             return Location{station, false, ""}

//         default:
//             // Unknown state, just use the raw location
//             return Location{loc, false, ""}

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
