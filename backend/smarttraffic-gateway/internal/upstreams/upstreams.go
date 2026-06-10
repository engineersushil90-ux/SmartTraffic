package upstreams

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Spec struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	HealthURL string `json:"healthUrl"`
}

type Status struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	HealthURL string `json:"healthUrl"`
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

type Checker struct {
	specs  []Spec
	client *http.Client
}

func NewChecker(specs []Spec) *Checker {
	return &Checker{
		specs: specs,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func ServiceSpecs(atccURL string, ptzURL string) []Spec {
	return []Spec{
		{
			Name:      "ATCC",
			URL:       strings.TrimRight(atccURL, "/"),
			HealthURL: joinURL(atccURL, "/healthz"),
		},
		{
			Name:      "PTZ Camera",
			URL:       strings.TrimRight(ptzURL, "/"),
			HealthURL: joinURL(ptzURL, "/healthz"),
		},
	}
}

func (c *Checker) CheckAll(ctx context.Context) []Status {
	statuses := make([]Status, 0, len(c.specs))
	for _, spec := range c.specs {
		statuses = append(statuses, c.Check(ctx, spec))
	}
	return statuses
}

func (c *Checker) Check(ctx context.Context, spec Spec) Status {
	status := Status{
		Name:      spec.Name,
		URL:       spec.URL,
		HealthURL: spec.HealthURL,
	}
	if strings.TrimSpace(spec.HealthURL) == "" {
		status.Error = "health URL is empty"
		return status
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spec.HealthURL, nil)
	if err != nil {
		status.Error = err.Error()
		return status
	}

	resp, err := c.client.Do(req)
	if err != nil {
		status.Error = err.Error()
		return status
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		status.Error = fmt.Sprintf("health check returned %s", resp.Status)
		return status
	}

	status.Connected = true
	return status
}

func joinURL(base string, path string) string {
	return strings.TrimRight(base, "/") + path
}
