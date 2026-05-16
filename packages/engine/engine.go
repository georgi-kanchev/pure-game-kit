// The entry point of the game and the pump of the game loop.
package engine

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Initialize(title string, maxFps uint16, vsync, antialias bool) {
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

	internal.Init()
}

func Run(gameLoop func()) {
	for !rl.WindowShouldClose() {
		internal.UpdateWindowData()
		internal.UpdateAudio()

		internal.FrameDelta = rl.GetFrameTime()
		internal.AccumulateAndSyncInput()

		internal.ResetBatches()
		internal.UpdateTimeData()
		gameLoop()
		internal.CloseBatch()

		rl.BeginDrawing()
		rl.EnableDepthTest()
		rl.ClearBackground(rl.Black)
		internal.Draw()
		rl.DrawFPS(10, 10)
		rl.DisableDepthTest()
		rl.EndDrawing()
	}
	rl.CloseWindow()
}

func Stop() {
	rl.CloseWindow()
}
