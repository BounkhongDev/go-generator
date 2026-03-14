# Install go-gen-r for Windows: build and add this directory to user PATH.
# Requires PowerShell. Run: .\install.ps1

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot

Write-Host "Building go-gen-r..." -ForegroundColor Cyan
go build -o go-gen-r.exe ./cmd/go-gen-r
if (-not (Test-Path "go-gen-r.exe")) {
    Write-Error "Build failed."
    exit 1
}

$binDir = (Get-Location).Path
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -split ";" -notcontains $binDir) {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$binDir", "User")
    Write-Host "Added to user PATH: $binDir" -ForegroundColor Green
} else {
    Write-Host "Already in user PATH: $binDir" -ForegroundColor Yellow
}

Write-Host "Installed: $binDir\go-gen-r.exe" -ForegroundColor Green
Write-Host "Open a new terminal and run: go-gen-r" -ForegroundColor Cyan
