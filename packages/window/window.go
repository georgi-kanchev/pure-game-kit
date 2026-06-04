package window

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Mode uint8

const ModeFloating, ModeMaximized, ModeFullscreen, ModeFullscreenBorderless Mode = 0, 1, 2, 3

func Create(title string, vsync, antialias bool) {
	var flags uint32 = rl.FlagWindowResizable
	if vsync {
		flags |= rl.FlagVsyncHint
	}
	if antialias {
		flags |= rl.FlagMsaa4xHint
	}

	rl.SetConfigFlags(flags)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.InitWindow(1600, 900, title)
	rl.SetExitKey(rl.KeyNull)
	rl.MaximizeWindow()
	MoveToMonitor(0)
	isInit = true

	internal.Init()
}
func KeepOpen() bool {
	if !isInit {
		debug.Print("[window.KeepOpen]: Window not yet created. Call `window.Create()`.")
		return false
	}
	if terminate {
		return false
	}

	internal.CloseBatch()

	rl.BeginDrawing()
	rl.EnableDepthTest()
	rl.ClearScreenBuffers()
	internal.Draw()
	rl.DrawFPS(10, 10)
	rl.DisableDepthTest()
	rl.EndDrawing()

	//=================================================================

	internal.UpdateWindowData()
	internal.UpdateAudio()
	internal.UpdateTimeData()

	internal.FrameDelta = rl.GetFrameTime()
	internal.CacheInput()

	internal.ResetBatches()
	return !rl.WindowShouldClose()
}
func Close() {
	rl.CloseWindow()
	terminate = true
}

//=================================================================

func ApplyMode(mode Mode) {
	var curr = CurrentMode()
	if curr == mode {
		return
	}

	if mode == ModeFullscreen {
		rl.RestoreWindow()

		if curr == ModeFullscreenBorderless {
			rl.ToggleBorderlessWindowed()
		}

		var m = rl.GetCurrentMonitor()
		var w, h = rl.GetMonitorWidth(m), rl.GetMonitorHeight(m)

		rl.SetWindowSize(w, h)
		rl.ToggleFullscreen()
		return
	}

	if curr == ModeFullscreen {
		rl.ToggleFullscreen()
	}

	if curr == ModeFullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}
	rl.RestoreWindow()

	if mode == ModeFullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}

	if mode == ModeFullscreenBorderless || mode == ModeMaximized {
		rl.MaximizeWindow()
	}

	if mode == ModeFloating {
		var m = rl.GetCurrentMonitor()
		var pos = rl.GetMonitorPosition(m)
		var ww, wh = Size()
		rl.SetWindowPosition(int(pos.X+ww/4), int(pos.Y+wh/4))
	}
}
func MoveToMonitor(monitor int) {
	var wasMax = rl.IsWindowMaximized()
	if wasMax {
		rl.RestoreWindow()
	}

	rl.SetWindowMonitor(monitor)

	if wasMax {
		rl.MaximizeWindow()
	}
}
func SetIcon(imagePath string) {
	var img = rl.LoadImage(imagePath)
	rl.SetWindowIcon(*img)
}
func SetTargetFPS(fps uint8) {
	targetFPS = fps
	rl.SetTargetFPS(int32(fps))
}

//=================================================================

func Size() (width, height float32) {
	return internal.WindowWidth, internal.WindowHeight
}
func Monitors() (info []string, current int) {
	var count = rl.GetMonitorCount()
	info = make([]string, count)
	for i := range count {
		var refreshRate = rl.GetMonitorRefreshRate(i)
		var name = rl.GetMonitorName(i)
		var w, h = rl.GetMonitorWidth(i), rl.GetMonitorHeight(i)
		info[i] = text.New(name, " (", w, "x", h, ", ", refreshRate, "Hz)")
	}
	return info, rl.GetCurrentMonitor()
}
func CurrentMode() Mode {
	var fs = rl.IsWindowFullscreen()
	var bor = rl.IsWindowState(rl.FlagBorderlessWindowedMode)
	var max = rl.IsWindowMaximized()

	if fs && !bor && !max {
		return ModeFullscreen
	}
	if !fs && bor && max {
		return ModeFullscreenBorderless
	}
	if !fs && !bor && max {
		return ModeMaximized
	}

	return ModeFloating
}
func TargetFPS() uint8 {
	return targetFPS
}

func IsFocused() bool {
	return internal.WindowFocused
}
func IsJustResized() bool {
	return internal.WindowJustResized
}

// private ========================================================

var targetFPS uint8
var terminate, isInit bool
