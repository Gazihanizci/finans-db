# Market Service (Doviz + Altin)

Bu servis USD/TRY, EUR/TRY ve Gram Altin (TRY) fiyatlarini ceker, konsola yazdirir ve REST API olarak sunar.

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
