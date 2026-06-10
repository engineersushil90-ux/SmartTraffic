# SmartTraffic Backend Bundle

SmartTraffic backend is a bundle of service executables managed by the SmartTraffic Gateway.

Install and start the gateway service, and it starts/connects the bundled services such as ATCC.

## Build Bundle

```powershell
cd backend
.\scripts\build-bundle.ps1
```

## Install Gateway Service

Run PowerShell as Administrator:

```powershell
cd backend
.\scripts\install-gateway-service.ps1
```

Installed Windows service:

```text
SmartTrafficGatewayService
```

The gateway starts:

```text
backend/atcc-service/atcc-service.exe
```

## Uninstall Gateway Service

Run PowerShell as Administrator:

```powershell
cd backend
.\scripts\uninstall-gateway-service.ps1
```

## Health

```text
GET http://localhost:8080/healthz
```

The health response includes `managedServices`, showing which services the gateway started and whether they are reachable.
