package server

import (
	"encoding/json"
	"net/http"

	"smarttraffic/smarttraffic-service/internal/config"
	"smarttraffic/smarttraffic-service/internal/manager"
)

func New(cfg config.Config, services *manager.Manager) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, map[string]any{"ok": true, "service": "smarttraffic-service", "services": services.Statuses()})
	})
	return &http.Server{Addr: cfg.Addr, Handler: mux, ReadHeaderTimeout: cfg.ReadHeaderTimeout}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
