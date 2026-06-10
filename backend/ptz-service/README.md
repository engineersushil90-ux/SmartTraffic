# SmartTraffic PTZ Service

Standalone PTZ Camera server for SmartTraffic.

The parent `Smarttraffic-Service` starts this process. The gateway calls it as a client/proxy.

## Endpoints

```text
GET  /healthz
GET  /api/ptz-cameras
GET  /api/ptz-cameras/{cameraId}
POST /api/ptz/{cameraId}
```

Default address:

```text
http://localhost:8092
```

## Run In Console

```powershell
cd backend/ptz-service
go run ./cmd/ptz-service
```

## Build

```powershell
cd backend/ptz-service
go build -o ptz-service.exe ./cmd/ptz-service
```
