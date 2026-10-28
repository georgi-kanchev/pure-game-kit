package window

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/text"
	f "pure-game-kit/packages/window/filter"
	m "pure-game-kit/packages/window/mode"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
	// rl.InitWindow(3840, 1080, title) // ultra-wide test, 32:9
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
		Create("game", true, true)
	}
	if terminate {
		return false
	}

	internal.Draw()
	rl.DisableDepthTest()
	rl.EndDrawing()

	//=================================================================

	rl.BeginDrawing()
	rl.EnableDepthTest()
	rl.ClearScreenBuffers()

	internal.QuadQueues = 0
	internal.FrameDelta = rl.GetFrameTime()
	internal.UpdateCommands()
	internal.UpdateWindowData()
	internal.UpdateAudio()
	internal.UpdateTimeData()
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

func SetMode(mode m.Enum) {
	var curr = CurrentMode()
	if curr == mode {
		return
	}

	if mode == m.Fullscreen {
		rl.RestoreWindow()

		if curr == m.FullscreenBorderless {
			rl.ToggleBorderlessWindowed()
		}

		var m = rl.GetCurrentMonitor()
		var w, h = rl.GetMonitorWidth(m), rl.GetMonitorHeight(m)

		rl.SetWindowSize(w, h)
		rl.ToggleFullscreen()
		return
	}

	if curr == m.Fullscreen {
		rl.ToggleFullscreen()
	}

	if curr == m.FullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}
	rl.RestoreWindow()

	if mode == m.FullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}

	if mode == m.FullscreenBorderless || mode == m.Maximized {
		rl.MaximizeWindow()
	}

	if mode == m.Floating {
		var m = rl.GetCurrentMonitor()
		var pos = rl.GetMonitorPosition(m)
		var ww, wh = Size()
		rl.SetWindowPosition(int(pos.X+ww/4), int(pos.Y+wh/4))
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
func SetQuality(pixelScale float32, filter f.Enum) {
	internal.PixelScale, internal.Filter = max(pixelScale, 1), uint8(filter)
}

func MoveToMonitor(monitor int) {
	if rl.GetCurrentMonitor() == monitor {
		return
	}

	var wasMax = rl.IsWindowMaximized()
	if wasMax {
		rl.RestoreWindow()
	}

	rl.SetWindowMonitor(monitor)

	if wasMax {
		rl.MaximizeWindow()
	}
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
func CurrentMode() m.Enum {
	var fs = rl.IsWindowFullscreen()
	var bor = rl.IsWindowState(rl.FlagBorderlessWindowedMode)
	var max = rl.IsWindowMaximized()

	if fs && !bor && !max {
		return m.Fullscreen
	}
	if !fs && bor && max {
		return m.FullscreenBorderless
	}
	if !fs && !bor && max {
		return m.Maximized
	}

	return m.Floating
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
func Quality() (pixelScale float32, filter f.Enum) {
	return internal.PixelScale, f.Enum(internal.Filter)
}

// private ========================================================

var gameLogicStart time.Time
var terminate, isInit bool
