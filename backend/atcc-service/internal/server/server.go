package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"smarttraffic/atcc-service/internal/atcc"
	"smarttraffic/atcc-service/internal/config"
)

type Server struct {
	atcc *atcc.Service
}

func New(cfg config.Config, service *atcc.Service) *http.Server {
	app := &Server{atcc: service}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/api/atcc", withCORS(app.handleDevices))
	mux.HandleFunc("/api/atcc/", withCORS(app.handleDevice))
	mux.HandleFunc("/api/atcc-events", withCORS(app.handleEvents))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}
}

func Shutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}

	writeJSON(w, map[string]any{
		"ok":      true,
		"service": "atcc",
		"summary": s.atcc.Summary(),
	})
}

func (s *Server) handleDevices(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}

	writeJSON(w, map[string]any{
		"devices": s.atcc.ListDevices(),
		"summary": s.atcc.Summary(),
	})
}

func (s *Server) handleDevice(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}

	deviceID := strings.TrimPrefix(r.URL.Path, "/api/atcc/")
	device, ok := s.atcc.GetDevice(deviceID)
	if !ok {
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	writeJSON(w, device)
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}

	writeJSON(w, map[string]any{"events": s.atcc.ListEvents()})
}

func allowMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	if r.Method != method {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	return true
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
