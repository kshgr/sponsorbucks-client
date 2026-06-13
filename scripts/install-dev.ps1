$ErrorActionPreference = "Stop"

go build -o sponsorbucks.exe ./cmd/sponsorbucks

$target = Join-Path $HOME ".sponsorbucks\bin"
New-Item -ItemType Directory -Force -Path $target | Out-Null
Move-Item -Force sponsorbucks.exe (Join-Path $target "sponsorbucks.exe")

Write-Host "Installed sponsorbucks to $target\sponsorbucks.exe"
Write-Host "Add $target to your PATH."
