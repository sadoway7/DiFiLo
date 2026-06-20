#!/bin/bash
# [D]iFiLo — STOP (macOS)
# Double-click this file to stop the running [D]iFiLo server.
cd "$(dirname "$0")/.."

STOPPED=0
if [ -f difilo.pid ]; then
  PID=$(cat difilo.pid)
  if kill -0 "$PID" 2>/dev/null; then
    kill "$PID"
    echo "Stopped [D]iFiLo (PID $PID)."
    STOPPED=1
  fi
  rm -f difilo.pid
fi
if [ "$STOPPED" = "0" ]; then
  if pkill -f "[D]IFI-LOCAL --mirror"; then
    echo "Stopped [D]iFiLo."
  else
    echo "[D]iFiLo was not running."
  fi
fi
echo "Press Return to close."; read -r x
