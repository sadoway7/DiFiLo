@echo off
REM DIFI-LOCAL - START (Windows)
REM Double-click this file to launch DIFI-LOCAL and open it in your browser.
cd /d "%~dp0\.."

tasklist /fi "imagename eq DiFiLo.exe" 2>nul | find /i "DiFiLo.exe" >nul && (
  echo DIFI-LOCAL is already running.
) || (
  if not exist DiFiLo.exe (
    echo DiFiLo.exe not found. Building...
    where go >nul 2>&1
    if %errorlevel% equ 0 (
      go build -o DiFiLo.exe ./cmd/difilo
      if %errorlevel% neq 0 (
        echo Build failed.
        pause
        exit /b 1
      )
      echo Build successful.
    ) else (
      echo Could not find DiFiLo.exe and Go is not installed.
      pause
      exit /b 1
    )
  )
  start "" /b DiFiLo.exe --mirror .\mirror --port 8000 > difilo.log 2>&1
  echo Started DIFI-LOCAL.
)

echo Waiting for DIFI-LOCAL to be ready (first run builds a search index, can take a minute or two)...
for /L %%i in (1,1,180) do (
  powershell -NoProfile -Command "try{Invoke-WebRequest -UseBasicParsing -Uri 'http://localhost:8000/' -TimeoutSec 2 | Out-Null;exit 0}catch{exit 1}" >nul 2>&1 && goto :ready
  timeout /t 1 /nobreak >nul
)
:ready
start "" "http://localhost:8000/"
echo.
echo DIFI-LOCAL is running at http://localhost:8000/
echo To stop it, double-click: scripts\stop-windows.bat
echo You can close this window.
timeout /t 3 /nobreak >nul
