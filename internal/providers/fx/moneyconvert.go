package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type MoneyConvertProvider struct {
	client *http.Client
}

func NewMoneyConvertProvider(client *http.Client) *MoneyConvertProvider {
	return &MoneyConvertProvider{client: client}
}

func (p *MoneyConvertProvider) Name() string {
	return "MoneyConvert"
}

type moneyConvertResp struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func (p *MoneyConvertProvider) GetRates(ctx context.Context) (float64, float64, float64, error) {
	url := "https://cdn.moneyconvert.net/api/latest.json"

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
		return 0, 0, 0, fmt.Errorf("moneyconvert status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var payload moneyConvertResp
	if err := dec.Decode(&payload); err != nil {
		return 0, 0, 0, err
	}

	tryRate, ok := payload.Rates["TRY"]
	if !ok {
		return 0, 0, 0, fmt.Errorf("moneyconvert missing TRY")
	}
	eurRate, ok := payload.Rates["EUR"]
	if !ok || eurRate == 0 {
		return 0, 0, 0, fmt.Errorf("moneyconvert missing EUR")
	}
	gbpRate, ok := payload.Rates["GBP"]
	if !ok || gbpRate == 0 {
		return 0, 0, 0, fmt.Errorf("moneyconvert missing GBP")
	}

	usdtry := tryRate
	eurtry := tryRate / eurRate
	gbptry := tryRate / gbpRate
	return usdtry, eurtry, gbptry, nil
}
