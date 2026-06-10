# SmartTraffic Gateway

Go backend gateway for the SmartTraffic service bundle and live camera streaming.

It starts and connects bundled SmartTraffic backend services, then exposes them through one gateway API. It also reads an RTSP stream with FFmpeg, keeps a rolling byte buffer, and serves an HTTP-FLV stream for the Angular dashboard.

## Structure

```text
backend/smarttraffic-gateway/
  cmd/smarttraffic-gateway/  App entrypoint
  internal/config/           Environment config
  internal/server/           HTTP routes, CORS, health, PTZ stub
  internal/services/         Separate ATCC, VIDS, PTZ, CCTV, MET, VMS, VSDS services
  internal/stream/           FFmpeg runner and buffered stream hub
  internal/supervisor/       Starts and monitors bundled service processes
  go.mod
```

## Requirements

- Go 1.22+
- FFmpeg available on `PATH`
- RTSP input running, for example:

```text
rtsp://localhost:8554/webcam
```

## Run

```powershell
cd backend/smarttraffic-gateway
go run ./cmd/smarttraffic-gateway
```

Default output:

```text
http://localhost:8080/live
```

The dashboard first feed is configured for this URL with `streamType: 'flv'`.

## Config

Environment variables:

```text
STREAM_GATEWAY_ADDR=:8080
STREAM_INPUT_RTSP=rtsp://localhost:8554/webcam
STREAM_RTSP_TRANSPORT=tcp
STREAM_BUFFER_BYTES=4194304
FFMPEG_PATH=ffmpeg
ATCC_SERVICE_URL=http://localhost:8091
ATCC_SERVICE_EXE=..\atcc-service\atcc-service.exe
SMARTTRAFFIC_MANAGE_SERVICES=true
```

Example:

```powershell
$env:STREAM_INPUT_RTSP='rtsp://localhost:8554/webcam'
$env:STREAM_GATEWAY_ADDR=':8080'
go run ./cmd/smarttraffic-gateway
```

## Endpoints

```text
GET  /healthz
GET  /live
GET  /api/services
GET  /api/atcc
GET  /api/atcc/{deviceId}
GET  /api/vids
GET  /api/vids/{deviceId}
GET  /api/ptz-cameras
GET  /api/ptz-cameras/{cameraId}
GET  /api/cctv-cameras
GET  /api/cctv-cameras/{cameraId}
GET  /api/met
GET  /api/met/{deviceId}
GET  /api/vms
GET  /api/vms/{deviceId}
GET  /api/vsds
GET  /api/vsds/{deviceId}
POST /api/ptz/{cameraId}
```

`/api/services` returns status summaries for every backend service.

Each device endpoint returns a service-owned JSON payload with `id`, `name`, `category`, `location`, `status`, `lastSeen`, optional `streamUrl`, and service-specific `details`.

ATCC routes are proxied to the standalone ATCC service configured by `ATCC_SERVICE_URL`.

When `SMARTTRAFFIC_MANAGE_SERVICES=true`, the gateway starts the bundled ATCC service executable configured by `ATCC_SERVICE_EXE`. `/healthz` includes `managedServices` so you can see whether each bundled service is running.

`/api/ptz/{cameraId}` validates the camera through the PTZ Camera service before accepting a command. Wire `internal/services` to your database, ONVIF service, or camera vendor APIs when ready.

## Install Gateway As Windows Service

Run PowerShell as Administrator:

```powershell
cd backend
.\scripts\install-gateway-service.ps1
```

This builds:

```text
backend/atcc-service/atcc-service.exe
backend/smarttraffic-gateway/smarttraffic-gateway.exe
```

Then it installs only the SmartTraffic Gateway Windows service:

```text
SmartTrafficGatewayService
```

The gateway starts the ATCC service process automatically.

To uninstall:

```powershell
cd backend
.\scripts\uninstall-gateway-service.ps1
```
