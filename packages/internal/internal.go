package internal

import (
	"math"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/storage"
	"strings"

	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var GameBusyMicroSec, EngineBusyMicroSec int64

//=================================================================

var WindowWidth, WindowHeight float32
var PixelScale float32 = 1
var Filter uint8
var WindowHovered, WindowFocused, WindowJustResized bool
var WindowVsync, WindowAntialias bool
var WindowTargetFPS byte

//=================================================================

func Init() {
	for i := range 3600 {
		sineTable[i] = float32(math.Sin(float64(i) * math.Pi / 1800.0)) // convert index to radians (i / 10.0 * Pi / 180.0)
	}

	if Shader.ID == 0 {
		Shader = rl.LoadShaderFromMemory(string(shaderVert), string(shaderFrag))
		ShaderTileDataLoc = rl.GetLocationUniform(Shader.ID, "tileData")
		ShaderLoc = rl.GetLocationUniform(Shader.ID, "u")
	}
	DefaultMatrix = rl.MatrixIdentity()
	DefaultMaterial = rl.LoadMaterialDefault()

	var img = rl.LoadImageFromMemory(".png", defaultFontAtlas, int32(len(defaultFontAtlas)))
	var tex = rl.LoadTextureFromImage(img)
	Images[0] = ImageData{Texture: tex, CropX: 0, CropY: 0, CropWidth: float32(img.Width - 1), CropHeight: float32(img.Height - 1)}
	rl.UnloadImage(img)
	rl.SetTextureFilter(tex, rl.FilterTrilinear)

	var font = string(storage.DecompressGZIP(defaultFont))
	var fontData = &FontJSON{}
	storage.FromJSON(font, fontData)
	LoadFont(fontData, 0)

	var theme GuiTheme
	theme.Image = GuiImage{Col: "#ffffff", BorSz: -5, BorCol: "#000000"}
	theme.Label = GuiText{Margin: "10 10", Align: "0.5 0.5", Gap: "0 0", Col: "#ffffff", OutCol: "#000000",
		ShCol: "#000000", ShBlur: 0.15, ShOff: "1 1"}
	theme.Text = GuiText{LineH: 50, Margin: "20 20", Align: "0 0", Gap: "0 0", Col: "#ffffff", OutCol: "#000000",
		ShCol: "#000000", ShBlur: 0.15, ShOff: "1 1"}
	theme.Button.Body.GuiImage = GuiImage{Rnds: 0.5, Col: "#949494", BorSz: -8, BorCol: "#808080"}
	theme.Button.Body.Focused = GuiImage{Col: "#a8a8a8", BorCol: "#949494"}
	theme.Button.Body.Clicked = GuiImage{Col: "#808080", BorCol: "#6c6c6c"}
	theme.Button.Body.Disabled = GuiImage{Col: "#464646", BorCol: "#323232"}
	theme.Button.Value.GuiText = GuiText{Margin: "10 10", Align: "0.5 0.5", Gap: "0 0", Col: "#ffffff", OutCol: "#000000",
		ShCol: "#000000", ShBlur: 0.15, ShOff: "1 1"}
	theme.Scroll.Body.Size, theme.Scroll.Body.GuiImage = 10, GuiImage{Col: "#00000080"}
	theme.Scroll.Handle.Speed, theme.Scroll.Handle.GuiImage = 40, GuiImage{Rnds: 1, Col: "#bfbfbf"}
	theme.Scroll.Handle.Focused, theme.Scroll.Handle.Clicked = GuiImage{Col: "#ffffff"}, GuiImage{Col: "#7f7f7f"}
	theme.Slider.Body.GuiImage = GuiImage{Rnds: 1, Col: "#949494", BorSz: -8, BorCol: "#808080"}
	theme.Slider.Body.Disabled = GuiImage{Col: "#323232", BorCol: "#464646"}
	theme.Slider.Hnd.GuiImage = GuiImage{Rnds: 1, Col: "#ebebeb", BorSz: -8, BorCol: "#d7d7d7"}
	theme.Slider.Hnd.Focused = GuiImage{Col: "#ffffff", BorCol: "#ebebeb"}
	theme.Slider.Hnd.Clicked = GuiImage{Col: "#d7d7d7", BorCol: "#c3c3c3"}
	theme.Slider.Hnd.Disabled = GuiImage{Col: "#828282", BorCol: "#6e6e6e"}
	theme.Inputbox.Body.GuiImage = GuiImage{Rnds: 0.3, Col: "#6c6c6c", BorSz: -8, BorCol: "#464646"}
	theme.Inputbox.Body.Typing, theme.Inputbox.Body.Focused = GuiImage{BorCol: "#949494"}, GuiImage{BorCol: "#6c6c6c"}
	theme.Inputbox.Body.Disabled = GuiImage{Col: "#6c6c6c", BorCol: "#464646"}
	theme.Inputbox.Value.GuiText = GuiText{Margin: "30 25", Align: "0 0.5", Gap: "0 0", Col: "#ffffff", OutCol: "#000000",
		ShCol: "#000000", ShBlur: 0.15, ShOff: "1 1"}
	theme.Inputbox.Value.Disabled = GuiText{Col: "#7f7f7f"}
	theme.Inputbox.Placeholder = GuiText{Margin: "30 25", Align: "0 0.5", Gap: "0 0", Col: "#464646", OutCol: "#000000",
		ShCol: "#00000000", ShBlur: 0.15, ShOff: "1 1"}
	theme.Inputbox.Selection = GuiImage{Rnds: 0.3, Col: "#007fff", BorSz: -4, BorCol: "#28a7ff"}
	theme.Inputbox.Cursor.GuiImage = GuiImage{Rnds: 1, Col: "#c3c3c3"}
	theme.Inputbox.Cursor.Width = 8
	Themes[0] = theme
}
func UpdateWindowData() {
	WindowWidth, WindowHeight = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
	PixelScale = max(PixelScale, 1)
	WindowHovered, WindowFocused, WindowJustResized = rl.IsCursorOnScreen(), rl.IsWindowFocused(), rl.IsWindowResized()
}

func Path(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
func SinCos(degrees float32) (sin, cos float32) {
	var idx = int(degrees * 10)                      // convert to index (0.1 degree precision)
	idx = ((idx % 3600) + 3600) % 3600               // and wrap 0-3599
	return sineTable[idx], sineTable[(idx+900)%3600] // sine is lookup, cosine is sine shifted by 90 degrees (900 indices)
}

// private ========================================================

var sineTable [3600]float32

func moveAtAngle(x, y, angle, step float32) (float32, float32) {
	var sin, cos = SinCos(angle)
	var dirX, dirY = cos, sin
	if dirX == 0 && dirY == 0 {
		return x, y
	}

	var length = number.SquareRoot(dirX*dirX + dirY*dirY)
	x += (dirX / length) * step
	y += (dirY / length) * step
	return x, y
}
