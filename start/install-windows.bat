@echo off
REM DIFI-LOCAL - INSTALL DEPENDENCIES & BUILD (Windows)
REM Double-click this file to install Go dependencies and build the binary.
cd /d "%~dp0\.."

echo ========================================
echo   DiFiLo - Install ^& Build (Windows)
echo ========================================
echo.

REM Check for Go
where go >nul 2>&1
if %errorlevel% neq 0 (
  echo ERROR: Go is not installed.
  echo.
  echo Install Go from: https://go.dev/dl/
  echo Download the Windows installer (.msi), run it, then re-run this script.
  echo.
  pause
  exit /b 1
)

go version
echo.

REM Download dependencies
echo Downloading dependencies...
go mod download
if %errorlevel% neq 0 (
  echo ERROR: Failed to download dependencies.
  pause
  exit /b 1
)
echo Dependencies installed.
echo.

REM Build
echo Building DiFiLo binary...
go build -o DiFiLo.exe ./cmd/difilo
if %errorlevel% neq 0 (
  echo ERROR: Build failed.
  pause
  exit /b 1
)
echo.

echo ========================================
echo   BUILD SUCCESSFUL!
echo ========================================
echo.
echo The DiFiLo.exe binary has been built.
echo.
echo To start the server, double-click:
echo   scripts\start-windows.bat
echo.
pause
