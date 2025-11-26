# Maajise Installation Update Script
# Updates all Maajise installations to the latest version
# Run as Administrator: powershell -ExecutionPolicy Bypass -File update-installations.ps1

$ErrorActionPreference = "Stop"

Write-Host "`n=== Maajise Installation Update ===" -ForegroundColor Cyan
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
    Write-Host "ERROR: maajise.exe not found in $scriptDir" -ForegroundColor Red
    Write-Host "Please build the project first with: go build -o maajise.exe" -ForegroundColor Yellow
    exit 1
}

Write-Host "Source: $sourceExe" -ForegroundColor Green
Write-Host ""

# Installation locations
$locations = @(
    "C:\Program Files\Maajise",
    "C:\Users\$env:USERNAME\AppData\Local\Programs\Maajise"
)

$updated = $false

foreach ($location in $locations) {
    $installExe = Join-Path $location "maajise.exe"

    if (Test-Path $installExe) {
        Write-Host "Updating: $location" -ForegroundColor Yellow

        # Check version
        $currentVersion = & "$installExe" version 2>&1 | Select-String "v[0-9.]+" -o | Select-Object -First 1
        $newVersion = & "$sourceExe" version 2>&1 | Select-String "v[0-9.]+" -o | Select-Object -First 1

        Write-Host "  Current: $currentVersion" -ForegroundColor Gray
        Write-Host "  New:     $newVersion" -ForegroundColor Gray

        # Copy executable
        Copy-Item -Path $sourceExe -Destination $installExe -Force
        Write-Host "  ✓ Updated" -ForegroundColor Green
        $updated = $true
        Write-Host ""
    }
}

if (-not $updated) {
    Write-Host "No existing installations found" -ForegroundColor Yellow
    Write-Host "Running standard install..." -ForegroundColor Yellow
    & (Join-Path $scriptDir "install.ps1")
    exit 0
}

Write-Host "Update complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Current installations:" -ForegroundColor Cyan
foreach ($location in $locations) {
    $installExe = Join-Path $location "maajise.exe"
    if (Test-Path $installExe) {
        Write-Host "  ✓ $location" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "Note: If you have multiple installations, consider removing duplicates:" -ForegroundColor Yellow
Write-Host "  Remove-Item -Path '$($locations[1])' -Recurse -Force" -ForegroundColor Gray
Write-Host ""
Write-Host "To keep only one installation, run:" -ForegroundColor Cyan
Write-Host "  Remove-Item -Path 'C:\Users\$env:USERNAME\AppData\Local\Programs\Maajise' -Recurse -Force" -ForegroundColor White
Write-Host ""
