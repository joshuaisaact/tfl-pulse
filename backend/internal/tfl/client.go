package tfl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Prediction struct {
	VehicleID       string    `json:"vehicleId"`
	StationName     string    `json:"stationName"`
	PlatformName    string    `json:"platformName"`
	TimeToStation   int       `json:"timeToStation"` // seconds
	CurrentLocation string    `json:"currentLocation"`
	Towards         string    `json:"towards"` // destination
	Timestamp       time.Time `json:"timestamp"`
}

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.tfl.gov.uk/Line/victoria/Arrivals/",
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetVictoriaPredictions(ctx context.Context) ([]Prediction, error) {
	var predictions []Prediction
	url := fmt.Sprintf("%s?app_key=%s", c.baseURL, c.apiKey)
	fmt.Println("Requesting URL:", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "tfl-pulse/1.0 (https://github.com/joshuaisaact/tfl-pulse)") // Critical to avoid TFL 403 errors - they hate requests without user agents!
	req.Header.Set("Accept", "*/*")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response status: %d\n", resp.StatusCode)
		fmt.Printf("Response body: %s\n", string(body))
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&predictions); err != nil {
		return nil, fmt.Errorf("error parsing the predictions: %w", err)
	}

	return predictions, nil
}
