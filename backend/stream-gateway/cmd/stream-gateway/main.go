package main

import (
	"context"
	"log"

	"smarttraffic/stream-gateway/internal/config"
	"smarttraffic/stream-gateway/internal/server"
	"smarttraffic/stream-gateway/internal/stream"
)

func main() {
	cfg := config.Load()
	hub := stream.NewHub(cfg.BufferBytes)
	runner := stream.NewFFmpegRunner(cfg.FFmpegPath, cfg.InputURL, cfg.RTSPTransport, hub)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runner.RunLoop(ctx)

	app := server.New(cfg, hub)
	log.Fatal(app.ListenAndServe())
}
