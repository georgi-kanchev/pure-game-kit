package engine

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/window"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Run(tps uint16, gameLoop func()) {
	if !internal.WindowReady || !rl.IsWindowReady() {
		window.Recreate()
	}

	tps = max(tps, 1)
	// rl.SetTargetFPS(60)
	internal.TargetTPS = tps
	var channel = make(chan internal.DrawData, 1)
	go func() { // updater
		var ticker = time.NewTicker(time.Second / time.Duration(tps))

		for range ticker.C {
			internal.Update()

			var start = time.Now()
			gameLoop()
			internal.TickBusy = float32(time.Since(start).Seconds())

			var drawData = internal.DrawData{}
			select { // pass the draw data to the renderer
			case channel <- drawData:
			default:
				<-channel
				channel <- drawData
			}
		}
	}()

	//=================================================================

	var view = graphics.NewView(1)
	var currDrawData = internal.DrawData{}
	for window.KeepOpen() { // renderer
		select { // pick up the latest draw data if the updater served one
		case latest := <-channel:
			currDrawData = latest
		default:
		}

		_ = currDrawData

		var dt = rl.GetFrameTime()
		internal.FPS = 1.0 / dt
		internal.FrameDelta = dt

		view.DrawTextDebug(true, false, false, false)
	}
}
