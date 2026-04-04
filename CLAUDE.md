# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Project Is

`pure-game-kit` is a modular 2D game engine written in Go, built on top of [raylib-go](https://github.com/gen2brain/raylib-go). It provides a layered set of packages covering graphics, GUI, input, geometry, motion, execution flow, and utilities. `main.go` is an interactive hub that demonstrates all features via a GUI menu launching named example functions in `examples/systems/`.

## Build Commands

There are no tests and no linter configuration.

## Architecture

### Package Layers

```
window/        – OS window creation and the main game loop; message pump into internal/
internal/      – Shared runtime state (time, asset cache, input state); used only by engine packages to avoid import cycles
graphics/      – Drawing objects (using assets but still lightweight) and primitives through a Camera directly to a window
gui/           – Hybrid data-driven + immediate-mode UI; layouts defined via chained builder calls that produce XML, parsed at runtime
input/         – Frame-based polling (not event callbacks): mouse/, keyboard/
geometry/      – Line-based primitives, Shape collision (ray casting), ShapeGrid, A* pathfinding, path following
motion/        – Tween (chainable value animation), Animation[T] (generic frame sequences), ParticleSystem
execution/     – StateMachine, Command/Condition combinators, Script (Lua via gopher-lua), Screens state machine
data/          – Asset loading (fonts, textures, tilesets, tilemaps, sound, music); embedded defaults; file/folder/storage helpers
utility/       – Pure helper packages: angle, color, collection, direction, flag, is, naming, noise, number, point, random, text, time
```

### Key Design Decisions

**Minimizing structs** Keeping API shallow and wide rather than thin and deep. No multi-purpose vectors etc. Prefer working with a few values rather than a small struct. Big structs are fine but nesting is best avoided.

**`internal/` as the messenger** Engine packages (window, graphics, gui, input, motion) share state through `internal/` rather than importing each other, preventing dependency cycles. Do not import `internal/` from user-facing packages.

**`var` instead of `:=`** Wherever possible.

### Entry Point Flow

1. `main()` in `main.go` calls `window.Open(...)` which starts the raylib loop.
2. Each frame, `window` updates `internal.DeltaTime` and calls registered update/draw callbacks.
3. Example functions (in `examples/systems/`) each call `window.Open(...)` independently — only one runs at a time.
