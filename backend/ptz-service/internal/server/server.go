package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"smarttraffic/ptz-service/internal/config"
	"smarttraffic/ptz-service/internal/ptz"
)

type Server struct {
	ptz *ptz.Service
}

type commandRequest struct {
	Command string `json:"command"`
}

func New(cfg config.Config, service *ptz.Service) *http.Server {
	app := &Server{ptz: service}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/api/ptz-cameras", withCORS(app.handleCameras))
	mux.HandleFunc("/api/ptz-cameras/", withCORS(app.handleCamera))
	mux.HandleFunc("/api/ptz/", withCORS(app.handleCommand))
	return &http.Server{Addr: cfg.Addr, Handler: mux, ReadHeaderTimeout: cfg.ReadHeaderTimeout}
}

func Shutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, map[string]any{"ok": true, "service": "ptz", "summary": s.ptz.Summary()})
}

func (s *Server) handleCameras(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, map[string]any{"devices": s.ptz.ListCameras(), "summary": s.ptz.Summary()})
}

func (s *Server) handleCamera(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	cameraID := strings.TrimPrefix(r.URL.Path, "/api/ptz-cameras/")
	camera, ok := s.ptz.GetCamera(cameraID)
	if !ok {
		http.Error(w, "camera not found", http.StatusNotFound)
		return
	}
	writeJSON(w, camera)
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	cameraID := strings.TrimPrefix(r.URL.Path, "/api/ptz/")
	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	result, err := s.ptz.SendCommand(cameraID, req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, result)
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
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
