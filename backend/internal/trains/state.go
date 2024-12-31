package trains

import (
	"fmt"
	"log"
	"strings"
)

type TrainState int

const (
	Unknown     TrainState = iota // Default state when location can't be determined
	AtStation                     // Train is stopped at a station
	AtPlatform                    // Train is stopped at the platform
	Between                       // Train is between stations
	Approaching                   // Train is approaching next station
	Departed                      // Train has just left a station
)

var stateStrings = map[TrainState]string{
	Unknown:     "UNKNOWN",
	AtStation:   "AT_STATION",
	AtPlatform:  "AT_PLATFORM",
	Between:     "BETWEEN",
	Approaching: "APPROACHING",
	Departed:    "DEPARTED",
}

func (s TrainState) String() string {
	if str, ok := stateStrings[s]; ok {
		return str
	}
	return fmt.Sprintf("INVALID_STATE(%d)", s)
}

func (s TrainState) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s.String())), nil
}

func (s *TrainState) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	for state, stateStr := range stateStrings {
		if stateStr == str {
			*s = state
			return nil
		}
	}
	log.Printf("failed to unmarshal train state: %s", str)
	return fmt.Errorf("invalid train state: %s", str)
}

// DetectState determines the train state from a location description
// Returns Unknown if the location format is unexpected
func DetectState(location string) TrainState {
	if location == "" {
		return Unknown
	}
	location = strings.TrimSpace(location)

	switch {
	case strings.HasPrefix(location, "At "):
		return AtStation
	case strings.HasPrefix(location, "Platform"):
		return AtPlatform
	case strings.HasPrefix(location, "Between "):
		return Between
	case strings.HasPrefix(location, "Approaching "):
		return Approaching
	case strings.HasPrefix(location, "Left "):
		return Departed
	default:
		return Unknown
	}
}
