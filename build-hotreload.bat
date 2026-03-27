@echo off
REM filepath: f:\dev\godxc-v12-wrk\dxgoc\build-hotreload.bat
REM ============================================================
REM   Hot Reload Demo - Pure Go Build (NO CGO, NO C COMPILER)
REM ============================================================

echo.
echo ============================================================
echo   Hot Reload Demo - Pure Go Build (No CGO Required!)
echo ============================================================
echo.

set "PROJECT_DIR=%~dp0"
cd /d "%PROJECT_DIR%"

set "BIN_DIR=%PROJECT_DIR%bin"

REM ============================================
REM Build Pure Go Version (No CGO!)
REM ============================================

echo [1/2] Building hot reload demo (no CGO)...

REM Disable CGO completely
set CGO_ENABLED=0
REM Disable VCS stamping (prevents git calls)
set GOFLAGS=-mod=mod

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"

REM Build the pure Go version
go build -v -ldflags="-s -w" -trimpath -buildvcs=false -o "%BIN_DIR%\hotreload-demo.exe" .\cmd\hotreload-demo

if %ERRORLEVEL% neq 0 (
    echo.
    echo BUILD FAILED!
    echo.
    echo Make sure Go is installed: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo   Built: %BIN_DIR%\hotreload-demo.exe
echo.

REM ============================================
REM Done
REM ============================================

echo [2/2] Build complete!

echo.
echo ============================================================
echo   BUILD SUCCESSFUL!
echo ============================================================
echo.
echo  Output: %BIN_DIR%\hotreload-demo.exe
echo.
echo  This is the hot reload demonstration tool.
echo.
echo  Usage:
echo    bin\hotreload-demo.exe                    (show help)
echo    bin\hotreload-demo.exe -shader shader.hlsl -watch
echo    bin\hotreload-demo.exe -dir shaders/ -watch
echo.

pause
