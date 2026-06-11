package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ATCCExecutable    string
	PTZExecutable     string
	GatewayExecutable string
	ATCCURL           string
	PTZURL            string
	GatewayURL        string
}

func Load() Config {
	backendRoot := findBackendRoot()
	return Config{
		ATCCExecutable:    envPath("ATCC_SERVICE_EXE", filepath.Join(backendRoot, "atcc-service", "atcc-service.exe")),
		PTZExecutable:     envPath("PTZ_SERVICE_EXE", filepath.Join(backendRoot, "ptz-service", "ptz-service.exe")),
		GatewayExecutable: envPath("SMARTTRAFFIC_GATEWAY_EXE", filepath.Join(backendRoot, "smarttraffic-gateway", "smarttraffic-gateway.exe")),
		ATCCURL:           envURL("ATCC_SERVICE_URL", "http://localhost:8091"),
		PTZURL:            envURL("PTZ_SERVICE_URL", "http://localhost:8092"),
		GatewayURL:        envURL("SMARTTRAFFIC_GATEWAY_URL", "http://localhost:8080"),
	}
}

func envPath(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return filepath.Clean(fallback)
	}
	return filepath.Clean(value)
}

func envURL(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return strings.TrimRight(value, "/")
}

func findBackendRoot() string {
	wd, err := os.Getwd()
	if err == nil {
		for dir := filepath.Clean(wd); ; dir = filepath.Dir(dir) {
			if hasBackendServices(dir) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		parent := filepath.Clean(filepath.Join(exeDir, ".."))
		if hasBackendServices(parent) {
			return parent
		}
	}

	return filepath.Clean(filepath.Join("..", ".."))
}

func hasBackendServices(dir string) bool {
	for _, name := range []string{"atcc-service", "ptz-service", "smarttraffic-gateway"} {
		if info, err := os.Stat(filepath.Join(dir, name)); err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}
