package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr           string
	InputURL       string
	FFmpegPath     string
	RTSPTransport  string
	VideoCodec     string
	BufferBytes    int
	ATCCServiceURL string
	PTZServiceURL  string
}

func Load() Config {
	return Config{
		Addr:           env("STREAM_GATEWAY_ADDR", ":8080"),
		InputURL:       env("STREAM_INPUT_RTSP", "rtsp://localhost:8554/live"),
		FFmpegPath:     env("FFMPEG_PATH", "ffmpeg"),
		RTSPTransport:  env("STREAM_RTSP_TRANSPORT", "tcp"),
		VideoCodec:     env("STREAM_VIDEO_CODEC", "libx264"),
		BufferBytes:    envInt("STREAM_BUFFER_BYTES", 4*1024*1024),
		ATCCServiceURL: env("ATCC_SERVICE_URL", "http://localhost:8091"),
		PTZServiceURL:  env("PTZ_SERVICE_URL", "http://localhost:8092"),
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
