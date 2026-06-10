$ErrorActionPreference = "Stop"

$backendRoot = Resolve-Path "$PSScriptRoot\.."
$gatewayRoot = Join-Path $backendRoot "smarttraffic-gateway"
$gatewayExe = Join-Path $gatewayRoot "smarttraffic-gateway.exe"

if (Test-Path $gatewayExe) {
  Push-Location $gatewayRoot
  try {
    & $gatewayExe -service stop
    & $gatewayExe -service uninstall
  }
  finally {
    Pop-Location
  }
}
else {
  Write-Error "Cannot find $gatewayExe. Build the gateway first or remove SmartTrafficGatewayService manually."
}
