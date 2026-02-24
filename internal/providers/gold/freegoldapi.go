package gold

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FreeGoldAPIProvider struct {
	client *http.Client
}

func NewFreeGoldAPIProvider(client *http.Client) *FreeGoldAPIProvider {
	return &FreeGoldAPIProvider{client: client}
}

func (p *FreeGoldAPIProvider) Name() string {
	return "FreeGoldAPI"
}

type freeGoldEntry struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
}

func (p *FreeGoldAPIProvider) GetGoldUSDPerOunce(ctx context.Context) (float64, int64, error) {
	url := "https://freegoldapi.com/data/latest.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, fmt.Errorf("freegoldapi status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var payload []freeGoldEntry
	if err := dec.Decode(&payload); err != nil {
		return 0, 0, err
	}
	if len(payload) == 0 {
		return 0, 0, fmt.Errorf("freegoldapi empty payload")
	}

	entry := payload[len(payload)-1]
	if entry.Price == 0 {
		return 0, 0, fmt.Errorf("freegoldapi missing price")
	}

	var updatedAt int64
	if entry.Date != "" {
		if t, err := time.Parse("2006-01-02", entry.Date); err == nil {
			updatedAt = t.Unix()
		}
	}

	return entry.Price, updatedAt, nil
}
