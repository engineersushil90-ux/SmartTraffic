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

	"smarttraffic/smarttraffic-service/internal/config"
	"smarttraffic/smarttraffic-service/internal/manager"
	"smarttraffic/smarttraffic-service/internal/server"
)

const (
	serviceName        = "Smarttraffic-Service"
	serviceDisplayName = "Smarttraffic Service"
	serviceDescription = "Smarttraffic parent service that starts gateway, ATCC, PTZ, and other backend services."
)

type windowsService struct{}

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
		if err := run(serviceName, windowsService{}); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := run(context.Background()); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (s windowsService) Execute(_ []string, requests <-chan svc.ChangeRequest, changes chan<- svc.Status) (bool, uint32) {
	changes <- svc.Status{State: svc.StartPending}
	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() { errs <- run(ctx) }()
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
			}
		case err := <-errs:
			if err != nil && err != http.ErrServerClosed {
				return false, 1
			}
			return false, 0
		}
	}
}

func run(ctx context.Context) error {
	cfg := config.Load()
	services := manager.New([]manager.Spec{
		{Name: "atcc", URL: cfg.ATCCURL, HealthURL: cfg.ATCCURL + "/healthz", Executable: cfg.ATCCExecutable},
		{Name: "ptz", URL: cfg.PTZURL, HealthURL: cfg.PTZURL + "/healthz", Executable: cfg.PTZExecutable},
		{Name: "gateway", URL: cfg.GatewayURL, HealthURL: cfg.GatewayURL + "/healthz", Executable: cfg.GatewayExecutable, Env: []string{"ATCC_SERVICE_URL=" + cfg.ATCCURL, "PTZ_SERVICE_URL=" + cfg.PTZURL}},
	})
	services.StartAll(ctx)
	defer services.StopAll()

	srv := server.New(cfg, services)
	errs := make(chan error, 1)
	go func() {
		log.Printf("Smarttraffic-Service listening on %s", cfg.Addr)
		errs <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
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
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	existing, err := m.OpenService(serviceName)
	if err == nil {
		existing.Close()
		return fmt.Errorf("%s is already installed", serviceName)
	}
	service, err := m.CreateService(serviceName, exePath, mgr.Config{DisplayName: serviceDisplayName, Description: serviceDescription, StartType: mgr.StartAutomatic})
	if err != nil {
		return err
	}
	defer service.Close()
	_ = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	return nil
}

func uninstallService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	service, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()
	_, _ = service.Control(svc.Stop)
	if err := service.Delete(); err != nil {
		return err
	}
	_ = eventlog.Remove(serviceName)
	return nil
}

func startService() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	service, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()
	return service.Start()
}

func controlService(command svc.Cmd, label string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	service, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer service.Close()
	_, err = service.Control(command)
	return err
}
