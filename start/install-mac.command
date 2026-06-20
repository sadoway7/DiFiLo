#!/bin/bash
# DIFI-LOCAL — INSTALL DEPENDENCIES & BUILD (macOS)
# Double-click this file to install Go dependencies and build the binary.
cd "$(dirname "$0")/.."

echo "========================================"
echo "  DiFiLo — Install & Build (macOS)"
echo "========================================"
echo ""

# Check for Go
if ! command -v go &> /dev/null; then
  echo "ERROR: Go is not installed."
  echo ""
  echo "Install Go from: https://go.dev/dl/"
  echo "Download the macOS installer (.pkg), run it, then re-run this script."
  echo ""
  echo "Press Return to close."; read -r x; exit 1
fi

GO_VERSION=$(go version)
echo "Found: $GO_VERSION"
echo ""

# Download dependencies
echo "Downloading dependencies..."
go mod download
if [ $? -ne 0 ]; then
  echo "ERROR: Failed to download dependencies."
  echo "Press Return to close."; read -r x; exit 1
fi
echo "Dependencies installed."
echo ""

# Build
echo "Building DiFiLo binary..."
go build -o DiFiLo ./cmd/difilo
if [ $? -ne 0 ]; then
  echo "ERROR: Build failed."
  echo "Press Return to close."; read -r x; exit 1
fi
echo ""

echo "========================================"
echo "  BUILD SUCCESSFUL!"
echo "========================================"
echo ""
echo "The DiFiLo binary has been built."
echo ""
echo "To start the server, double-click:"
echo "  start/start-mac.command"
echo ""
echo "Press Return to close."; read -r x
