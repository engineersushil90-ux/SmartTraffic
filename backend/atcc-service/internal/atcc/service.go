package atcc

import (
	"strings"
	"time"
)

type Status string

const (
	StatusConnected    Status = "connected"
	StatusDisconnected Status = "disconnected"
	StatusWarning      Status = "warning"
)

type Device struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Location  string            `json:"location"`
	Status    Status            `json:"status"`
	LastSeen  string            `json:"lastSeen"`
	Details   map[string]string `json:"details"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type VehicleEvent struct {
	ID           int       `json:"id"`
	Junction     string    `json:"junction"`
	Timestamp    time.Time `json:"timestamp"`
	VehicleClass string    `json:"vehicleClass"`
	Direction    string    `json:"direction"`
	CameraID     string    `json:"cameraId"`
	Lane         int       `json:"lane"`
	Speed        string    `json:"speed"`
	Color        string    `json:"color"`
}

type Service struct {
	devices []Device
	events  []VehicleEvent
}

func NewService() *Service {
	now := time.Now().UTC()

	return &Service{
		devices: []Device{
			{
				ID:        "ATCC-ROH-01",
				Name:      "ATCC Rohini Mainline",
				Location:  "NH 44 - Rohini",
				Status:    StatusConnected,
				LastSeen:  "live",
				UpdatedAt: now,
				Details: map[string]string{
					"lanes":        "4",
					"vehicleCount": "36200 veh/hr",
					"avgSpeed":     "42 km/h",
				},
			},
			{
				ID:        "ATCC-DND-01",
				Name:      "ATCC DND Flyway",
				Location:  "NH 24 - DND Flyway",
				Status:    StatusWarning,
				LastSeen:  "2 minutes ago",
				UpdatedAt: now.Add(-2 * time.Minute),
				Details: map[string]string{
					"lanes":        "3",
					"vehicleCount": "18480 veh/hr",
					"avgSpeed":     "37 km/h",
				},
			},
		},
		events: []VehicleEvent{
			{ID: 17613, Junction: "ATCC-Test", Timestamp: now.Add(-1 * time.Second), VehicleClass: "Car", Direction: "Toward phil", CameraID: "192.168.2.97", Lane: 2, Speed: "-", Color: "-"},
			{ID: 17612, Junction: "ATCC-Test", Timestamp: now.Add(-2 * time.Second), VehicleClass: "Car", Direction: "Toward phil", CameraID: "192.168.2.97", Lane: 1, Speed: "-", Color: "-"},
			{ID: 17611, Junction: "ATCC-Test", Timestamp: now.Add(-3 * time.Second), VehicleClass: "Car", Direction: "Toward phil", CameraID: "192.168.2.97", Lane: 1, Speed: "-", Color: "-"},
			{ID: 17610, Junction: "ATCC-Test", Timestamp: now.Add(-4 * time.Second), VehicleClass: "Bus/Truck", Direction: "Toward phil", CameraID: "192.168.2.97", Lane: 2, Speed: "-", Color: "-"},
		},
	}
}

func (s *Service) ListDevices() []Device {
	devices := make([]Device, len(s.devices))
	copy(devices, s.devices)
	return devices
}

func (s *Service) GetDevice(id string) (Device, bool) {
	for _, device := range s.devices {
		if strings.EqualFold(device.ID, id) {
			return device, true
		}
	}

	return Device{}, false
}

func (s *Service) ListEvents() []VehicleEvent {
	events := make([]VehicleEvent, len(s.events))
	copy(events, s.events)
	return events
}

func (s *Service) Summary() map[string]any {
	status := map[string]int{
		"connected":    0,
		"disconnected": 0,
		"warning":      0,
	}

	for _, device := range s.devices {
		status[string(device.Status)]++
	}

	return map[string]any{
		"service": "atcc",
		"total":   len(s.devices),
		"status":  status,
	}
}
