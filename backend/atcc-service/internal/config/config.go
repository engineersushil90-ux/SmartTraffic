package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Addr              string
	ReadHeaderTimeout time.Duration
}

func Load() Config {
	return Config{
		Addr:              env("ATCC_SERVICE_ADDR", ":8091"),
		ReadHeaderTimeout: time.Duration(envInt("ATCC_READ_HEADER_TIMEOUT_SECONDS", 5)) * time.Second,
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
