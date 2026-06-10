package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Addr                  string
	InputURL              string
	FFmpegPath            string
	RTSPTransport         string
	BufferBytes           int
	ManageServices        bool
	ATCCServiceURL        string
	ATCCServiceExecutable string
}

func Load() Config {
	return Config{
		Addr:                  env("STREAM_GATEWAY_ADDR", ":8080"),
		InputURL:              env("STREAM_INPUT_RTSP", "rtsp://localhost:8554/live"),
		FFmpegPath:            env("FFMPEG_PATH", "ffmpeg"),
		RTSPTransport:         env("STREAM_RTSP_TRANSPORT", "tcp"),
		BufferBytes:           envInt("STREAM_BUFFER_BYTES", 4*1024*1024),
		ManageServices:        envBool("SMARTTRAFFIC_MANAGE_SERVICES", true),
		ATCCServiceURL:        env("ATCC_SERVICE_URL", "http://localhost:8091"),
		ATCCServiceExecutable: env("ATCC_SERVICE_EXE", defaultATCCServiceExecutable()),
	}
}

func env(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
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

func envBool(name string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func defaultATCCServiceExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		return filepath.Clean(filepath.Join("..", "atcc-service", "atcc-service.exe"))
	}

	return filepath.Clean(filepath.Join(filepath.Dir(exePath), "..", "atcc-service", "atcc-service.exe"))
}
