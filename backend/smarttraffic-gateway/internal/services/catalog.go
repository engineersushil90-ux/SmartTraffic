package services

func NewATCCService() *DeviceService {
	return NewDeviceService("atcc", []Device{
		{
			ID:       "ATCC-ROH-01",
			Name:     "ATCC Rohini Mainline",
			Category: "ATCC",
			Location: "NH 44 - Rohini",
			Status:   StatusConnected,
			LastSeen: "live",
			Details: map[string]string{
				"lanes":        "4",
				"vehicleCount": "36200 veh/hr",
				"avgSpeed":     "42 km/h",
			},
		},
		{
			ID:       "ATCC-DND-01",
			Name:     "ATCC DND Flyway",
			Category: "ATCC",
			Location: "NH 24 - DND Flyway",
			Status:   StatusWarning,
			LastSeen: "2 minutes ago",
			Details: map[string]string{
				"lanes":        "3",
				"vehicleCount": "18480 veh/hr",
				"avgSpeed":     "37 km/h",
			},
		},
	})
}

func NewVIDSService() *DeviceService {
	return NewDeviceService("vids", []Device{
		{
			ID:        "VIDS-ROH-01",
			Name:      "VIDS Rohini",
			Category:  "VIDS",
			Location:  "NH 44 - Rohini",
			Status:    StatusConnected,
			LastSeen:  "live",
			StreamURL: "http://localhost:8080/live",
			Details: map[string]string{
				"analytics": "incident detection, stopped vehicle",
				"stream":    "flv",
			},
		},
		{
			ID:        "VIDS-WAZ-01",
			Name:      "VIDS Wazirpur",
			Category:  "VIDS",
			Location:  "Ring Road - Wazirpur",
			Status:    StatusDisconnected,
			LastSeen:  "stream unavailable",
			StreamURL: "/streams/ring-road-wazirpur/mjpeg",
			Details: map[string]string{
				"analytics": "wrong-way detection, queue detection",
				"stream":    "mjpeg",
			},
		},
	})
}

func NewPTZCameraService() *PTZService {
	return NewPTZService([]Device{
		{
			ID:       "ptz-rohini-01",
			Name:     "PTZ Rohini",
			Category: "PTZ Camera",
			Location: "NH 44 - Rohini",
			Status:   StatusConnected,
			LastSeen: "live",
			Details: map[string]string{
				"pan":  "enabled",
				"tilt": "enabled",
				"zoom": "enabled",
			},
		},
		{
			ID:       "ptz-wazirpur-01",
			Name:     "PTZ Wazirpur",
			Category: "PTZ Camera",
			Location: "Ring Road - Wazirpur",
			Status:   StatusDisconnected,
			LastSeen: "10:21 AM",
			Details: map[string]string{
				"pan":  "enabled",
				"tilt": "enabled",
				"zoom": "enabled",
			},
		},
	})
}

func NewCCTVCameraService() *DeviceService {
	return NewDeviceService("cctv-cameras", []Device{
		{
			ID:       "CCTV-ND-01",
			Name:     "CCTV New Delhi",
			Category: "CCTV Camera",
			Location: "New Delhi",
			Status:   StatusConnected,
			LastSeen: "live",
			Details: map[string]string{
				"resolution": "1080p",
				"purpose":    "traffic monitoring",
			},
		},
		{
			ID:       "CCTV-RJG-01",
			Name:     "CCTV Rajouri Garden",
			Category: "CCTV Camera",
			Location: "NH 48 - Rajouri Garden",
			Status:   StatusWarning,
			LastSeen: "5 minutes ago",
			Details: map[string]string{
				"resolution": "720p",
				"issue":      "intermittent packets",
			},
		},
	})
}

func NewMETService() *DeviceService {
	return NewDeviceService("met", []Device{
		{
			ID:       "MET-ND-01",
			Name:     "MET New Delhi",
			Category: "MET",
			Location: "New Delhi",
			Status:   StatusConnected,
			LastSeen: "live",
			Details: map[string]string{
				"temperature": "28 C",
				"humidity":    "54%",
				"wind":        "12 km/h NE",
				"visibility":  "6.5 km",
			},
		},
	})
}

func NewVMSService() *DeviceService {
	return NewDeviceService("vms", []Device{
		{
			ID:       "VMS-IIT-01",
			Name:     "VMS IIT Flyover",
			Category: "VMS",
			Location: "IIT Flyover",
			Status:   StatusConnected,
			LastSeen: "10:42 AM",
			Details: map[string]string{
				"message": "DRIVE SAFE / ARRIVE SAFE",
			},
		},
	})
}

func NewVSDSService() *DeviceService {
	return NewDeviceService("vsds", []Device{
		{
			ID:       "VSDS-NH48-01",
			Name:     "VSDS NH 48",
			Category: "VSDS",
			Location: "NH 48 - Rajouri Garden",
			Status:   StatusConnected,
			LastSeen: "live",
			Details: map[string]string{
				"averageSpeed": "95.7 km/h",
				"maxSpeed":     "213 km/h",
				"violations":   "36097",
			},
		},
	})
}
