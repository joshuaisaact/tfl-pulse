package poller

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/trains"
)

type Poller struct {
	client    tfl.Client
	trainData map[string]trains.TrainInfo
	mu        sync.RWMutex
}

func New(client *tfl.Client) *Poller {
	return &Poller{
		client:    client,
		trainData: make(map[string]trains.TrainInfo),
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Do an initial poll immediately
	p.poll()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.poll()
		}
	}
}

func (p *Poller) poll() {
	predictions, err := p.client.GetVictoriaPredictions(context.Background())
	if err != nil {
		log.Printf("Error getting predictions: %v", err)
		return
	}

	p.mu.Lock()
	p.trainData = trains.ProcessPredictions(predictions)
	p.mu.Unlock()
}

func (p *Poller) GetTrains() map[string]trains.TrainInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]trains.TrainInfo)
	for k, v := range p.trainData {
		result[k] = v
	}
	return result
}
