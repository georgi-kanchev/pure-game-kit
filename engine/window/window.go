package window

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type State byte

const (
	StateWindowed State = iota
	StateWindowedBorderless
	StateFullscreen
	StateFullscreenBorderless
	StateMaximized
	StateMinimized
)

var Title = ""
var Color = struct{ R, G, B byte }{0, 0, 0}
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
	var wasMax = rl.IsWindowMaximized()
	if wasMax {
		rl.RestoreWindow()
	}

	rl.SetWindowMonitor(int(monitor))

	if wasMax {
		rl.MaximizeWindow()
	}
}

func IsFocused() bool {
	return rl.IsWindowFocused()
}

func SetState(state State) {
	var currentState = GetState()
	if currentState == state {
		return
	}

	if state == StateMinimized {
		rl.MinimizeWindow()
		return
	}

	if state == StateFullscreen {
		rl.RestoreWindow()

		if currentState == StateFullscreenBorderless || currentState == StateWindowedBorderless {
			rl.ToggleBorderlessWindowed()
		}

		var m = rl.GetCurrentMonitor()
		var w = rl.GetMonitorWidth(m)
		var h = rl.GetMonitorHeight(m)
		rl.SetWindowSize(w, h)
		rl.ToggleFullscreen()
		return
	}

	// restore to windowed
	if currentState == StateFullscreen {
		rl.ToggleFullscreen()
	}
	if currentState == StateFullscreenBorderless || currentState == StateWindowedBorderless {
		rl.ToggleBorderlessWindowed()
	}
	rl.RestoreWindow()

	// after resotre
	if state == StateFullscreenBorderless || state == StateWindowedBorderless {
		rl.ToggleBorderlessWindowed()
	}
	if state == StateFullscreenBorderless || state == StateMaximized {
		rl.MaximizeWindow()
	}

}
func GetState() State {
	var fs = rl.IsWindowFullscreen()
	var bor = rl.IsWindowState(rl.FlagBorderlessWindowedMode)
	var max = rl.IsWindowMaximized()
	var min = rl.IsWindowMinimized()

	if min {
		return StateMinimized
	}
	if fs && !bor && !max {
		return StateFullscreen
	}
	if !fs && bor && max {
		return StateFullscreenBorderless
	}
	if !fs && bor && !max {
		return StateWindowedBorderless
	}
	if !fs && !bor && max {
		return StateMaximized
	}

	return StateWindowed
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
	if Title != currTitle {
		currTitle = Title
		rl.SetWindowTitle(Title)
	}

	if TargetFrameRate != currTargetFPS {
		rl.SetTargetFPS(int32(TargetFrameRate))
	}
}

//endregion
