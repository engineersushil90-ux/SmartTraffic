# SmartTraffic Backend Bundle

SmartTraffic backend is a bundle managed by the parent `Smarttraffic-Service` Windows service.

`Smarttraffic-Service` starts backend server processes such as ATCC and PTZ, then starts the SmartTraffic Gateway. The gateway works as a client/proxy to ATCC and PTZ.

## Build Services

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

The parent service starts ATCC, PTZ, and Gateway.

## Uninstall Parent Service

Run PowerShell as Administrator:

```powershell
cd backend/smarttraffic-service
.\smarttraffic-service.exe -service stop
.\smarttraffic-service.exe -service uninstall
```

Use the gateway health endpoint to verify the running bundle:

```text
GET http://localhost:8080/healthz
```
