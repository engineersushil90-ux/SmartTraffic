$ErrorActionPreference = "Stop"

$serviceRoot = Resolve-Path "$PSScriptRoot\.."
$exePath = Join-Path $serviceRoot "atcc-service.exe"

Push-Location $serviceRoot
try {
  go build -o $exePath .\cmd\atcc-service
  & $exePath -service install
  & $exePath -service start
}
finally {
  Pop-Location
}
