# Market Service (Doviz + Altin)

Bu servis USD/TRY, EUR/TRY, GBP/TRY, JPY/TRY, CHF/TRY, CAD/TRY, AUD/TRY, NZD/TRY, SEK/TRY, NOK/TRY, DKK/TRY, PLN/TRY ve Gram Altin (TRY) fiyatlarini ceker, konsola yazdirir ve REST API olarak sunar.

## Calistirma

Tek komut:

```bash
# secenek 1
go run .

# secenek 2
go run ./cmd/server
```

Calistiginda terminalde ornek cikti verir ve HTTP sunucusu 8090 portunda acilir.

## API

- `GET /healthz` -> `ok`
- `GET /api/market/latest`

Ornek cevap:

```json
{
  "ok": true,
  "ts": 1700000000,
  "source": "fx:Frankfurter, gold:SILV.DATA",
  "data": {
    "USDTRY": {"value": 43.60},
    "EURTRY": {"value": 51.73},
    "GBPTRY": {"value": 61.20},
    "JPYTRY": {"value": 0.2890},
    "CHFTRY": {"value": 55.90},
    "CADTRY": {"value": 31.70},
    "AUDTRY": {"value": 28.40},
    "NZDTRY": {"value": 26.10},
    "SEKTRY": {"value": 4.20},
    "NOKTRY": {"value": 4.00},
    "DKKTRY": {"value": 7.50},
    "PLNTRY": {"value": 12.80},
    "GRAM_ALTIN_TRY": {"value": 7025.85}
  }
}
```

## Kaynaklar ve Donusum

- Doviz: Frankfurter (primary) ve MoneyConvert (fallback)
- Altin (USD/ons): SILV.DATA (primary) ve FreeGoldAPI (fallback)
- Gram Altin (TRY) hesaplama:
  - USD/ons -> USD/gram: `USD_per_gram = USD_per_ounce / 31.1034768`
  - TRY/gram: `USD_per_gram * USDTRY`

## Notlar

- 5 saniye TTL in-memory cache vardir.
- Es zamanli isteklerde tek fetch calisir, digerleri cache sonucunu alir.
- Timeout, context ve hata yonetimi eklenmistir.
- Windows ile uyumludur.
