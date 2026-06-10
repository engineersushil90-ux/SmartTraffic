# SmartTraffic Gateway

Gateway API for the SmartTraffic service bundle and live camera streaming.

The gateway works as a client/proxy to backend server services like ATCC and PTZ. It also reads an RTSP stream with FFmpeg, keeps a rolling byte buffer, and serves an HTTP-FLV stream for the Angular dashboard.

## Structure

```text
backend/smarttraffic-gateway/
  cmd/smarttraffic-gateway/  App entrypoint
  internal/config/           Environment config
  internal/server/           HTTP routes, CORS, health, PTZ stub
  internal/services/         Local fallback/static services for unsplit domains
  internal/stream/           FFmpeg runner and buffered stream hub
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
PTZ_SERVICE_URL=http://localhost:8092
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

ATCC and PTZ routes are proxied to their server services. The parent `Smarttraffic-Service` is responsible for starting those service processes.

`/api/ptz/{cameraId}` validates the camera through the PTZ Camera service before accepting a command. Wire `internal/services` to your database, ONVIF service, or camera vendor APIs when ready.

## Run Gateway Directly

Normally, `Smarttraffic-Service` starts the gateway. For development:

```powershell
cd backend/smarttraffic-gateway
go run ./cmd/smarttraffic-gateway
```
