@echo off
REM ============================================================
REM   Go - DXC Compiler - Pure Go Build (NO CGO, NO C COMPILER)
REM   This uses syscall.LoadDLL - works with just Go!
REM ============================================================

echo.
echo ============================================================
echo   Go - DXC Compiler - Pure Go Build (No CGO Required!)
echo ============================================================
echo.

set "PROJECT_DIR=%~dp0"
cd /d "%PROJECT_DIR%"

set "BIN_DIR=%PROJECT_DIR%bin"

REM ============================================
REM Check for dxcompiler.dll (optional)
REM ============================================

echo [1/3] Checking for dxcompiler.dll...

REM Check multiple locations
set "DLL_FOUND=0"

REM Check if dxc is in PATH
where dxc >nul 2>&1
if %ERRORLEVEL% equ 0 (
    for /f "tokens=*" %%i in ('where dxc') do set "DXC_PATH=%%i"
    for %%i in ("%DXC_PATH%") do set "DXC_DIR=%%~dpi"
    if exist "%DXC_DIR%dxcompiler.dll" (
        echo   Found via dxc command: %DXC_DIR%dxcompiler.dll
        set "DLL_FOUND=1"
    )
)

REM Check common installation paths
if %DLL_FOUND% equ 0 (
    if exist "C:\Program Files\Microsoft\DirectXShaderCompiler\bin\dxcompiler.dll" (
        echo   Found in Program Files
        set "DLL_FOUND=1"
    )
)

if %DLL_FOUND% equ 0 (
    echo.
    echo   dxcompiler.dll not found in common locations.
    echo   Build will succeed, but you need to:
    echo   - Place dxcompiler.dll in the same directory as dxgoc.exe, OR
    echo   - Add DXC's bin directory to your PATH, OR
    echo   - Use -dll flag to specify the path
    echo.
) else (
    echo   dxcompiler.dll is available
)

echo.

REM ============================================
REM Build Pure Go Version (No CGO!)
REM ============================================

echo [2/3] Building pure Go version (no CGO)...

REM Disable CGO completely
set CGO_ENABLED=0
REM Disable VCS stamping (prevents git calls)
set GOFLAGS=-mod=mod

if not exist "%BIN_DIR%" mkdir "%BIN_DIR%"

REM Build the pure Go version
go build -v -ldflags="-s -w" -trimpath -buildvcs=false -o "%BIN_DIR%\dxgoc.exe" .\cmd\dxgoc

if %ERRORLEVEL% neq 0 (
    echo.
    echo BUILD FAILED!
    echo.
    echo Make sure Go is installed: https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo   Built: %BIN_DIR%\dxgoc.exe
echo.

REM ============================================
REM Done - No DLL copying needed!
REM ============================================

echo [3/3] Build complete!

echo.
echo ============================================================
echo   BUILD SUCCESSFUL!
echo ============================================================
echo.
echo  Output: %BIN_DIR%\dxgoc.exe
echo.
echo  This build uses syscall.LoadDLL - NO C compiler needed!
echo.
echo  The program will automatically find dxcompiler.dll from:
echo   - Same directory as executable
echo   - PATH (where dxc command is found)
echo   - Common installation directories
echo   - Or use -dll "path\to\dxcompiler.dll" flag
echo.
echo  Usage:
echo    bin\dxgoc.exe -version
echo    bin\dxgoc.exe -check
echo    bin\dxgoc.exe -i shaders\basic_vertex.hlsl -e main -t vs_6_0
echo    bin\dxgoc.exe -dll "C:\path\to\dxcompiler.dll" -i shader.hlsl -e main -t vs_6_0
echo.

pause
