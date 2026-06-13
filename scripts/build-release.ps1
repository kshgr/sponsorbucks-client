param(
  [string]$Goos = "windows",
  [string]$Goarch = "amd64",
  [string]$OutDir = "dist",
  [string]$Version = "1.0.0-preview",
  [string]$BuildId = "manual",
  [string]$BuildChannel = "release"
)

$ErrorActionPreference = "Stop"
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
$ext = if ($Goos -eq "windows") { ".exe" } else { "" }
$outFile = Join-Path $OutDir "sponsorbucks-$Goos-$Goarch$ext"
$ldflags = "-s -w -X sponsorbucks-client/internal/buildinfo.Version=$Version -X sponsorbucks-client/internal/buildinfo.BuildID=$BuildId -X sponsorbucks-client/internal/buildinfo.BuildChannel=$BuildChannel"
$env:GOOS = $Goos
$env:GOARCH = $Goarch
go build -trimpath -ldflags $ldflags -o $outFile ./cmd/sponsorbucks
