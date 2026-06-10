package supervisor

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type ServiceSpec struct {
	Name       string
	URL        string
	HealthURL  string
	Executable string
	Args       []string
	Env        []string
}

type ServiceStatus struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Managed    bool   `json:"managed"`
	Running    bool   `json:"running"`
	Started    bool   `json:"started"`
	Executable string `json:"executable,omitempty"`
	Error      string `json:"error,omitempty"`
}

type Supervisor struct {
	enabled bool
	specs   []ServiceSpec
	mu      sync.Mutex
	procs   map[string]*exec.Cmd
	status  map[string]ServiceStatus
	client  *http.Client
}

func New(enabled bool, specs []ServiceSpec) *Supervisor {
	status := make(map[string]ServiceStatus, len(specs))
	for _, spec := range specs {
		status[spec.Name] = ServiceStatus{
			Name:       spec.Name,
			URL:        spec.URL,
			Managed:    enabled,
			Executable: spec.Executable,
		}
	}

	return &Supervisor{
		enabled: enabled,
		specs:   specs,
		procs:   make(map[string]*exec.Cmd),
		status:  status,
		client:  &http.Client{Timeout: 2 * time.Second},
	}
}

func (s *Supervisor) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, spec := range s.specs {
		status := s.status[spec.Name]
		status.Running = s.isHealthy(spec)

		if !s.enabled {
			status.Error = "service management disabled"
			s.status[spec.Name] = status
			continue
		}

		if status.Running {
			status.Started = false
			status.Error = ""
			s.status[spec.Name] = status
			continue
		}

		if _, err := os.Stat(spec.Executable); err != nil {
			status.Error = "service executable not found"
			s.status[spec.Name] = status
			continue
		}

		cmd := exec.CommandContext(ctx, spec.Executable, spec.Args...)
		cmd.Dir = filepath.Dir(spec.Executable)
		cmd.Env = append(os.Environ(), spec.Env...)
		if err := cmd.Start(); err != nil {
			status.Error = err.Error()
			s.status[spec.Name] = status
			continue
		}

		s.procs[spec.Name] = cmd
		status.Started = true
		status.Error = ""
		s.status[spec.Name] = status

		go s.wait(spec.Name, cmd)
		log.Printf("started managed service name=%s pid=%d", spec.Name, cmd.Process.Pid)
	}
}

func (s *Supervisor) Stop() {
	s.mu.Lock()
	procs := make(map[string]*exec.Cmd, len(s.procs))
	for name, cmd := range s.procs {
		procs[name] = cmd
	}
	s.mu.Unlock()

	for name, cmd := range procs {
		if cmd.Process == nil {
			continue
		}
		if err := cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Printf("failed to stop managed service name=%s err=%v", name, err)
		}
	}
}

func (s *Supervisor) Statuses() []ServiceStatus {
	s.mu.Lock()
	defer s.mu.Unlock()

	statuses := make([]ServiceStatus, 0, len(s.specs))
	for _, spec := range s.specs {
		status := s.status[spec.Name]
		status.Running = s.isHealthy(spec)
		s.status[spec.Name] = status
		statuses = append(statuses, status)
	}

	return statuses
}

func (s *Supervisor) wait(name string, cmd *exec.Cmd) {
	err := cmd.Wait()

	s.mu.Lock()
	defer s.mu.Unlock()

	status := s.status[name]
	status.Running = false
	if err != nil {
		status.Error = err.Error()
	}
	delete(s.procs, name)
	s.status[name] = status
}

func (s *Supervisor) isHealthy(spec ServiceSpec) bool {
	healthURL := spec.HealthURL
	if healthURL == "" {
		healthURL = spec.URL + "/healthz"
	}

	resp, err := s.client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
