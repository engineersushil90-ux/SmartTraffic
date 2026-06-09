package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr        string
	InputURL    string
	FFmpegPath  string
	BufferBytes int
}

func Load() Config {
	return Config{
		Addr:        env("STREAM_GATEWAY_ADDR", ":8080"),
		InputURL:    env("STREAM_INPUT_RTSP", "rtsp://localhost:8554/webcam"),
		FFmpegPath:  env("FFMPEG_PATH", "ffmpeg"),
		BufferBytes: envInt("STREAM_BUFFER_BYTES", 4*1024*1024),
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
