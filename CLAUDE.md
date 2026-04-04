# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Is

`pure-game-kit` is a modular 2D game engine written in Go, built on top of [raylib-go](https://github.com/gen2brain/raylib-go). It provides a layered set of packages covering graphics, GUI, input, geometry, motion, execution flow, and utilities. `main.go` is an interactive hub that demonstrates all features via a GUI menu launching named example functions in `examples/systems/`.

## Build Commands

```bash
# Run (Linux)
go run .

# Build for Linux
bash build-linux-to-linux.bash   # outputs ./app

# Build for Windows (cross-compile from Linux, requires mingw-w64 and cmake)
bash build-linux-to-windows.bash  # outputs ./app.exe

# Check heap/stack allocation for a file
go build -gcflags="-m" <package>/<file>.go
```

There are no tests (`*_test.go` files do not exist in this project) and no linter configuration.

## Architecture

### Package Layers

```
window/        – OS window creation and the main game loop; pumps DeltaTime into internal/
internal/      – Shared runtime state (time, asset cache, input state); used only by engine packages to avoid import cycles
graphics/      – Drawing: Camera, Sprite, NinePatch, TextBox, TileMap, blend modes
gui/           – Hybrid data-driven + immediate-mode UI; layouts defined via chained builder calls that produce XML, parsed at runtime
input/         – Frame-based polling (not event callbacks): mouse/, keyboard/
geometry/      – Line-based primitives, Shape collision (ray casting), ShapeGrid, A* pathfinding, path following
motion/        – Tween (chainable value animation), Animation[T] (generic frame sequences), ParticleSystem
execution/     – StateMachine, Command/Condition combinators, Script (Lua via gopher-lua), Screens
data/          – Asset loading (fonts, textures, tilesets, XML UI layouts); embedded defaults; file/folder/storage helpers
utility/       – Pure helper packages: angle, color, collection, direction, flag, is, naming, noise, number, point, random, text, time
```

### Key Design Decisions

**No pointers for value types.** Recent commits moved most structs off the heap. Prefer value receivers unless mutation across frames is needed. Use `go build -gcflags="-m"` to verify stack vs. heap allocation.

**`internal/` as the messenger.** Engine packages (window, graphics, gui, input, motion) share state through `internal/` rather than importing each other, preventing dependency cycles. Do not import `internal/` from user-facing packages.

**GUI is XML under the hood.** `gui.NewFromXMLs()` takes chained builder calls (`gui.Container(...)`, `gui.Button(...)`, `gui.Theme(...)`, etc.) that serialize to XML and are parsed at runtime. Dynamic values (camera dimensions, owner sizes) use string expressions from `gui/dynamic/` (e.g., `d.CameraCenterX+"-400"`). Field names are constants in `gui/field/`. GUI XML does not support the symbols `>` and `&` — use `internal.Placeholder` substitution instead.

**Line as the sole geometry primitive.** `geometry.Line` (`Ax, Ay, Bx, By`) is the base for all shapes. `Shape` is a polygon built from corner points; collision uses ray casting. `ShapeGrid` provides spatial chunking for queries.

**Generic animation.** `motion.Animation[T]` is a generic frame-sequence type. `motion.Tween` chains value animations with easing functions from `motion/easing/`.

**Assets are embedded.** Default fonts, UI textures, cursors, and patterns are embedded in `data/assets/` and lazy-loaded on first use.

### Entry Point Flow

1. `main()` in `main.go` calls `window.Open(...)` which starts the raylib loop.
2. Each frame, `window` updates `internal.DeltaTime` and calls registered update/draw callbacks.
3. Example functions (in `examples/systems/`) each call `window.Open(...)` independently — only one runs at a time.
