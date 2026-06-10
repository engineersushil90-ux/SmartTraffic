package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"smarttraffic/gateway/internal/config"
	"smarttraffic/gateway/internal/services"
	"smarttraffic/gateway/internal/stream"
)

type Server struct {
	cfg       config.Config
	hub       *stream.Hub
	services  *services.Registry
	atccProxy *httputil.ReverseProxy
	ptzProxy  *httputil.ReverseProxy
}

type statusResponse struct {
	OK          bool             `json:"ok"`
	InputURL    string           `json:"inputUrl"`
	StreamURL   string           `json:"streamUrl"`
	BufferBytes int              `json:"bufferBytes"`
	ATCCService string           `json:"atccService"`
	PTZService  string           `json:"ptzService"`
	Services    []map[string]any `json:"services"`
}

type ptzRequest struct {
	Command string `json:"command"`
}

func New(cfg config.Config, hub *stream.Hub, registry *services.Registry) *http.Server {
	app := &Server{cfg: cfg, hub: hub, services: registry}
	app.atccProxy = newReverseProxy(cfg.ATCCServiceURL, "ATCC service")
	app.ptzProxy = newReverseProxy(cfg.PTZServiceURL, "PTZ service")
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/live", withCORS(app.hub.HandleFLV))
	mux.HandleFunc("/api/services", withCORS(app.handleServiceSummaries))
	mux.HandleFunc("/api/atcc", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/atcc/", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/atcc-events", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/vids", withCORS(app.handleDeviceCollection(app.services.VIDS)))
	mux.HandleFunc("/api/vids/", withCORS(app.handleDeviceDetail(app.services.VIDS, "/api/vids/")))
	mux.HandleFunc("/api/ptz-cameras", withCORS(app.handlePTZProxy))
	mux.HandleFunc("/api/ptz-cameras/", withCORS(app.handlePTZProxy))
	mux.HandleFunc("/api/cctv-cameras", withCORS(app.handleDeviceCollection(app.services.CCTVCameras)))
	mux.HandleFunc("/api/cctv-cameras/", withCORS(app.handleDeviceDetail(app.services.CCTVCameras, "/api/cctv-cameras/")))
	mux.HandleFunc("/api/met", withCORS(app.handleDeviceCollection(app.services.MET)))
	mux.HandleFunc("/api/met/", withCORS(app.handleDeviceDetail(app.services.MET, "/api/met/")))
	mux.HandleFunc("/api/vms", withCORS(app.handleDeviceCollection(app.services.VMS)))
	mux.HandleFunc("/api/vms/", withCORS(app.handleDeviceDetail(app.services.VMS, "/api/vms/")))
	mux.HandleFunc("/api/vsds", withCORS(app.handleDeviceCollection(app.services.VSDS)))
	mux.HandleFunc("/api/vsds/", withCORS(app.handleDeviceDetail(app.services.VSDS, "/api/vsds/")))
	mux.HandleFunc("/api/ptz/", withCORS(app.handlePTZProxy))

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func newReverseProxy(targetURL string, label string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("%s proxy disabled: invalid url %q", label, targetURL)
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("%s proxy error path=%s err=%v", label, r.URL.Path, err)
		http.Error(w, label+" unavailable", http.StatusBadGateway)
	}
	return proxy
}

func (s *Server) handleATCCProxy(w http.ResponseWriter, r *http.Request) {
	if !allowProxyMethod(w, r, http.MethodGet) {
		return
	}
	if s.atccProxy == nil {
		http.Error(w, "ATCC service proxy is not configured", http.StatusServiceUnavailable)
		return
	}
	s.atccProxy.ServeHTTP(w, r)
}

func (s *Server) handlePTZProxy(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.ptzProxy == nil {
		http.Error(w, "PTZ service proxy is not configured", http.StatusServiceUnavailable)
		return
	}
	s.ptzProxy.ServeHTTP(w, r)
}

func allowProxyMethod(w http.ResponseWriter, r *http.Request, method string) bool {
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

func (s *Server) handleServiceSummaries(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, map[string]any{"services": s.services.Summaries()})
}

func (s *Server) handleDeviceCollection(service *services.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		writeJSON(w, map[string]any{"devices": service.List(), "summary": service.Summary()})
	}
}

func (s *Server) handleDeviceDetail(service *services.DeviceService, prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		deviceID := strings.TrimPrefix(r.URL.Path, prefix)
		device, ok := service.Get(deviceID)
		if !ok {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		writeJSON(w, device)
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
		ATCCService: s.cfg.ATCCServiceURL,
		PTZService:  s.cfg.PTZServiceURL,
		Services:    s.services.Summaries(),
	})
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
