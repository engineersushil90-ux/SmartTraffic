package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"smarttraffic/ptz-service/internal/config"
	"smarttraffic/ptz-service/internal/ptz"
	"smarttraffic/ptz-service/internal/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg := config.Load()
	srv := server.New(cfg, ptz.NewService())
	errs := make(chan error, 1)
	go func() {
		log.Printf("PTZ service listening on %s", cfg.Addr)
		errs <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx, srv)
	case err := <-errs:
		return err
	}
}
