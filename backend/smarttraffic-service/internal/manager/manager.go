package manager

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	procs  map[string]*managedProcess
	status map[string]Status
	client *http.Client
}

type managedProcess struct {
	cmd  *exec.Cmd
	done chan error
}

func New(specs []Spec) *Manager {
	status := make(map[string]Status, len(specs))
	for _, spec := range specs {
		status[spec.Name] = Status{Name: spec.Name, URL: spec.URL, Executable: spec.Executable}
	}
	return &Manager{
		specs:  specs,
		procs:  make(map[string]*managedProcess),
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
	for index := len(m.specs) - 1; index >= 0; index-- {
		name := m.specs[index].Name
		m.mu.Lock()
		proc := m.procs[name]
		m.mu.Unlock()
		if proc == nil || proc.cmd.Process == nil {
			continue
		}

		select {
		case <-proc.done:
			m.markStopped(name, "")
			continue
		default:
		}

		if err := proc.cmd.Process.Kill(); err != nil {
			if !isProcessAlreadyFinished(err) {
				log.Printf("failed to stop %s: %v", name, err)
			}
			m.markStopped(name, err.Error())
			continue
		}

		select {
		case err := <-proc.done:
			if err != nil && !isProcessAlreadyFinished(err) {
				log.Printf("service %s stopped: %v", name, err)
			}
		case <-time.After(5 * time.Second):
			log.Printf("timed out waiting for %s to stop", name)
		}
	}
}

func (m *Manager) markStopped(name string, errorMessage string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	status := m.status[name]
	status.Running = false
	status.PID = 0
	status.Error = errorMessage
	m.status[name] = status
	delete(m.procs, name)
}

func isProcessAlreadyFinished(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "invalid argument") ||
		strings.Contains(message, "process already finished") ||
		strings.Contains(message, "no process")
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

	if ctx.Err() != nil {
		status.Error = ctx.Err().Error()
		m.status[spec.Name] = status
		return
	}

	cmd := exec.Command(spec.Executable)
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
	proc := &managedProcess{cmd: cmd, done: make(chan error, 1)}
	m.procs[spec.Name] = proc
	m.status[spec.Name] = status

	go m.wait(spec.Name, proc)
	log.Printf("started service name=%s pid=%d", spec.Name, cmd.Process.Pid)
}

func (m *Manager) wait(name string, proc *managedProcess) {
	err := proc.cmd.Wait()
	proc.done <- err
	m.mu.Lock()
	defer m.mu.Unlock()

	if current := m.procs[name]; current != proc {
		return
	}

	status := m.status[name]
	status.Running = false
	status.PID = 0
	if err != nil && !isProcessAlreadyFinished(err) {
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
