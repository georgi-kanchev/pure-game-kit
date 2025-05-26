package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Title = ""
var Color = struct{ R, G, B byte }{0, 0, 0}
var IsMaximized = false
var IsMinimized = false
var IsFullscreen = false
var IsBorderless = false
var IsVerticallySynchronized = false
var IsAntialiased = false
var TargetFrameRate byte = 60
var IsOpen = false

func Recreate() {
	Close()
	tryCreate()
}
func KeepOpen() bool {
	if terminate {
		return false
	}

	tryCreate()
	tryUpdateProperties()

	rl.EndDrawing()
	rl.BeginDrawing()
	rl.ClearBackground(rl.NewColor(Color.R, Color.G, Color.B, 255))
	rl.DrawFPS(0, 0)
	return !rl.WindowShouldClose()
}
func Close() {
	IsOpen = false
	terminate = true
	rl.CloseWindow()
}

func MoveToMonitor(monitor byte) {
	if IsMaximized {
		rl.RestoreWindow()
	}
	rl.SetWindowMonitor(int(monitor))
	if IsMaximized {
		rl.MaximizeWindow()
	}
}

func IsFocused() bool {
	return rl.IsWindowFocused()
}

// region private
var terminate = false
var currTitle = ""
var currTargetFPS byte = 60

func tryCreate() {
	if rl.IsWindowReady() {
		return
	}

	var flags uint32 = rl.FlagWindowResizable
	if IsVerticallySynchronized {
		flags |= rl.FlagVsyncHint
	}
	if IsAntialiased {
		flags |= rl.FlagMsaa4xHint
	}

	rl.SetConfigFlags(flags)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.InitWindow(1280, 720, "")
	rl.SetExitKey(rl.KeyNull)
	rl.MaximizeWindow()
	rl.SetTargetFPS(60)

	tryUpdateProperties()
	IsOpen = true
	terminate = false
}

func tryUpdateProperties() {
	if !IsBorderless && IsFullscreen != rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}

	if IsBorderless != rl.IsWindowState(rl.FlagBorderlessWindowedMode) {
		rl.ToggleBorderlessWindowed()

		if IsFullscreen && !IsMaximized {
			IsMaximized = true // borderless fullscreen
		}
	}

	if IsMaximized != rl.IsWindowMaximized() {
		if IsMaximized {
			rl.MaximizeWindow()
		} else {
			rl.RestoreWindow()
		}
	}

	if IsMinimized != rl.IsWindowMinimized() {
		if IsMinimized {
			rl.MinimizeWindow()
		} else {
			rl.RestoreWindow()
		}
	}

	if Title != currTitle {
		currTitle = Title
		rl.SetWindowTitle(Title)
	}

	if TargetFrameRate != currTargetFPS {
		rl.SetTargetFPS(int32(TargetFrameRate))
	}
}

//endregion
