#!/bin/bash
# DIFI-LOCAL — START (macOS)
# Double-click this file to launch DIFI-LOCAL and open it in your browser.
cd "$(dirname "$0")/.."

if [ -f difilo.pid ] && kill -0 "$(cat difilo.pid)" 2>/dev/null; then
  echo "DIFI-LOCAL is already running (PID $(cat difilo.pid))."
else
  if [ ! -x ./DiFiLo ]; then
    echo "DiFiLo binary not found. Building..."
    if command -v go &> /dev/null; then
      go build -o DiFiLo ./cmd/difilo || { echo "Build failed. Press Return to close."; read -r x; exit 1; }
      echo "Build successful."
    else
      echo "Could not find DiFiLo binary and Go is not installed."
      echo "Press Return to close."; read -r x; exit 1
    fi
  fi
  nohup ./DiFiLo --mirror ./mirror --port 8000 > difilo.log 2>&1 &
  echo $! > difilo.pid
  echo "Started DIFI-LOCAL (PID $(cat difilo.pid))."
fi

sleep 1
echo "Waiting for DIFI-LOCAL to be ready (first run builds a search index, can take a minute or two)..."
for i in $(seq 1 180); do
  if curl -s -o /dev/null "http://localhost:8000/"; then
    echo "Ready!"
    break
  fi
  sleep 1
done
open "http://localhost:8000/"
echo ""
echo "DIFI-LOCAL is running at http://localhost:8000/"
echo "To stop it, double-click: scripts/stop-mac.command"
echo "You can close this window."
sleep 2
