package market

import "context"

type RatesProvider interface {
	Name() string
	GetRates(ctx context.Context) (usdtry float64, eurtry float64, gbptry float64, jpytry float64, chftry float64, cadtry float64, audtry float64, nzdtry float64, sektry float64, noktry float64, dkktry float64, plntry float64, err error)
}

type GoldProvider interface {
	Name() string
	GetGoldUSDPerOunce(ctx context.Context) (price float64, updatedAtUnix int64, err error)
}
