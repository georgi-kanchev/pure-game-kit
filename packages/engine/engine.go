// The entry point of the game and the pump of the game loop.
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
	var ready = make(chan *internal.Bus, 1)
	var pool = make(chan *internal.Bus, 3)

	for range 3 { // pre-fill the pool with 3 buses (triple buffering)
		pool <- &internal.Bus{}
	}

	go func() {
		var ticker = time.NewTicker(time.Second / time.Duration(internal.TargetTPS))

		for range ticker.C {
			if terminate {
				return
			}

			internal.UpdateWindowData()
			internal.UpdateTimeData()
			internal.UpdateAudio()

			var bus = <-pool           // grab a bus from the pool (carries input snapshot from main thread)
			bus.Reset()                // clear previous frame's data (resets slices to length 0, keeps capacity)
			bus.SyncAccumulatedInput() // apply the input snapshot that the main thread captured onto the bus
			internal.ActiveBus = bus   // set as the global active bus so gameLoop can call Queue
			gameLoop()
			bus.Finalize() // close out the final active batch inside the bus

			select {
			case ready <- bus: // successfully sent to the renderer (main thread)
			default:
				select { // renderer (main thread) is lagging - swap out the 'ready' frame for the newer one
				case stale := <-ready:
					pool <- stale
					ready <- bus
				default:
					ready <- bus
				}
			}
		}
	}()

	//=================================================================

	var activeBus *internal.Bus
	var dirtyDraw = false
	for !rl.WindowShouldClose() {
		if terminate {
			return
		}

		activeBus.AccumulateInput()

		select { // check for a new frame from the ticker
		case latest := <-ready:
			if activeBus != nil {
				activeBus.CopyInputToBus() // snapshot accumulated input onto the bus before returning it
				pool <- activeBus          // return used bus to pool (carries input snapshot to ticker)
			}
			activeBus = latest
			dirtyDraw = true
		default: // no new frame? keep drawing active bus
			dirtyDraw = false
		}

		rl.BeginDrawing()
		rl.EnableDepthTest()

		rl.ClearBackground(rl.Black)

		if activeBus != nil {
			if len(activeBus.PendingWork) > 0 {
				for _, workID := range activeBus.PendingWork {
					var work, exists = internal.Work[workID]
					if exists {
						work()
					}
				}
			}

			activeBus.Draw(dirtyDraw)
		} // draw all batches stored in the bus for this frame

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
