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
var WindowHovered, WindowFocused, WindowJustResized bool
var WindowVsync, WindowAntialias bool
var WindowTargetFPS byte

//=================================================================

func Init() {
	for i := range 3600 {
		var rad = float64(i) * math.Pi / 1800.0 // convert index to radians (i / 10.0 * Pi / 180.0)
		sineTable[i] = float32(math.Sin(rad))
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

	var theme = GuiTheme{
		Image: GuiImage{Color: "#ffffff", BorderColor: "#ffffff"},
		Label: GuiText{Margin: "20 20", Align: "0.5 0.5", Gap: "0 0", Color: "#ffffff", OutlineColor: "#ffffff",
			ShadowColor: "#000000", ShadowBlur: 20, ShadowOffset: "30 30"},
		Text: GuiText{LineHeight: 50, Margin: "20 20", Align: "0 0", Gap: "0 0", Color: "#ffffff", OutlineColor: "#ffffff",
			ShadowColor: "#000000", ShadowBlur: 20, ShadowOffset: "30 30"}}
	theme.Button.Body.GuiImage = GuiImage{Roundness: 0.5, Color: "#808080", BorderSize: -8, BorderColor: "#949494"}
	theme.Button.Body.Focused = GuiImage{Color: "#949494", BorderColor: "#a8a8a8"}
	theme.Button.Body.Clicked = GuiImage{Color: "#6c6c6c", BorderColor: "#808080"}
	theme.Button.Body.Disabled = GuiImage{Color: "#323232", BorderColor: "#464646"}
	Themes[0] = theme
}
func UpdateWindowData() {
	WindowWidth, WindowHeight = float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
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
