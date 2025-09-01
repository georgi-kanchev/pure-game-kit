package window

import (
	"path/filepath"
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/symbols"
	"strings"

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

var Title = strings.TrimSuffix(filepath.Base(file.PathOfExecutable()), filepath.Ext(file.PathOfExecutable()))
var Color uint = 0
var IsVSynced = false
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

	rl.EndScissorMode()
	rl.EndDrawing()
	rl.BeginDrawing()
	rl.ClearBackground(rl.GetColor(Color))

	internal.Update()

	return !rl.WindowShouldClose()
}
func Close() {
	if !rl.IsWindowReady() {
		return
	}

	IsOpen = false
	terminate = true
	rl.CloseWindow()
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
func Monitors() (info []string, current int) {
	var count = rl.GetMonitorCount()
	info = make([]string, count)
	for i := 0; i < count; i++ {
		var refreshRate = rl.GetMonitorRefreshRate(i)
		var name = rl.GetMonitorName(i)
		var w, h = rl.GetMonitorWidth(i), rl.GetMonitorHeight(i)
		info[i] = symbols.New(name, " [", w, "x", h, ", ", refreshRate, "Hz]")
	}
	return info, rl.GetCurrentMonitor()
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

func IsHovered() bool {
	return rl.IsCursorOnScreen()
}
func IsFocused() bool {
	return rl.IsWindowFocused()
}
func Size() (width, height int) {
	return rl.GetScreenWidth(), rl.GetScreenHeight()
}

func SetIcon(assetId string) {
	var texture, fullTexture = internal.Textures[assetId]
	var texX, texY float32 = 0.0, 0.0

	if !fullTexture {
		var rect, has = internal.AtlasRects[assetId]
		if !has {
			return
		}

		var atlas, has2 = internal.Atlases[rect.AtlasId]
		if !has2 {
			return
		}

		var tex, has3 = internal.Textures[atlas.TextureId]
		if !has3 {
			return
		}

		texture = tex
		texX = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
		texY = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
	}

	var texW, texH = internal.AssetSize(assetId)
	var rect = rl.Rectangle{X: texX, Y: texY, Width: float32(texW), Height: float32(texH)}
	var imgPtr = rl.LoadImageFromTexture(*texture)

	rl.ImageCrop(imgPtr, rect)
	rl.SetWindowIcon(*imgPtr)
}

// region private
var terminate = false
var currTitle = ""
var currTargetFPS byte = 0

func tryCreate() {
	if rl.IsWindowReady() {
		return
	}

	var flags uint32 = rl.FlagWindowResizable
	if IsVSynced {
		flags |= rl.FlagVsyncHint
	}
	if IsAntialiased {
		flags |= rl.FlagMsaa4xHint
	}

	rl.SetConfigFlags(flags)
	rl.SetTraceLogLevel(rl.LogNone)
	rl.InitWindow(1280, 720, Title)
	rl.SetExitKey(rl.KeyNull)
	rl.MaximizeWindow()
	internal.WindowReady = true

	tryUpdateProperties()
	IsOpen = true
	terminate = false
}

func tryUpdateProperties() {
	if Title != currTitle {
		currTitle = Title
		rl.SetWindowTitle(Title)
	}

	if TargetFrameRate == 0 {
		rl.SetTargetFPS(99999999)
	} else if TargetFrameRate != currTargetFPS {
		rl.SetTargetFPS(int32(TargetFrameRate))
	}
}

//endregion
