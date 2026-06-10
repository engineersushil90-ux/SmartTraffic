$ErrorActionPreference = "Stop"

$serviceRoot = Resolve-Path "$PSScriptRoot\.."
$exePath = Join-Path $serviceRoot "atcc-service.exe"

if (Test-Path $exePath) {
  & $exePath -service stop
  & $exePath -service uninstall
}
else {
  Write-Error "Cannot find $exePath. Build the service first or remove SmartTrafficATCCService manually."
}
