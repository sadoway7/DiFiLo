@echo off
REM DIFI-LOCAL - STOP (Windows)
REM Double-click this file to stop the running DIFI-LOCAL server.

taskkill /im DiFiLo.exe /f >nul 2>&1
if %errorlevel%==0 (
  echo Stopped DIFI-LOCAL.
) else (
  echo DIFI-LOCAL was not running.
)
pause
