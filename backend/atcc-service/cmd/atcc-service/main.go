package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"smarttraffic/atcc-service/internal/atcc"
	"smarttraffic/atcc-service/internal/config"
	"smarttraffic/atcc-service/internal/server"
)

const (
	serviceName        = "SmartTrafficATCCService"
	serviceDisplayName = "SmartTraffic ATCC Service"
	serviceDescription = "SmartTraffic ATCC backend service for traffic classification and count data."
)

type atccWindowsService struct{}

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
		if err := run(serviceName, atccWindowsService{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := runHTTP(context.Background()); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (m atccWindowsService) Execute(_ []string, requests <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	changes <- svc.Status{State: svc.StartPending}

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() {
		errs <- runHTTP(ctx)
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
				case <-time.After(10 * time.Second):
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

func runHTTP(ctx context.Context) error {
	cfg := config.Load()
	srv := server.New(cfg, atcc.NewService())
	errs := make(chan error, 1)

	go func() {
		log.Printf("ATCC service listening on %s", cfg.Addr)
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
