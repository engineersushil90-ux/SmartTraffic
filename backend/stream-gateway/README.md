# SmartTraffic Stream Gateway

Go backend for camera streaming.

It reads an RTSP stream with FFmpeg, keeps a rolling byte buffer, and serves an HTTP-FLV stream for the Angular dashboard.

## Structure

```text
backend/stream-gateway/
  cmd/stream-gateway/        App entrypoint
  internal/config/           Environment config
  internal/server/           HTTP routes, CORS, health, PTZ stub
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
cd backend/stream-gateway
go run ./cmd/stream-gateway
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
```

Example:

```powershell
$env:STREAM_INPUT_RTSP='rtsp://localhost:8554/webcam'
$env:STREAM_GATEWAY_ADDR=':8080'
go run ./cmd/stream-gateway
```

## Endpoints

```text
GET  /healthz
GET  /live
POST /api/ptz/{cameraId}
```

`/api/ptz/{cameraId}` is a stub for now. Wire it to your camera vendor API or ONVIF service when ready.
