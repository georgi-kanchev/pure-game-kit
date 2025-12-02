#!/bin/bash
set -e

# --- CONFIGURATION ---
GOOS=windows
GOARCH=amd64
CC=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++
GO_PROGRAM_DIR=$(pwd)
OUTPUT_EXE=app_debug.exe

# --- CREATE TEMPORARY WORKDIR ---
WORKDIR=$(mktemp -d)
echo "Using temporary workdir: $WORKDIR"

# --- CLONE AND BUILD RAYLIB ---
echo "Cloning Raylib..."
git clone --depth 1 https://github.com/raysan5/raylib.git $WORKDIR/raylib

echo "Building Raylib static library for Windows..."
mkdir -p $WORKDIR/raylib/build-windows
cd $WORKDIR/raylib/build-windows
cmake -G "Unix Makefiles" \
      -DCMAKE_SYSTEM_NAME=Windows \
      -DCMAKE_C_COMPILER=$CC \
      -DCMAKE_CXX_COMPILER=$CXX \
      -DBUILD_SHARED_LIBS=OFF \
      -DCMAKE_BUILD_TYPE=Debug \
      ..
make

RAYLIB_INCLUDE="$WORKDIR/raylib/src"
RAYLIB_LIB="$WORKDIR/raylib/build-windows/libraylib.a"

# --- TEMPORARILY INJECT CGO FLAGS ---
cd $GO_PROGRAM_DIR
TEMP_GO_FILE="$GO_PROGRAM_DIR/.tmp_build.go"

echo "package main" > $TEMP_GO_FILE
echo "// Code injected by cross-compile script" >> $TEMP_GO_FILE
echo "// #cgo CFLAGS: -I$RAYLIB_INCLUDE" >> $TEMP_GO_FILE
echo "// #cgo LDFLAGS: $RAYLIB_LIB -lopengl32 -lgdi32 -lwinmm -lshell32" >> $TEMP_GO_FILE
echo "import \"C\"" >> $TEMP_GO_FILE
echo >> $TEMP_GO_FILE

# Include all Go files except the temp file
find . -name "*.go" -type f ! -name ".tmp_build.go" -print0 | xargs -0 cat >> $TEMP_GO_FILE

# --- BUILD GO PROGRAM (DEBUG MODE) ---
echo "Building Go program for Windows (debug)..."
export GOOS=$GOOS
export GOARCH=$GOARCH
export CC=$CC
export CXX=$CXX

# Debug flags: disable optimizations and inlining
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build \
    -gcflags "all=-N -l" \
    -o $OUTPUT_EXE

# --- CLEAN UP ---
rm $TEMP_GO_FILE
rm -rf $WORKDIR

echo "Debug build complete: $GO_PROGRAM_DIR/$OUTPUT_EXE"
echo "Press Enter to exit..."
read

