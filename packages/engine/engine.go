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
	var channel = make(chan internal.DrawData, 1)
	go func() { // updater
		var ticker = time.NewTicker(time.Second / time.Duration(internal.TargetTPS))

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
	for !rl.WindowShouldClose() { // renderer
		select { // pick up the latest draw data if the updater served one
		case latest := <-channel:
			currDrawData = latest
		default:
		}

		_ = currDrawData

		var dt = rl.GetFrameTime()
		internal.FPS = 1.0 / dt
		internal.FrameDelta = dt

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		view.DrawTextDebug(true, false, false, false)

		rl.DisableDepthTest()
		rl.EndShaderMode()
		rl.EndBlendMode()
		rl.EndScissorMode()
		rl.EndDrawing()
	}
}
