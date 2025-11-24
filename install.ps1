# Maajise Installation Script for Windows
# Run as Administrator: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"

Write-Host "`n=== Maajise Installation ===" -ForegroundColor Cyan
Write-Host ""

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "ERROR: Please run as Administrator" -ForegroundColor Red
    Write-Host "Right-click PowerShell and select 'Run as Administrator'" -ForegroundColor Yellow
    exit 1
}

# Get current directory
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$sourceExe = Join-Path $scriptDir "maajise.exe"

# Check if maajise.exe exists
if (-not (Test-Path $sourceExe)) {
    Write-Host "Building maajise.exe..." -ForegroundColor Yellow
    Push-Location $scriptDir
    & go build -o maajise.exe
    Pop-Location
    
    if (-not (Test-Path $sourceExe)) {
        Write-Host "ERROR: Build failed" -ForegroundColor Red
        exit 1
    }
}

# Install location
$installDir = "C:\Program Files\Maajise"
$installExe = Join-Path $installDir "maajise.exe"

Write-Host "Installing to: $installDir" -ForegroundColor Green

# Create install directory
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Write-Host "Created directory: $installDir" -ForegroundColor Green
}

# Copy executable
Copy-Item -Path $sourceExe -Destination $installExe -Force
Write-Host "Copied maajise.exe" -ForegroundColor Green

# Add to PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*$installDir*") {
    $newPath = "$currentPath;$installDir"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    Write-Host "Added to system PATH" -ForegroundColor Green
    Write-Host ""
    Write-Host "IMPORTANT: Restart terminal for PATH changes" -ForegroundColor Yellow
} else {
    Write-Host "Already in system PATH" -ForegroundColor Green
}

Write-Host ""
Write-Host "Installation complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Usage:" -ForegroundColor Cyan
Write-Host "  maajise <project-name>" -ForegroundColor White
Write-Host "  maajise --help" -ForegroundColor White
Write-Host ""
