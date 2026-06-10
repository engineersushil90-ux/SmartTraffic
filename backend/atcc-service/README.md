# SmartTraffic ATCC Service

Standalone ATCC backend service for traffic classification/count devices.

In the SmartTraffic bundle, this service is normally started by `SmartTrafficGatewayService`. You can still run it directly during development, or install it separately only when you intentionally want ATCC to be managed outside the gateway.

## Endpoints

```text
GET /healthz
GET /api/atcc
GET /api/atcc/{deviceId}
GET /api/atcc-events
```

Default address:

```text
http://localhost:8091
```

Configure with:

```text
ATCC_SERVICE_ADDR=:8091
ATCC_READ_HEADER_TIMEOUT_SECONDS=5
```

## Run In Console

```powershell
cd backend/atcc-service
go run ./cmd/atcc-service
```

## Build

```powershell
cd backend/atcc-service
go build -o atcc-service.exe ./cmd/atcc-service
```

## Install As Windows Service

Use this only if ATCC should run independently from the SmartTraffic Gateway bundle.

Run PowerShell as Administrator:

```powershell
cd backend/atcc-service
.\scripts\install-service.ps1
```

Manual equivalent:

```powershell
go build -o atcc-service.exe ./cmd/atcc-service
.\atcc-service.exe -service install
.\atcc-service.exe -service start
```

## Stop And Uninstall

Run PowerShell as Administrator:

```powershell
cd backend/atcc-service
.\scripts\uninstall-service.ps1
```

Manual equivalent:

```powershell
.\atcc-service.exe -service stop
.\atcc-service.exe -service uninstall
```
