package poller

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joshuaisaact/tfl-pulse/backend/internal/tfl"
	"github.com/joshuaisaact/tfl-pulse/backend/internal/trains"
)

type Poller struct {
	client    *tfl.Client
	trainData map[string]trains.TrainInfo
	mu        sync.RWMutex
	onUpdate  func()
}

func New(client *tfl.Client) *Poller {
	return &Poller{
		client:    client,
		trainData: make(map[string]trains.TrainInfo),
	}
}

func (p *Poller) SetUpdateCallback(callback func()) {
	p.onUpdate = callback
}

func (p *Poller) Start(ctx context.Context) {
	// Poll every 6 seconds (10 times per minute)
	// This gives us plenty of headroom under the 500 requests/minute limit
	ticker := time.NewTicker(6 * time.Second)
	defer ticker.Stop()

	// Log polling starts
	log.Printf("Starting poller with 6-second interval")

	// Do an initial poll immediately
	if err := p.poll(); err != nil {
		log.Printf("Error in initial poll: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Poller stopping due to context cancellation")
			return
		case <-ticker.C:
			if err := p.poll(); err != nil {
				log.Printf("Error polling: %v", err)
				// Continue running despite error
			}
		}
	}
}

func (p *Poller) poll() error {
	predictions, err := p.client.GetVictoriaPredictions(context.Background())
	if err != nil {
		return fmt.Errorf("error getting predictions: %w", err)
	}

	p.mu.Lock()
	p.trainData = trains.ProcessPredictions(predictions)
	p.mu.Unlock()

	if p.onUpdate != nil {
		p.onUpdate()
	}

	return nil
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
