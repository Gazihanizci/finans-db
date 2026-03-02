package market

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Snapshot struct {
	USDTRY        float64
	EURTRY        float64
	GBPTRY        float64
	JPYTRY        float64
	CHFTRY        float64
	CADTRY        float64
	AUDTRY        float64
	NZDTRY        float64
	SEKTRY        float64
	NOKTRY        float64
	DKKTRY        float64
	PLNTRY        float64
	GramAltinTRY  float64
	Source        string
	FetchedAtUnix int64
}

type Service struct {
	ratesProviders []RatesProvider
	goldProviders  []GoldProvider
	ttl            time.Duration
	cache          *cache
}

func NewService(ttl time.Duration, rates []RatesProvider, gold []GoldProvider) *Service {
	return &Service{
		ratesProviders: rates,
		goldProviders:  gold,
		ttl:            ttl,
		cache:          newCache(ttl),
	}
}

func (s *Service) GetLatest(ctx context.Context) (Snapshot, error) {
	return s.cache.get(ctx, s.fetch)
}

func (s *Service) fetch(ctx context.Context) (Snapshot, error) {
	if len(s.ratesProviders) == 0 {
		return Snapshot{}, errors.New("no rates providers configured")
	}
	if len(s.goldProviders) == 0 {
		return Snapshot{}, errors.New("no gold providers configured")
	}

	usdtry, eurtry, gbptry, jpytry, chftry, cadtry, audtry, nzdtry, sektry, noktry, dkktry, plntry, rateSource, err := s.fetchRates(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	goldUSDPerOunce, goldSource, err := s.fetchGold(ctx)
	if err != nil {
		return Snapshot{}, err
	}

	usdPerGram := goldUSDPerOunce / 31.1034768
	gramAltinTRY := usdPerGram * usdtry

	fetchedAt := time.Now().Unix()
	source := fmt.Sprintf("fx:%s, gold:%s", rateSource, goldSource)

	return Snapshot{
		USDTRY:        usdtry,
		EURTRY:        eurtry,
		GBPTRY:        gbptry,
		JPYTRY:        jpytry,
		CHFTRY:        chftry,
		CADTRY:        cadtry,
		AUDTRY:        audtry,
		NZDTRY:        nzdtry,
		SEKTRY:        sektry,
		NOKTRY:        noktry,
		DKKTRY:        dkktry,
		PLNTRY:        plntry,
		GramAltinTRY:  gramAltinTRY,
		Source:        source,
		FetchedAtUnix: fetchedAt,
	}, nil
}

func (s *Service) fetchRates(ctx context.Context) (usdtry float64, eurtry float64, gbptry float64, jpytry float64, chftry float64, cadtry float64, audtry float64, nzdtry float64, sektry float64, noktry float64, dkktry float64, plntry float64, source string, err error) {
	var lastErr error
	for _, p := range s.ratesProviders {
		u, e, g, j, c, a, au, nz, se, no, dk, pl, eErr := p.GetRates(ctx)
		if eErr == nil {
			return u, e, g, j, c, a, au, nz, se, no, dk, pl, p.Name(), nil
		}
		lastErr = eErr
	}
	if lastErr == nil {
		lastErr = errors.New("rates providers failed")
	}
	return 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, "", lastErr
}

func (s *Service) fetchGold(ctx context.Context) (usdPerOunce float64, source string, err error) {
	var lastErr error
	for _, p := range s.goldProviders {
		price, _, eErr := p.GetGoldUSDPerOunce(ctx)
		if eErr == nil {
			return price, p.Name(), nil
		}
		lastErr = eErr
	}
	if lastErr == nil {
		lastErr = errors.New("gold providers failed")
	}
	return 0, "", lastErr
}
