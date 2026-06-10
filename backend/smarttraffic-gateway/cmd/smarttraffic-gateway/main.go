package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"smarttraffic/gateway/internal/config"
	"smarttraffic/gateway/internal/server"
	"smarttraffic/gateway/internal/services"
	"smarttraffic/gateway/internal/stream"
)

const (
	serviceName        = "SmartTrafficGatewayService"
	serviceDisplayName = "SmartTraffic Gateway Service"
	serviceDescription = "SmartTraffic gateway service that starts and connects SmartTraffic backend services."
)

type gatewayWindowsService struct{}

func main() {
	serviceAction := flag.String("service", "", "Windows service action: install, uninstall, start, stop")
	flag.Parse()

	if *serviceAction != "" {
		if err := handleServiceAction(*serviceAction); err != nil {
			log.Fatal(err)
		}
		return
	}

	isService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatal(err)
	}

	if isService {
		run := svc.Run
		if len(os.Args) > 1 && os.Args[1] == "debug" {
			run = debug.Run
		}
		if err := run(serviceName, gatewayWindowsService{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := runGateway(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (m gatewayWindowsService) Execute(_ []string, requests <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() {
		errs <- runGateway(ctx)
	}()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	for {
		select {
		case request := <-requests:
			switch request.Cmd {
			case svc.Interrogate:
				changes <- request.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				cancel()
				select {
				case <-errs:
				case <-time.After(15 * time.Second):
				}
				return false, 0
			default:
				continue
			}
		case err := <-errs:
			if err != nil && err != http.ErrServerClosed {
				return false, 1
			}
			return false, 0
		}
	}
}

func runGateway(ctx context.Context) error {
	cfg := config.Load()
	hub := stream.NewHub(cfg.BufferBytes)
	runner := stream.NewFFmpegRunner(cfg.FFmpegPath, cfg.InputURL, cfg.RTSPTransport, hub)

	registry := services.NewRegistry()

	go runner.RunLoop(ctx)

	app := server.New(cfg, hub, registry)
	errs := make(chan error, 1)

	go func() {
		log.Printf("SmartTraffic gateway listening on %s", cfg.Addr)
		errs <- app.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := app.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	case err := <-errs:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

func handleServiceAction(action string) error {
	switch action {
	case "install":
		return installService()
	case "uninstall":
		return uninstallService()
	case "start":
		return startService()
	case "stop":
		return controlService(svc.Stop, "stop")
	default:
		return fmt.Errorf("unknown service action %q", action)
	}
}

func installService() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	manager, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer manager.Disconnect()

	existing, err := manager.OpenService(serviceName)
	if err == nil {
		existing.Close()
		return fmt.Errorf("%s is already installed", serviceName)
	}

	service, err := manager.CreateService(serviceName, exePath, mgr.Config{
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		StartType:   mgr.StartAutomatic,
	})
	if err != nil {
		return err
	}
	defer service.Close()

	_ = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	log.Printf("installed %s from %s", serviceName, exePath)
	return nil
}

func uninstallService() error {
	manager, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()

	_, _ = service.Control(svc.Stop)
	if err := service.Delete(); err != nil {
		return err
	}

	_ = eventlog.Remove(serviceName)
	log.Printf("uninstalled %s", serviceName)
	return nil
}

func startService() error {
	manager, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()

	if err := service.Start(); err != nil {
		return err
	}
	log.Printf("start requested for %s", serviceName)
	return nil
}

func controlService(command svc.Cmd, label string) error {
	manager, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer manager.Disconnect()

	service, err := manager.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()

	_, err = service.Control(command)
	if err != nil {
		return err
	}
	log.Printf("%s requested for %s", label, serviceName)
	return nil
}
