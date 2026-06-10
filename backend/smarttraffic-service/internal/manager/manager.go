package manager

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Spec struct {
	Name       string
	Executable string
	URL        string
	HealthURL  string
	Env        []string
}

type Status struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Executable string `json:"executable"`
	Running    bool   `json:"running"`
	Started    bool   `json:"started"`
	PID        int    `json:"pid,omitempty"`
	Error      string `json:"error,omitempty"`
}

type Manager struct {
	specs  []Spec
	mu     sync.Mutex
	procs  map[string]*exec.Cmd
	status map[string]Status
	client *http.Client
}

func New(specs []Spec) *Manager {
	status := make(map[string]Status, len(specs))
	for _, spec := range specs {
		status[spec.Name] = Status{Name: spec.Name, URL: spec.URL, Executable: spec.Executable}
	}
	return &Manager{
		specs:  specs,
		procs:  make(map[string]*exec.Cmd),
		status: status,
		client: &http.Client{Timeout: 2 * time.Second},
	}
}

func (m *Manager) StartAll(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, spec := range m.specs {
		m.startLocked(ctx, spec)
	}
}

func (m *Manager) StopAll() {
	m.mu.Lock()
	procs := make(map[string]*exec.Cmd, len(m.procs))
	for name, cmd := range m.procs {
		procs[name] = cmd
	}
	m.mu.Unlock()

	for name, cmd := range procs {
		if cmd.Process == nil {
			continue
		}
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("failed to stop %s: %v", name, err)
		}
	}
}

func (m *Manager) Statuses() []Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	statuses := make([]Status, 0, len(m.specs))
	for _, spec := range m.specs {
		status := m.status[spec.Name]
		status.Running = m.isHealthy(spec)
		statuses = append(statuses, status)
		m.status[spec.Name] = status
	}
	return statuses
}

func (m *Manager) startLocked(ctx context.Context, spec Spec) {
	status := m.status[spec.Name]
	status.Running = m.isHealthy(spec)
	if status.Running {
		status.Error = ""
		m.status[spec.Name] = status
		return
	}

	if _, err := os.Stat(spec.Executable); err != nil {
		status.Error = "executable not found"
		m.status[spec.Name] = status
		return
	}

	cmd := exec.CommandContext(ctx, spec.Executable)
	cmd.Dir = filepath.Dir(spec.Executable)
	cmd.Env = append(os.Environ(), spec.Env...)
	if err := cmd.Start(); err != nil {
		status.Error = err.Error()
		m.status[spec.Name] = status
		return
	}

	status.Started = true
	status.PID = cmd.Process.Pid
	status.Error = ""
	m.procs[spec.Name] = cmd
	m.status[spec.Name] = status

	go m.wait(spec.Name, cmd)
	log.Printf("started service name=%s pid=%d", spec.Name, cmd.Process.Pid)
}

func (m *Manager) wait(name string, cmd *exec.Cmd) {
	err := cmd.Wait()
	m.mu.Lock()
	defer m.mu.Unlock()

	status := m.status[name]
	status.Running = false
	status.PID = 0
	if err != nil {
		status.Error = err.Error()
	}
	delete(m.procs, name)
	m.status[name] = status
}

func (m *Manager) isHealthy(spec Spec) bool {
	healthURL := spec.HealthURL
	if healthURL == "" {
		healthURL = spec.URL + "/healthz"
	}
	resp, err := m.client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
