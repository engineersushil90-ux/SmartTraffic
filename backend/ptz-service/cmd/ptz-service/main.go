package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"strings"
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
	go registerWithGateway(ctx, cfg)

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx, srv)
	case err := <-errs:
		return err
	}
}

func registerWithGateway(ctx context.Context, cfg config.Config) {
	payload := map[string]string{
		"name":      "ptz",
		"url":       strings.TrimRight(cfg.PublicURL, "/"),
		"healthUrl": strings.TrimRight(cfg.PublicURL, "/") + "/healthz",
	}
	endpoint := strings.TrimRight(cfg.GatewayURL, "/") + "/api/clients"
	registered := false

	for {
		if registerClient(ctx, endpoint, payload) == nil {
			if !registered {
				log.Printf("PTZ registered with gateway at %s", cfg.GatewayURL)
			}
			registered = true
			if !wait(ctx, 10*time.Second) {
				return
			}
			continue
		}

		if !wait(ctx, 2*time.Second) {
			return
		}
	}
}

func wait(ctx context.Context, interval time.Duration) bool {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func registerClient(ctx context.Context, endpoint string, payload map[string]string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("gateway returned %s", resp.Status)
	}
	return nil
}
