package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"scraper/internal/httpapi"
	"scraper/internal/market"
	"scraper/internal/providers/fx"
	"scraper/internal/providers/gold"
	"scraper/internal/storage"
)

func Run() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ratesProviders := []market.RatesProvider{
		fx.NewFrankfurterProvider(client),
		fx.NewMoneyConvertProvider(client),
	}
	goldProviders := []market.GoldProvider{
		gold.NewSilvDataProvider(client),
		gold.NewFreeGoldAPIProvider(client),
	}

	service := market.NewService(1*time.Hour, ratesProviders, goldProviders)

	var store *storage.MySQLStore
	var lastPrinted int64
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		s, err := storage.NewMySQLStore(dsn)
		if err != nil {
			log.Printf("mysql baglanamadi: %v", err)
		} else {
			store = s
			log.Printf("mysql baglandi")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := store.EnsureSchema(ctx); err != nil {
				log.Printf("db schema uyarisi: %v", err)
			}
			cancel()
		}
	} else {
		log.Printf("MYSQL_DSN bos, DB guncellemesi kapali")
	}

	printSnapshot := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if snap, err := service.GetLatest(ctx); err == nil {
			if snap.FetchedAtUnix == lastPrinted {
				return
			}
			lastPrinted = snap.FetchedAtUnix
			fmt.Printf("USD/TRY: %.4f\n", snap.USDTRY)
			fmt.Printf("EUR/TRY: %.4f\n", snap.EURTRY)
			fmt.Printf("GBP/TRY: %.4f\n", snap.GBPTRY)
			fmt.Printf("JPY/TRY: %.4f\n", snap.JPYTRY)
			fmt.Printf("CHF/TRY: %.4f\n", snap.CHFTRY)
			fmt.Printf("CAD/TRY: %.4f\n", snap.CADTRY)
			fmt.Printf("AUD/TRY: %.4f\n", snap.AUDTRY)
			fmt.Printf("NZD/TRY: %.4f\n", snap.NZDTRY)
			fmt.Printf("SEK/TRY: %.4f\n", snap.SEKTRY)
			fmt.Printf("NOK/TRY: %.4f\n", snap.NOKTRY)
			fmt.Printf("DKK/TRY: %.4f\n", snap.DKKTRY)
			fmt.Printf("PLN/TRY: %.4f\n", snap.PLNTRY)
			fmt.Printf("Gram Altin (TRY): %.2f\n", snap.GramAltinTRY)
			fmt.Printf("Guncelleme: %s (source: %s)\n", time.Unix(snap.FetchedAtUnix, 0).Format(time.RFC3339), snap.Source)

			if store != nil {
				if err := store.UpsertSnapshot(ctx, snap); err != nil {
					log.Printf("db guncelleme hatasi: %v", err)
				}
			}
		} else {
			log.Printf("fetch failed: %v", err)
		}
	}

	printSnapshot()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			printSnapshot()
		}
	}()

	addr := ":8090"
	server := httpapi.NewServer(addr, service)
	log.Printf("listening on %s", addr)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
