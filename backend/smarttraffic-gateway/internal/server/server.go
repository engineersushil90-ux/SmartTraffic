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
	"smarttraffic/gateway/internal/supervisor"
)

type Server struct {
	cfg             config.Config
	hub             *stream.Hub
	services        *services.Registry
	serviceStatuses func() []supervisor.ServiceStatus
	atccProxy       *httputil.ReverseProxy
}

type statusResponse struct {
	OK          bool                       `json:"ok"`
	InputURL    string                     `json:"inputUrl"`
	StreamURL   string                     `json:"streamUrl"`
	BufferBytes int                        `json:"bufferBytes"`
	ATCCService string                     `json:"atccService"`
	Services    []map[string]any           `json:"services"`
	Managed     []supervisor.ServiceStatus `json:"managedServices"`
}

type ptzRequest struct {
	Command string `json:"command"`
}

func New(cfg config.Config, hub *stream.Hub, serviceStatuses func() []supervisor.ServiceStatus) *http.Server {
	app := &Server{cfg: cfg, hub: hub, services: services.NewRegistry(), serviceStatuses: serviceStatuses}
	app.atccProxy = newReverseProxy(cfg.ATCCServiceURL, "ATCC service")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/live", withCORS(app.hub.HandleFLV))
	mux.HandleFunc("/api/services", withCORS(app.handleServiceSummaries))
	mux.HandleFunc("/api/atcc", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/atcc/", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/atcc-events", withCORS(app.handleATCCProxy))
	mux.HandleFunc("/api/vids", withCORS(app.handleDeviceCollection(app.services.VIDS)))
	mux.HandleFunc("/api/vids/", withCORS(app.handleDeviceDetail(app.services.VIDS, "/api/vids/")))
	mux.HandleFunc("/api/ptz-cameras", withCORS(app.handleDeviceCollection(app.services.PTZCameras.DeviceService)))
	mux.HandleFunc("/api/ptz-cameras/", withCORS(app.handleDeviceDetail(app.services.PTZCameras.DeviceService, "/api/ptz-cameras/")))
	mux.HandleFunc("/api/cctv-cameras", withCORS(app.handleDeviceCollection(app.services.CCTVCameras)))
	mux.HandleFunc("/api/cctv-cameras/", withCORS(app.handleDeviceDetail(app.services.CCTVCameras, "/api/cctv-cameras/")))
	mux.HandleFunc("/api/met", withCORS(app.handleDeviceCollection(app.services.MET)))
	mux.HandleFunc("/api/met/", withCORS(app.handleDeviceDetail(app.services.MET, "/api/met/")))
	mux.HandleFunc("/api/vms", withCORS(app.handleDeviceCollection(app.services.VMS)))
	mux.HandleFunc("/api/vms/", withCORS(app.handleDeviceDetail(app.services.VMS, "/api/vms/")))
	mux.HandleFunc("/api/vsds", withCORS(app.handleDeviceCollection(app.services.VSDS)))
	mux.HandleFunc("/api/vsds/", withCORS(app.handleDeviceDetail(app.services.VSDS, "/api/vsds/")))
	mux.HandleFunc("/api/ptz/", withCORS(app.handlePTZ))

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
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.atccProxy == nil {
		http.Error(w, "ATCC service proxy is not configured", http.StatusServiceUnavailable)
		return
	}

	s.atccProxy.ServeHTTP(w, r)
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
		Services:    s.services.Summaries(),
		Managed:     s.managedServiceStatuses(),
	})
}

func (s *Server) managedServiceStatuses() []supervisor.ServiceStatus {
	if s.serviceStatuses == nil {
		return nil
	}

	return s.serviceStatuses()
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

	result, err := s.services.PTZCameras.SendCommand(cameraID, req.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("PTZ command camera=%s command=%s", cameraID, req.Command)
	writeJSON(w, result)
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
