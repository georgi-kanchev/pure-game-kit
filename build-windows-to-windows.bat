@echo off
setlocal enabledelayedexpansion

:: --- CONFIGURATION ---
set "OUTPUT_EXE=app.exe"
set "GO_PROGRAM_DIR=%cd%"
set "WORKDIR=%TEMP%\raylib_build_%RANDOM%"

echo Using temporary workdir: %WORKDIR%

:: --- 1. CHECK FOR TOOLS ---
where git >nul 2>nul || (echo ERROR: git not found. Install Git for Windows. && pause && exit /b 1)
where go >nul 2>nul || (echo ERROR: go not found. Install Go. && pause && exit /b 1)
where cmake >nul 2>nul || (echo ERROR: cmake not found. Install CMake. && pause && exit /b 1)
where mingw32-make >nul 2>nul || (echo ERROR: mingw32-make not found. Ensure MinGW/bin is in your PATH. && pause && exit /b 1)

:: --- 2. CLONE AND BUILD RAYLIB ---
echo Cloning Raylib...
git clone --depth 1 https://github.com/raysan5/raylib.git "%WORKDIR%\raylib" || (echo Clone failed. && pause && exit /b 1)

echo Building Raylib static library...
mkdir "%WORKDIR%\raylib\build"
cd /d "%WORKDIR%\raylib\build"

cmake -G "MinGW Makefiles" ^
      -DCMAKE_BUILD_TYPE=Release ^
      -DBUILD_SHARED_LIBS=OFF ^
      -DPLATFORM=Desktop ^
      .. || (echo CMake configuration failed. && pause && exit /b 1)

mingw32-make || (echo Raylib compilation failed. && pause && exit /b 1)

:: Set paths (Converting backslashes to forward slashes for CGO)
set "RAYLIB_INCLUDE=%WORKDIR%\raylib\src"
set "RAYLIB_LIB=%WORKDIR%\raylib\build\raylib\libraylib.a"

set "RAYLIB_INCLUDE_CGO=%RAYLIB_INCLUDE:\=/%"
set "RAYLIB_LIB_CGO=%RAYLIB_LIB:\=/%"

:: --- 3. INJECT CGO FLAGS ---
cd /d "%GO_PROGRAM_DIR%"
set "TEMP_GO_FILE=.tmp_build.go"

echo package main > "%TEMP_GO_FILE%"
echo // #cgo CFLAGS: -I%RAYLIB_INCLUDE_CGO% >> "%TEMP_GO_FILE%"
echo // #cgo LDFLAGS: %RAYLIB_LIB_CGO% -lopengl32 -lgdi32 -lwinmm -lshell32 >> "%TEMP_GO_FILE%"
echo import "C" >> "%TEMP_GO_FILE%"

:: --- 4. BUILD GO PROGRAM ---
echo Building Go program...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

:: Using "go build -v" to see exactly what is happening
go build -v -ldflags="-s -w -H=windowsgui" -o "%OUTPUT_EXE%" || (echo Go build failed. Check the errors above. && pause && exit /b 1)

:: --- 5. CLEAN UP ---
if exist "%TEMP_GO_FILE%" del "%TEMP_GO_FILE%"
:: Note: I commented out workdir deletion so you can inspect it if it fails
:: if exist "%WORKDIR%" rd /s /q "%WORKDIR%"

echo ---------------------------------------
echo SUCCESS! Build complete: %GO_PROGRAM_DIR%\%OUTPUT_EXE%
pause