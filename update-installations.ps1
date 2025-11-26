# Maajise Installation Update Script
# Updates all Maajise installations to the latest version
# Run as Administrator: powershell -ExecutionPolicy Bypass -File update-installations.ps1

$ErrorActionPreference = "Continue"

Write-Host "`n=== Maajise Installation Update ===" -ForegroundColor Cyan
Write-Host ""

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "WARNING: Not running as Administrator" -ForegroundColor Yellow
    Write-Host "You may not be able to update all installations" -ForegroundColor Yellow
    Write-Host "Right-click PowerShell and select 'Run as Administrator' for full access" -ForegroundColor Yellow
    Write-Host ""
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

        # Get file info
        $currentFile = Get-Item $installExe
        $newFile = Get-Item $sourceExe

        Write-Host "  Current: $($currentFile.LastWriteTime)" -ForegroundColor Gray
        Write-Host "  New:     $($newFile.LastWriteTime)" -ForegroundColor Gray

        # Copy executable with error handling
        try {
            Copy-Item -Path $sourceExe -Destination $installExe -Force -ErrorAction Stop
            Write-Host "  Updated successfully" -ForegroundColor Green
            $updated = $true
        } catch {
            Write-Host "  Failed to update: $_" -ForegroundColor Red
        }
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
        Write-Host "  $location" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "Note: If you have multiple installations, consider removing duplicates:" -ForegroundColor Yellow
Write-Host "  Remove-Item -Path '$($locations[1])' -Recurse -Force" -ForegroundColor Gray
Write-Host ""
Write-Host "To keep only one installation, run:" -ForegroundColor Cyan
Write-Host "  Remove-Item -Path 'C:\Users\$env:USERNAME\AppData\Local\Programs\Maajise' -Recurse -Force" -ForegroundColor White
Write-Host ""
