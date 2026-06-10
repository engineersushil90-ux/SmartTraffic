$ErrorActionPreference = "Stop"

$backendRoot = Resolve-Path "$PSScriptRoot\.."
$atccRoot = Join-Path $backendRoot "atcc-service"
$gatewayRoot = Join-Path $backendRoot "smarttraffic-gateway"

Push-Location $atccRoot
try {
  go build -o atcc-service.exe .\cmd\atcc-service
}
finally {
  Pop-Location
}

Push-Location $gatewayRoot
try {
  go build -o smarttraffic-gateway.exe .\cmd\smarttraffic-gateway
}
finally {
  Pop-Location
}

Write-Host "Built SmartTraffic bundle services."
