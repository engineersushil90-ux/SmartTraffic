# Smarttraffic-Service

Parent Windows service for the SmartTraffic backend bundle.

This service starts and monitors backend server processes:

```text
ATCC Service      -> http://localhost:8091
PTZ Service       -> http://localhost:8092
Gateway Client    -> http://localhost:8080
```

The gateway works as a client/proxy. ATCC and PTZ run as server processes.

## Build

```powershell
cd backend/atcc-service
go build -o atcc-service.exe ./cmd/atcc-service

cd ../ptz-service
go build -o ptz-service.exe ./cmd/ptz-service

cd ../smarttraffic-gateway
go build -o smarttraffic-gateway.exe ./cmd/smarttraffic-gateway

cd ../smarttraffic-service
go build -o smarttraffic-service.exe ./cmd/smarttraffic-service
```

## Install Parent Service

Run PowerShell as Administrator:

```powershell
cd backend/smarttraffic-service
.\smarttraffic-service.exe -service install
.\smarttraffic-service.exe -service start
```

Installed Windows service:

```text
Smarttraffic-Service
```

## Health

```text
GET http://localhost:8079/healthz
```

## Uninstall

Run PowerShell as Administrator:

```powershell
cd backend/smarttraffic-service
.\smarttraffic-service.exe -service stop
.\smarttraffic-service.exe -service uninstall
```
