package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"scraper/internal/market"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) Close() error {
	return s.db.Close()
}

func (s *MySQLStore) EnsureSchema(ctx context.Context) error {
	const q = `
SELECT COUNT(*)
FROM information_schema.statistics
WHERE table_schema = DATABASE()
  AND table_name = 'piyasa_fiyat'
  AND column_name = 'piyasa_id'
  AND non_unique = 0
`
	var count int
	if err := s.db.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("piyasa_fiyat.piyasa_id icin UNIQUE index bulunamadi")
	}
	return nil
}

func (s *MySQLStore) UpsertSnapshot(ctx context.Context, snap market.Snapshot) error {
	items := []struct {
		piyasaAdi  string
		fiyat      float64
		paraBirimi string
	}{
		{piyasaAdi: "USDTRY", fiyat: snap.USDTRY, paraBirimi: "TRY"},
		{piyasaAdi: "EURTRY", fiyat: snap.EURTRY, paraBirimi: "TRY"},
		{piyasaAdi: "GBPTRY", fiyat: snap.GBPTRY, paraBirimi: "TRY"},
		{piyasaAdi: "JPYTRY", fiyat: snap.JPYTRY, paraBirimi: "TRY"},
		{piyasaAdi: "CHFTRY", fiyat: snap.CHFTRY, paraBirimi: "TRY"},
		{piyasaAdi: "CADTRY", fiyat: snap.CADTRY, paraBirimi: "TRY"},
		{piyasaAdi: "AUDTRY", fiyat: snap.AUDTRY, paraBirimi: "TRY"},
		{piyasaAdi: "NZDTRY", fiyat: snap.NZDTRY, paraBirimi: "TRY"},
		{piyasaAdi: "SEKTRY", fiyat: snap.SEKTRY, paraBirimi: "TRY"},
		{piyasaAdi: "NOKTRY", fiyat: snap.NOKTRY, paraBirimi: "TRY"},
		{piyasaAdi: "DKKTRY", fiyat: snap.DKKTRY, paraBirimi: "TRY"},
		{piyasaAdi: "PLNTRY", fiyat: snap.PLNTRY, paraBirimi: "TRY"},
		// NOTE: Assumes XAUTRY maps to gram altin TRY.
		{piyasaAdi: "XAUTRY", fiyat: snap.GramAltinTRY, paraBirimi: "TRY"},
	}

	for _, item := range items {
		if err := s.UpsertPrice(ctx, item.piyasaAdi, item.fiyat, item.paraBirimi, snap.Source, time.Unix(snap.FetchedAtUnix, 0)); err != nil {
			return err
		}
	}
	return nil
}

func (s *MySQLStore) UpsertPrice(ctx context.Context, piyasaAdi string, fiyat float64, paraBirimi string, kaynak string, zaman time.Time) error {
	var piyasaID int64
	row := s.db.QueryRowContext(ctx, `SELECT piyasa_id FROM piyasa_adlari WHERE piyasa_adi = ? LIMIT 1`, piyasaAdi)
	if err := row.Scan(&piyasaID); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("piyasa_adi bulunamadi: %s", piyasaAdi)
		}
		return err
	}

	updateRes, err := s.db.ExecContext(ctx, `
UPDATE piyasa_fiyat
SET guncel_fiyat = ?, para_birimi = ?, zaman = ?, kaynak = ?
WHERE piyasa_id = ?
`, fiyat, paraBirimi, zaman, kaynak, piyasaID)
	if err != nil {
		return err
	}
	affected, err := updateRes.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 1 {
		return fmt.Errorf("piyasa_id %d icin birden fazla satir var (%d)", piyasaID, affected)
	}
	if affected == 1 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, `
INSERT INTO piyasa_fiyat (piyasa_id, guncel_fiyat, para_birimi, zaman, kaynak)
VALUES (?, ?, ?, ?, ?)
`, piyasaID, fiyat, paraBirimi, zaman, kaynak)
	return err
}
