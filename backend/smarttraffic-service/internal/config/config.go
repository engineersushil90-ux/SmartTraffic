package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ATCCExecutable    string
	PTZExecutable     string
	GatewayExecutable string
	ATCCURL           string
	PTZURL            string
	GatewayURL        string
}

func Load() Config {
	baseDir := executableDir()
	return Config{
		Addr:              env("SMARTTRAFFIC_SERVICE_ADDR", ":8079"),
		ReadHeaderTimeout: time.Duration(envInt("SMARTTRAFFIC_READ_HEADER_TIMEOUT_SECONDS", 5)) * time.Second,
		ATCCExecutable:    env("ATCC_SERVICE_EXE", filepath.Join(baseDir, "..", "atcc-service", "atcc-service.exe")),
		PTZExecutable:     env("PTZ_SERVICE_EXE", filepath.Join(baseDir, "..", "ptz-service", "ptz-service.exe")),
		GatewayExecutable: env("SMARTTRAFFIC_GATEWAY_EXE", filepath.Join(baseDir, "..", "smarttraffic-gateway", "smarttraffic-gateway.exe")),
		ATCCURL:           env("ATCC_SERVICE_URL", "http://localhost:8091"),
		PTZURL:            env("PTZ_SERVICE_URL", "http://localhost:8092"),
		GatewayURL:        env("SMARTTRAFFIC_GATEWAY_URL", "http://localhost:8080"),
	}
}

func env(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return filepath.Clean(fallback)
	}
	return value
}

func envInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func executableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			return "."
		}
		return wd
	}
	return filepath.Dir(exePath)
}
