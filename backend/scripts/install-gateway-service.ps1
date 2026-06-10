$ErrorActionPreference = "Stop"

$backendRoot = Resolve-Path "$PSScriptRoot\.."
$gatewayRoot = Join-Path $backendRoot "smarttraffic-gateway"
$gatewayExe = Join-Path $gatewayRoot "smarttraffic-gateway.exe"

& "$PSScriptRoot\build-bundle.ps1"

Push-Location $gatewayRoot
try {
  & $gatewayExe -service install
  & $gatewayExe -service start
}
finally {
  Pop-Location
}

Write-Host "SmartTraffic Gateway service installed and started."
