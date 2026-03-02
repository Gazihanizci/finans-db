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

func (p *MoneyConvertProvider) GetRates(ctx context.Context) (float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, float64, error) {
	url := "https://cdn.moneyconvert.net/api/latest.json"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert status %d", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var payload moneyConvertResp
	if err := dec.Decode(&payload); err != nil {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, err
	}

	tryRate, ok := payload.Rates["TRY"]
	if !ok {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing TRY")
	}
	eurRate, ok := payload.Rates["EUR"]
	if !ok || eurRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing EUR")
	}
	gbpRate, ok := payload.Rates["GBP"]
	if !ok || gbpRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing GBP")
	}
	jpyRate, ok := payload.Rates["JPY"]
	if !ok || jpyRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing JPY")
	}
	chfRate, ok := payload.Rates["CHF"]
	if !ok || chfRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing CHF")
	}
	cadRate, ok := payload.Rates["CAD"]
	if !ok || cadRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing CAD")
	}
	audRate, ok := payload.Rates["AUD"]
	if !ok || audRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing AUD")
	}
	nzdRate, ok := payload.Rates["NZD"]
	if !ok || nzdRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing NZD")
	}
	sekRate, ok := payload.Rates["SEK"]
	if !ok || sekRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing SEK")
	}
	nokRate, ok := payload.Rates["NOK"]
	if !ok || nokRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing NOK")
	}
	dkkRate, ok := payload.Rates["DKK"]
	if !ok || dkkRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing DKK")
	}
	plnRate, ok := payload.Rates["PLN"]
	if !ok || plnRate == 0 {
		return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, fmt.Errorf("moneyconvert missing PLN")
	}

	usdtry := tryRate
	eurtry := tryRate / eurRate
	gbptry := tryRate / gbpRate
	jpytry := tryRate / jpyRate
	chftry := tryRate / chfRate
	cadtry := tryRate / cadRate
	audtry := tryRate / audRate
	nzdtry := tryRate / nzdRate
	sektry := tryRate / sekRate
	noktry := tryRate / nokRate
	dkktry := tryRate / dkkRate
	plntry := tryRate / plnRate
	return usdtry, eurtry, gbptry, jpytry, chftry, cadtry, audtry, nzdtry, sektry, noktry, dkktry, plntry, nil
}
