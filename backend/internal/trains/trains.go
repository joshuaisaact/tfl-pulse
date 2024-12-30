package trains

import "github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"

type Direction bool

const (
	Northbound Direction = true
	Southbound Direction = false
)

type Location struct {
	Station     string
	IsBetween   bool
	PrevStation bool
}

type TrainInfo struct {
	Location  Location
	Direction Direction
}

// This will map vehicle id to Train Info
type TrainMap map[string]TrainInfo

func ProcessPredictions(predictions []tfl.Prediction) TrainMap {
	// TODO: Transform predictions into simple TrainMap MVP
	return nil
}
