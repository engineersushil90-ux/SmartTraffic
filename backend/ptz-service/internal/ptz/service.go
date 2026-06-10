package ptz

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusConnected    Status = "connected"
	StatusDisconnected Status = "disconnected"
)

type Camera struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Status   Status            `json:"status"`
	LastSeen string            `json:"lastSeen"`
	Details  map[string]string `json:"details"`
}

type CommandResult struct {
	OK        bool      `json:"ok"`
	CameraID  string    `json:"cameraId"`
	Command   string    `json:"command"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type Service struct {
	cameras []Camera
}

func NewService() *Service {
	return &Service{
		cameras: []Camera{
			{ID: "ptz-rohini-01", Name: "PTZ Rohini", Location: "NH 44 - Rohini", Status: StatusConnected, LastSeen: "live", Details: map[string]string{"pan": "enabled", "tilt": "enabled", "zoom": "enabled"}},
			{ID: "ptz-wazirpur-01", Name: "PTZ Wazirpur", Location: "Ring Road - Wazirpur", Status: StatusDisconnected, LastSeen: "10:21 AM", Details: map[string]string{"pan": "enabled", "tilt": "enabled", "zoom": "enabled"}},
		},
	}
}

func (s *Service) ListCameras() []Camera {
	cameras := make([]Camera, len(s.cameras))
	copy(cameras, s.cameras)
	return cameras
}

func (s *Service) GetCamera(id string) (Camera, bool) {
	for _, camera := range s.cameras {
		if strings.EqualFold(camera.ID, id) {
			return camera, true
		}
	}
	return Camera{}, false
}

func (s *Service) SendCommand(cameraID string, command string) (CommandResult, error) {
	camera, ok := s.GetCamera(cameraID)
	if !ok {
		return CommandResult{}, errors.New("camera not found")
	}
	if camera.Status != StatusConnected {
		return CommandResult{}, fmt.Errorf("camera %s is %s", cameraID, camera.Status)
	}
	command = strings.TrimSpace(command)
	if command == "" {
		return CommandResult{}, errors.New("command is required")
	}
	return CommandResult{OK: true, CameraID: camera.ID, Command: command, Status: "accepted", Timestamp: time.Now().UTC()}, nil
}

func (s *Service) Summary() map[string]any {
	status := map[string]int{"connected": 0, "disconnected": 0}
	for _, camera := range s.cameras {
		status[string(camera.Status)]++
	}
	return map[string]any{"service": "ptz", "total": len(s.cameras), "status": status}
}
