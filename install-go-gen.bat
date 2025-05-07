@echo off
setlocal

:: Set target install path
set "TARGET_DIR=C:\go-generator"
set "TARGET_EXE=go-gen.exe"

:: Create directory if it doesn't exist
if not exist "%TARGET_DIR%" (
    echo Creating %TARGET_DIR%...
    mkdir "%TARGET_DIR%"
)

:: Copy go-gen.exe to the target directory
echo Copying %TARGET_EXE% to %TARGET_DIR%...
copy /Y "%~dp0%TARGET_EXE%" "%TARGET_DIR%"

:: Check if TARGET_DIR is already in the PATH
echo Checking if %TARGET_DIR% is in system PATH...
set "FOUND="
for /f "tokens=2*" %%A in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>nul') do (
    echo %%B | find /I "%TARGET_DIR%" >nul
    if not errorlevel 1 (
        set "FOUND=1"
    )
)

:: If not found, add it to the system PATH
if not defined FOUND (
    echo Adding %TARGET_DIR% to system PATH...
    setx /M Path "%PATH%;%TARGET_DIR%"
) else (
    echo %TARGET_DIR% already in system PATH.
)

echo.
echo ✅ go-gen installed successfully! You can now run:
echo     go-gen init
echo.
pause
