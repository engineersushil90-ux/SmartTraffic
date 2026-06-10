package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"smarttraffic/smarttraffic-service/internal/config"
	"smarttraffic/smarttraffic-service/internal/manager"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
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
			if err != nil {
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

	for _, status := range services.Statuses() {
		if status.Error != "" {
			log.Printf("service name=%s running=%t error=%s executable=%s", status.Name, status.Running, status.Error, status.Executable)
			continue
		}
		log.Printf("service name=%s running=%t pid=%d executable=%s", status.Name, status.Running, status.PID, status.Executable)
	}

	log.Printf("Smarttraffic-Service started all services")
	<-ctx.Done()
	log.Printf("Smarttraffic-Service stopping all services")
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
