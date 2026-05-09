package window

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"

	st "pure-game-kit/packages/window/state"
)

var Title = "game"
var IsVSynced = false     // Requires window recreation to take effect.
var IsAntialiased = false // Requires window recreation to take effect.

//=================================================================

func Recreate() {
	Close()
	tryCreate()
}
func KeepOpen() bool {
	if terminate {
		rl.CloseWindow()
		return false
	}

	tryCreate()
	tryUpdateProperties()

	rl.DisableDepthTest()
	rl.EndShaderMode()
	rl.EndBlendMode()
	rl.EndScissorMode()
	rl.EndDrawing()

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	w, h = rl.GetScreenWidth(), rl.GetScreenHeight()

	return !rl.WindowShouldClose()
}
func Close() {
	terminate = true
}

func ApplyState(state int) {
	tryCreate()

	var currentState = CurrentState()
	if currentState == state {
		return
	}

	if state == st.Fullscreen {
		rl.RestoreWindow()

		if currentState == st.FullscreenBorderless {
			rl.ToggleBorderlessWindowed()
		}

		var m = rl.GetCurrentMonitor()
		var w, h = rl.GetMonitorWidth(m), rl.GetMonitorHeight(m)

		rl.SetWindowSize(w, h)
		rl.ToggleFullscreen()
		return
	}

	if currentState == st.Fullscreen {
		rl.ToggleFullscreen()
	}

	if currentState == st.FullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}
	rl.RestoreWindow()

	if state == st.FullscreenBorderless {
		rl.ToggleBorderlessWindowed()
	}

	if state == st.FullscreenBorderless || state == st.Maximized {
		rl.MaximizeWindow()
	}

	if state == st.Floating {
		var m = rl.GetCurrentMonitor()
		var pos = rl.GetMonitorPosition(m)
		var ww, wh = Size()
		rl.SetWindowPosition(int(pos.X)+ww/4, int(pos.Y)+wh/4)
	}
}
func MoveToMonitor(monitor int) {
	tryCreate()

	var wasMax = rl.IsWindowMaximized()
	if wasMax {
		rl.RestoreWindow()
	}

	rl.SetWindowMonitor(monitor)

	if wasMax {
		rl.MaximizeWindow()
	}
}
func SetIcon(assetId string) {
	tryCreate()

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
	var imgPtr = rl.LoadImageFromTexture(texture)

	rl.ImageCrop(imgPtr, rect)
	rl.SetWindowIcon(*imgPtr)
}

//=================================================================

func Size() (width, height int) {
	return w, h
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
func CurrentState() int {
	var fs = rl.IsWindowFullscreen()
	var bor = rl.IsWindowState(rl.FlagBorderlessWindowedMode)
	var max = rl.IsWindowMaximized()

	if fs && !bor && !max {
		return st.Fullscreen
	}
	if !fs && bor && max {
		return st.FullscreenBorderless
	}
	if !fs && !bor && max {
		return st.Maximized
	}

	return st.Floating
}

func IsHovered() bool {
	return rl.IsCursorOnScreen()
}
func IsFocused() bool {
	return rl.IsWindowFocused()
}
func IsJustResized() bool {
	return rl.IsWindowResized()
}

// private ========================================================

var w, h = 0, 0
var terminate = false
var currTitle = ""

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
	terminate = false
	MoveToMonitor(0)

	internal.InitData()
}
func tryUpdateProperties() {
	if Title != currTitle {
		currTitle = Title
		rl.SetWindowTitle(Title)
	}
}
