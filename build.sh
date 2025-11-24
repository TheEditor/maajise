#!/bin/bash
# Build script for Maajise

echo "Building Maajise..."
go build -o maajise.exe

if [ $? -eq 0 ]; then
    echo ""
    echo "Build successful: maajise.exe"
    echo ""
    echo "To install globally (Windows):"
    echo "  powershell -ExecutionPolicy Bypass -File install.ps1"
    echo ""
else
    echo ""
    echo "Build failed. Make sure Go is installed."
    echo ""
    exit 1
fi
