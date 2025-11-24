@echo off
REM Build script for Maajise
REM Run this from maajise directory

echo Building Maajise...
go build -o maajise.exe

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful: maajise.exe
    echo.
    echo To install globally, run:
    echo   powershell -ExecutionPolicy Bypass -File install.ps1
    echo.
) else (
    echo.
    echo Build failed. Make sure Go is installed.
    echo.
)

pause
