// The entry point of the game and the pump of the game loop.
//
// Make sure to structure it like this:
//
//	// loading a file with settings...
//	engine.Initialize(...) // <- loaded settings
//	// any engine/game initializations or other loaded settings...
//	engine.Run(func { /* game loop code */ })
package engine

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/window"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Initialize(title string, tps, maxFps uint16, vsync, antialias bool) {
	var flags uint32 = rl.FlagWindowResizable
	if vsync {
		flags |= rl.FlagVsyncHint
	}
	if antialias {
		flags |= rl.FlagMsaa4xHint
	}

	rl.SetConfigFlags(flags)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.InitWindow(1280, 720, title)
	rl.SetExitKey(rl.KeyNull)
	rl.MaximizeWindow()
	rl.SetTargetFPS(int32(maxFps))
	window.MoveToMonitor(0)

	internal.TargetTPS = max(tps, 1)
	internal.Init()
}
func Run(gameLoop func()) {
	var ready = make(chan []internal.Batch, 1)
	var pool = make(chan []internal.Batch, 3)

	for range 3 {
		pool <- make([]internal.Batch, 0, 16)
	}

	go func() {
		var ticker = time.NewTicker(time.Second / time.Duration(internal.TargetTPS))
		var currentBatch = <-pool // grab the first buffer to start working

		for range ticker.C {
			if terminate {
				return
			}

			currentBatch = currentBatch[:0]
			internal.Batches = currentBatch // to avoid passing it into gameLoop()

			gameLoop()

			currentBatch = internal.Batches // update currentBatch to whatever gameLoop might have appended

			select {
			case ready <- currentBatch:
				currentBatch = <-pool // sent to render, now we need a new empty slice from the pool
			default: // renderer is busy and 'ready' is full
				var stale = <-ready
				pool <- stale // overwrite 'ready' by taking the old one out and putting it back in pool
				ready <- currentBatch
				currentBatch = <-pool
			}
		}
	}()

	//=================================================================

	var view = graphics.NewView(1)
	var activeBatches []internal.Batch

	for !rl.WindowShouldClose() {
		if terminate {
			return
		}

		select { // check if there is a new frame ready to draw
		case latest := <-ready:
			if activeBatches != nil {
				pool <- activeBatches // return the old batch to the pool for reuse
			}
			activeBatches = latest
		default: // keep drawing the last activeBatches if nothing new
		}

		internal.AccumulateInput()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		view.DrawTextDebug(true, false, false, false)

		for _, batch := range activeBatches {
			batch.Draw()
		}

		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func Stop() {
	terminate = true
}

// private ========================================================

var terminate bool
