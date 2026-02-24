package gold

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SilvDataProvider struct {
	client *http.Client
}

func NewSilvDataProvider(client *http.Client) *SilvDataProvider {
	return &SilvDataProvider{client: client}
}

func (p *SilvDataProvider) Name() string {
	return "SILV.DATA"
}

type silvDataResp struct {
	Commodities map[string]struct {
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		Unit        string  `json:"unit"`
		LastUpdated string  `json:"last_updated"`
		Timestamp   string  `json:"timestamp"`
	} `json:"commodities"`
}

func (p *SilvDataProvider) GetGoldUSDPerOunce(ctx context.Context) (float64, int64, error) {
	url := "https://data.silv.app/commodities.json"

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
		return 0, 0, fmt.Errorf("silvdata status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var payload silvDataResp
	if err := dec.Decode(&payload); err != nil {
		return 0, 0, err
	}

	gold, ok := payload.Commodities["gold"]
	if !ok {
		return 0, 0, fmt.Errorf("silvdata missing gold")
	}
	if gold.Price == 0 {
		return 0, 0, fmt.Errorf("silvdata missing price")
	}
	if gold.Currency != "USD" {
		return 0, 0, fmt.Errorf("silvdata unexpected currency %s", gold.Currency)
	}
	if gold.Unit != "troy_oz" {
		return 0, 0, fmt.Errorf("silvdata unexpected unit %s", gold.Unit)
	}

	var updatedAt int64
	if gold.LastUpdated != "" {
		if t, err := time.Parse(time.RFC3339, gold.LastUpdated); err == nil {
			updatedAt = t.Unix()
		}
	}
	if updatedAt == 0 && gold.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, gold.Timestamp); err == nil {
			updatedAt = t.Unix()
		}
	}

	return gold.Price, updatedAt, nil
}
