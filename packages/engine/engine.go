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
	var ready = make(chan *internal.BatchManager, 1)
	var pool = make(chan *internal.BatchManager, 3)

	for range 3 { // pre-fill the pool with 3 managers (triple buffering)
		pool <- &internal.BatchManager{}
	}

	go func() {
		var ticker = time.NewTicker(time.Second / time.Duration(internal.TargetTPS))

		for range ticker.C {
			if terminate {
				return
			}

			internal.SyncAccumulatedInput()
			internal.UpdateWindowData()
			internal.UpdateTimeData()
			internal.UpdateMusic()
			internal.UpdateScreens()

			var manager = <-pool        // grab a manager from the pool
			manager.Reset()             // clear previous frame's data (resets slices to length 0, keeps capacity)
			internal.Renderer = manager // set as the global active renderer so gameLoop can call Queue
			gameLoop()
			manager.Finalize() // close out the final active batch inside the manager

			select {
			case ready <- manager: // successfully sent to the renderer (main thread)
			default:
				select { // renderer (main thread) is lagging - swap out the 'ready' frame for the newer one
				case stale := <-ready:
					pool <- stale
					ready <- manager
				default:
					ready <- manager
				}
			}
		}
	}()

	//=================================================================

	var activeManager *internal.BatchManager
	for !rl.WindowShouldClose() {
		if terminate {
			return
		}

		select { // check for a new frame from the ticker
		case latest := <-ready:
			if activeManager != nil {
				pool <- activeManager // return used manager to pool
			}
			activeManager = latest
		default: // no new frame? keep drawing activeManager
		}

		internal.AccumulateInput()

		rl.BeginDrawing()
		rl.EnableDepthTest()

		rl.ClearBackground(rl.Black)

		if activeManager != nil {
			activeManager.Draw()
		} // draw all batches stored in the manager for this frame

		rl.DrawFPS(10, 10)

		rl.DisableDepthTest()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func Stop() {
	terminate = true
}

// private ========================================================

var terminate bool
