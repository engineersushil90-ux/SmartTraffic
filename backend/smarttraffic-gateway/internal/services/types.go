package services

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type DeviceStatus string

const (
	StatusConnected    DeviceStatus = "connected"
	StatusDisconnected DeviceStatus = "disconnected"
	StatusWarning      DeviceStatus = "warning"
)

type Device struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Category  string            `json:"category"`
	Location  string            `json:"location"`
	Status    DeviceStatus      `json:"status"`
	LastSeen  string            `json:"lastSeen"`
	StreamURL string            `json:"streamUrl,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

type DeviceService struct {
	category string
	devices  []Device
}

func NewDeviceService(category string, devices []Device) *DeviceService {
	return &DeviceService{category: category, devices: devices}
}

func (s *DeviceService) List() []Device {
	devices := make([]Device, len(s.devices))
	copy(devices, s.devices)
	return devices
}

func (s *DeviceService) Get(id string) (Device, bool) {
	for _, device := range s.devices {
		if strings.EqualFold(device.ID, id) {
			return device, true
		}
	}

	return Device{}, false
}

func (s *DeviceService) Summary() map[string]any {
	summary := map[string]int{
		"connected":    0,
		"disconnected": 0,
		"warning":      0,
	}

	for _, device := range s.devices {
		summary[string(device.Status)]++
	}

	return map[string]any{
		"category": s.category,
		"total":    len(s.devices),
		"status":   summary,
	}
}

type PTZCommandResult struct {
	OK        bool      `json:"ok"`
	CameraID  string    `json:"cameraId"`
	Command   string    `json:"command"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type PTZService struct {
	*DeviceService
}

func NewPTZService(devices []Device) *PTZService {
	return &PTZService{DeviceService: NewDeviceService("ptz-cameras", devices)}
}

func (s *PTZService) SendCommand(cameraID string, command string) (PTZCommandResult, error) {
	camera, ok := s.Get(cameraID)
	if !ok {
		return PTZCommandResult{}, errors.New("camera not found")
	}

	if camera.Status != StatusConnected {
		return PTZCommandResult{}, fmt.Errorf("camera %s is %s", cameraID, camera.Status)
	}

	command = strings.TrimSpace(command)
	if command == "" {
		return PTZCommandResult{}, errors.New("command is required")
	}

	return PTZCommandResult{
		OK:        true,
		CameraID:  camera.ID,
		Command:   command,
		Status:    "accepted",
		Timestamp: time.Now().UTC(),
	}, nil
}

type Registry struct {
	ATCC        *DeviceService
	VIDS        *DeviceService
	PTZCameras  *PTZService
	CCTVCameras *DeviceService
	MET         *DeviceService
	VMS         *DeviceService
	VSDS        *DeviceService
}

func NewRegistry() *Registry {
	return &Registry{
		ATCC:        NewATCCService(),
		VIDS:        NewVIDSService(),
		PTZCameras:  NewPTZCameraService(),
		CCTVCameras: NewCCTVCameraService(),
		MET:         NewMETService(),
		VMS:         NewVMSService(),
		VSDS:        NewVSDSService(),
	}
}

func (r *Registry) Summaries() []map[string]any {
	return []map[string]any{
		r.ATCC.Summary(),
		r.VIDS.Summary(),
		r.PTZCameras.Summary(),
		r.CCTVCameras.Summary(),
		r.MET.Summary(),
		r.VMS.Summary(),
		r.VSDS.Summary(),
	}
}
