package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type FrankfurterProvider struct {
	client *http.Client
}

func NewFrankfurterProvider(client *http.Client) *FrankfurterProvider {
	return &FrankfurterProvider{client: client}
}

func (p *FrankfurterProvider) Name() string {
	return "Frankfurter"
}

type frankfurterResp struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (p *FrankfurterProvider) GetRates(ctx context.Context) (float64, float64, float64, error) {
	url := "https://api.frankfurter.dev/v1/latest?base=USD&symbols=TRY,EUR,GBP"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, 0, fmt.Errorf("frankfurter status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var payload frankfurterResp
	if err := dec.Decode(&payload); err != nil {
		return 0, 0, 0, err
	}

	tryRate, ok := payload.Rates["TRY"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("frankfurter missing TRY")
	}
	eurRate, ok := payload.Rates["EUR"]
	if !ok || eurRate == 0 {
		return 0, 0, 0, fmt.Errorf("frankfurter missing EUR")
	}
	gbpRate, ok := payload.Rates["GBP"]
	if !ok || gbpRate == 0 {
		return 0, 0, 0, fmt.Errorf("frankfurter missing GBP")
	}

	usdtry := tryRate
	eurtry := tryRate / eurRate
	gbptry := tryRate / gbpRate
	return usdtry, eurtry, gbptry, nil
}
