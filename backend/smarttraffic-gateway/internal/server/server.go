package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"smarttraffic/gateway/internal/config"
	"smarttraffic/gateway/internal/services"
	"smarttraffic/gateway/internal/stream"
	"smarttraffic/gateway/internal/upstreams"
)

type Server struct {
	cfg      config.Config
	hub      *stream.Hub
	services *services.Registry
	clients  *ClientDirectory
}

type statusResponse struct {
	OK          bool               `json:"ok"`
	InputURL    string             `json:"inputUrl"`
	StreamURL   string             `json:"streamUrl"`
	BufferBytes int                `json:"bufferBytes"`
	ATCCService string             `json:"atccService"`
	PTZService  string             `json:"ptzService"`
	Upstreams   []upstreams.Status `json:"upstreams"`
	Services    []map[string]any   `json:"services"`
}

type clientRegistration struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	HealthURL string `json:"healthUrl"`
}

type ClientDirectory struct {
	mu      sync.RWMutex
	clients map[string]upstreams.Spec
}

func NewClientDirectory(defaults []upstreams.Spec) *ClientDirectory {
	directory := &ClientDirectory{clients: make(map[string]upstreams.Spec, len(defaults))}
	for _, spec := range defaults {
		directory.Register(spec)
	}
	return directory
}

func (d *ClientDirectory) Register(spec upstreams.Spec) {
	spec.Name = normalizeClientName(spec.Name)
	spec.URL = strings.TrimRight(spec.URL, "/")
	if strings.TrimSpace(spec.HealthURL) == "" && spec.URL != "" {
		spec.HealthURL = strings.TrimRight(spec.URL, "/") + "/healthz"
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	d.clients[spec.Name] = spec
}

func validateHTTPURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return url.InvalidHostError("scheme must be http or https")
	}
	if parsed.Host == "" {
		return url.InvalidHostError("host is required")
	}
	return nil
}

func (d *ClientDirectory) Get(name string) (upstreams.Spec, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	spec, ok := d.clients[normalizeClientName(name)]
	return spec, ok
}

func (d *ClientDirectory) Specs() []upstreams.Spec {
	d.mu.RLock()
	defer d.mu.RUnlock()
	specs := make([]upstreams.Spec, 0, len(d.clients))
	for _, name := range []string{"ATCC", "PTZ Camera"} {
		if spec, ok := d.clients[name]; ok {
			specs = append(specs, spec)
		}
	}
	return specs
}

func normalizeClientName(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "atcc":
		return "ATCC"
	case "ptz", "ptz-camera", "ptz-cameras":
		return "PTZ Camera"
	default:
		return strings.TrimSpace(name)
	}
}

func New(cfg config.Config, hub *stream.Hub, registry *services.Registry) *http.Server {
	clients := NewClientDirectory(upstreams.ServiceSpecs(cfg.ATCCServiceURL, cfg.PTZServiceURL))
	app := &Server{
		cfg:      cfg,
		hub:      hub,
		services: registry,
		clients:  clients,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", withCORS(app.handleHealth))
	mux.HandleFunc("/live", withCORS(app.hub.HandleFLV))
	mux.HandleFunc("/api/clients", withCORS(app.handleClients))
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

func proxyTo(w http.ResponseWriter, r *http.Request, targetURL string, label string) {
	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, label+" has invalid client url", http.StatusBadGateway)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("%s proxy error path=%s err=%v", label, r.URL.Path, err)
		http.Error(w, label+" unavailable", http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

func (s *Server) handleATCCProxy(w http.ResponseWriter, r *http.Request) {
	if !allowProxyMethod(w, r, http.MethodGet) {
		return
	}
	client, ok := s.clients.Get("atcc")
	if !ok || client.URL == "" {
		http.Error(w, "ATCC client is not registered", http.StatusServiceUnavailable)
		return
	}
	proxyTo(w, r, client.URL, "ATCC client")
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
	client, ok := s.clients.Get("ptz")
	if !ok || client.URL == "" {
		http.Error(w, "PTZ client is not registered", http.StatusServiceUnavailable)
		return
	}
	proxyTo(w, r, client.URL, "PTZ client")
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

	checkCtx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	checker := upstreams.NewChecker(s.clients.Specs())

	writeJSON(w, map[string]any{
		"upstreams": checker.CheckAll(checkCtx),
		"services":  s.services.Summaries(),
	})
}

func (s *Server) handleClients(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, map[string]any{"clients": s.clients.Specs()})
	case http.MethodPost:
		var req clientRegistration
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.URL) == "" {
			http.Error(w, "name and url are required", http.StatusBadRequest)
			return
		}
		if err := validateHTTPURL(strings.TrimSpace(req.URL)); err != nil {
			http.Error(w, "client url must be a valid http url", http.StatusBadRequest)
			return
		}
		s.clients.Register(upstreams.Spec{Name: req.Name, URL: req.URL, HealthURL: req.HealthURL})
		writeJSON(w, map[string]any{"ok": true, "client": normalizeClientName(req.Name)})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
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

	checkCtx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	checker := upstreams.NewChecker(s.clients.Specs())
	upstreamStatuses := checker.CheckAll(checkCtx)

	ok := true
	for _, status := range upstreamStatuses {
		if !status.Connected {
			ok = false
			break
		}
	}

	writeJSON(w, statusResponse{
		OK:          ok,
		InputURL:    s.cfg.InputURL,
		StreamURL:   "http://localhost" + s.cfg.Addr + "/live",
		BufferBytes: s.cfg.BufferBytes,
		ATCCService: s.cfg.ATCCServiceURL,
		PTZService:  s.cfg.PTZServiceURL,
		Upstreams:   upstreamStatuses,
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
