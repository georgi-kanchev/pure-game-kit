package window

import (
	"pure-kit/engine/utility/time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type State byte

const (
	Floating State = iota
	FloatingBorderless
	Fullscreen
	FullscreenBorderless
	Maximized
	Minimized
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

	time.Update()

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

func ApplyState(state State) {
	var currentState = CurrentState()
	if currentState == state {
		return
	}

	if state == Minimized {
		rl.MinimizeWindow()
		return
	}

	if state == Fullscreen {
		rl.RestoreWindow()

		if currentState == FullscreenBorderless || currentState == FloatingBorderless {
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
	if currentState == Fullscreen {
		rl.ToggleFullscreen()
	}
	if currentState == FullscreenBorderless || currentState == FloatingBorderless {
		rl.ToggleBorderlessWindowed()
	}
	rl.RestoreWindow()

	// after resotre
	if state == FullscreenBorderless || state == FloatingBorderless {
		rl.ToggleBorderlessWindowed()
	}
	if state == FullscreenBorderless || state == Maximized {
		rl.MaximizeWindow()
	}

}
func CurrentState() State {
	var fs = rl.IsWindowFullscreen()
	var bor = rl.IsWindowState(rl.FlagBorderlessWindowedMode)
	var max = rl.IsWindowMaximized()
	var min = rl.IsWindowMinimized()

	if min {
		return Minimized
	}
	if fs && !bor && !max {
		return Fullscreen
	}
	if !fs && bor && max {
		return FullscreenBorderless
	}
	if !fs && bor && !max {
		return FloatingBorderless
	}
	if !fs && !bor && max {
		return Maximized
	}

	return Floating
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
