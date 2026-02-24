package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"scraper/internal/market"
)

type MarketGetter interface {
	GetLatest(ctx context.Context) (market.Snapshot, error)
}

type Server struct {
	mux    *http.ServeMux
	market MarketGetter
	addr   string
}

func NewServer(addr string, market MarketGetter) *Server {
	s := &Server{
		mux:    http.NewServeMux(),
		market: market,
		addr:   addr,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux.HandleFunc("/api/market/latest", s.handleLatest)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleLatest(w http.ResponseWriter, r *http.Request) {
	snap, err := s.market.GetLatest(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	resp := map[string]interface{}{
		"ok":     true,
		"ts":     snap.FetchedAtUnix,
		"source": snap.Source,
		"data": map[string]interface{}{
			"USDTRY":         map[string]float64{"value": snap.USDTRY},
			"EURTRY":         map[string]float64{"value": snap.EURTRY},
			"GBPTRY":         map[string]float64{"value": snap.GBPTRY},
			"GRAM_ALTIN_TRY": map[string]float64{"value": snap.GramAltinTRY},
		},
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
