package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"smarttraffic/stream-gateway/internal/config"
	"smarttraffic/stream-gateway/internal/stream"
)

type Server struct {
	cfg config.Config
	hub *stream.Hub
}

type statusResponse struct {
	OK          bool   `json:"ok"`
	InputURL    string `json:"inputUrl"`
	StreamURL   string `json:"streamUrl"`
	BufferBytes int    `json:"bufferBytes"`
}

type ptzRequest struct {
	Command string `json:"command"`
}

func New(cfg config.Config, hub *stream.Hub) *http.Server {
	app := &Server{cfg: cfg, hub: hub}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/live", withCORS(app.hub.HandleFLV))
	mux.HandleFunc("/api/ptz/", withCORS(app.handlePTZ))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, statusResponse{
		OK:          true,
		InputURL:    s.cfg.InputURL,
		StreamURL:   "http://localhost" + s.cfg.Addr + "/live",
		BufferBytes: s.cfg.BufferBytes,
	})
}

func (s *Server) handlePTZ(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cameraID := strings.TrimPrefix(r.URL.Path, "/api/ptz/")
	var req ptzRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	log.Printf("PTZ command camera=%s command=%s", cameraID, req.Command)
	writeJSON(w, map[string]any{"ok": true, "cameraId": cameraID, "command": req.Command})
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
