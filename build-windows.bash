#!/bin/bash
set -e

# --- CONFIGURATION ---
GOOS=windows
GOARCH=amd64
CC=x86_64-w64-mingw32-gcc
GO_PROGRAM_DIR=$(pwd)  # Assumes script is in your Go project root
GO_FILE="main.go"      # Change if your main file has a different name
OUTPUT_EXE=myprogram.exe

# --- INSTALL DEPENDENCIES ---
echo "Installing dependencies..."
sudo apt update
sudo apt install -y mingw-w64 cmake make git

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
      -DBUILD_SHARED_LIBS=OFF \
      ..
make

RAYLIB_INCLUDE="$WORKDIR/raylib/src"
RAYLIB_LIB="$WORKDIR/raylib/build-windows/libraylib.a"

# --- TEMPORARILY INJECT CGO FLAGS ---
cd $GO_PROGRAM_DIR
TEMP_GO_FILE=$(mktemp)
echo "// Code injected by cross-compile script" > $TEMP_GO_FILE
echo "// #cgo CFLAGS: -I$RAYLIB_INCLUDE" >> $TEMP_GO_FILE
echo "// #cgo LDFLAGS: $RAYLIB_LIB -lopengl32 -lgdi32 -lwinmm -lshell32" >> $TEMP_GO_FILE
echo "import \"C\"" >> $TEMP_GO_FILE
echo >> $TEMP_GO_FILE
cat $GO_FILE >> $TEMP_GO_FILE

# --- BUILD GO PROGRAM ---
echo "Building Go program for Windows..."
export GOOS=$GOOS
export GOARCH=$GOARCH
export CC=$CC

go build -v -o $OUTPUT_EXE $TEMP_GO_FILE

# --- CLEAN UP ---
rm $TEMP_GO_FILE
rm -rf $WORKDIR

echo "Build complete: $GO_PROGRAM_DIR/$OUTPUT_EXE"
