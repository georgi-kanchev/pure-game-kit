package window

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/text"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Mode uint8

const ModeFloating, ModeMaximized, ModeFullscreen, ModeFullscreenBorderless Mode = 0, 1, 2, 3

func Create(title string, vsync, antialias bool) {
	var flags uint32 = rl.FlagWindowResizable
	if vsync {
		internal.WindowVsync = true
		flags |= rl.FlagVsyncHint
	}
	if antialias {
		internal.WindowAntialias = true
		flags |= rl.FlagMsaa4xHint
	}

	rl.SetConfigFlags(flags)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.InitWindow(1600, 900, title)
	rl.SetExitKey(rl.KeyNull)
	rl.MaximizeWindow()
	MoveToMonitor(0)
	SetTargetFPS(60)
	isInit = true

	internal.Init()
}
func KeepOpen() bool {
	internal.GameBusyMicroSec = time.Since(gameLogicStart).Microseconds()
	var engineFrameStart = time.Now()
	if !isInit {
		debug.Print("[window.KeepOpen]: Window not yet created. Call `window.Create()`.")
		return false
	}
	if terminate {
		return false
	}

	internal.UpdateCommands()

	internal.Draw()
	rl.DisableDepthTest()
	rl.EndDrawing()

	rl.BeginDrawing()
	rl.EnableDepthTest()
	rl.ClearScreenBuffers()

	//=================================================================

	internal.UpdateWindowData()
	internal.UpdateAudio()
	internal.UpdateTimeData()

	internal.FrameDelta = rl.GetFrameTime()
	internal.CacheInput()

	var shouldClose = !rl.WindowShouldClose()
	internal.EngineBusyMicroSec = time.Since(engineFrameStart).Microseconds()
	gameLogicStart = time.Now()
	return shouldClose
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
	internal.WindowTargetFPS = fps
	rl.SetTargetFPS(int32(fps))
}
func TakeScreenshot(pngPath string) {
	rl.TakeScreenshot(pngPath)
}

//=================================================================

func MousePosition() (x, y float32) {
	return internal.MouseX, internal.MouseY
}
func MouseDelta() (deltaX, deltaY float32) {
	return internal.MouseDeltaX, internal.MouseDeltaY
}

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
	return internal.WindowTargetFPS
}

func IsFocused() bool {
	return internal.WindowFocused
}
func IsJustResized() bool {
	return internal.WindowJustResized
}

// private ========================================================

var gameLogicStart time.Time
var terminate, isInit bool
