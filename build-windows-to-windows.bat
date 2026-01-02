@echo off
setlocal enabledelayedexpansion

:: --- CONFIGURATION ---
set "MSYS_ROOT=C:\msys64"
set "MINGW_PATH=%MSYS_ROOT%\mingw64"
set "PATH=%MINGW_PATH%\bin;%MSYS_ROOT%\usr\bin;%PATH%"
set "OUTPUT_EXE=app.exe"

echo Checking environment...

:: --- 1. CHECK FOR MSYS2 INSTALLATION ---
if not exist "%MSYS_ROOT%\msys2_shell.cmd" (
    echo [ERROR] MSYS2 not found at %MSYS_ROOT%
    echo Please install MSYS2 from https://www.msys2.org/ first.
    pause
    exit /b 1
)

:: --- 2. ENSURE GCC AND RAYLIB ARE INSTALLED ---
:: We use 'pacman' via the MSYS2 executable to install missing packages silently
where gcc >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [STATUS] GCC not found. Attempting to install via pacman...
    "%MSYS_ROOT%\usr\bin\bash" -lc "pacman -S --noconfirm mingw-w64-x86_64-toolchain"
) else (
    echo [OK] GCC is already installed.
)

:: Ensure Raylib library is also present in MSYS2
if not exist "%MINGW_PATH%\lib\libraylib.a" (
    echo [STATUS] Raylib library not found. Installing...
    "%MSYS_ROOT%\usr\bin\bash" -lc "pacman -S --noconfirm mingw-w64-x86_64-raylib"
) else (
    echo [OK] Raylib library found.
)

:: --- 3. CONFIGURE CGO VARIABLES ---
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

:: Pointing Go to the MSYS2 include and lib folders
set "CGO_CFLAGS=-I%MINGW_PATH%/include"
set "CGO_LDFLAGS=-L%MINGW_PATH%/lib -lraylib -lopengl32 -lgdi32 -lwinmm"

:: --- 4. BUILD GO PROGRAM ---
echo [STATUS] Building Go program: %OUTPUT_EXE%...
go build -v -ldflags="-s -w -H=windowsgui" -o "%OUTPUT_EXE%" .

if %ERRORLEVEL% EQU 0 (
    echo ---------------------------------------
    echo SUCCESS! Build complete.
) else (
    echo ---------------------------------------
    echo [ERROR] Build failed. Check the logs above.
)

pause