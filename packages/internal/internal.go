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

	var theme GUITheme
	theme.Image = GUIImage{Col: s("#ffffff"), BorSz: f(-5), BorCol: s("#000000")}
	theme.Label = GUIText{LineH: f(50), Margin: s("10 10"), Align: s("0.5 0.5"), Gap: s("0 0"), Col: s("#ffffff"),
		OutCol: s("#000000"), ShCol: s("#000000"), ShBlur: f(0.15), ShOff: s("1 1")}
	theme.Text = GUIText{LineH: f(50), Margin: s("20 20"), Align: s("0 0"), Gap: s("0 0"), Col: s("#ffffff"),
		OutCol: s("#000000"), ShCol: s("#000000"), ShBlur: f(0.15), ShOff: s("1 1")}
	theme.Button.Body.GUIImage = GUIImage{Rnds: f(0.5), Col: s("#949494"), BorSz: f(-8), BorCol: s("#808080")}
	theme.Button.Body.Focused = GUIImage{Col: s("#a8a8a8"), BorCol: s("#949494")}
	theme.Button.Body.Clicked = GUIImage{Col: s("#808080"), BorCol: s("#6c6c6c")}
	theme.Button.Body.Disabled = GUIImage{Col: s("#464646"), BorCol: s("#323232")}
	theme.Button.Value.GUIText = GUIText{LineH: f(50), Margin: s("10 10"), Align: s("0.5 0.5"), Gap: s("0 0"),
		Col: s("#ffffff"), OutCol: s("#000000"), ShCol: s("#000000"), ShBlur: f(0.15), ShOff: s("1 1")}
	theme.Scroll.Body.Size, theme.Scroll.Body.GUIImage = f(10), GUIImage{Col: s("#00000080")}
	theme.Scroll.Handle.Speed, theme.Scroll.Handle.GUIImage = f(40), GUIImage{Rnds: f(1), Col: s("#bfbfbf")}
	theme.Scroll.Handle.Focused, theme.Scroll.Handle.Clicked = GUIImage{Col: s("#ffffff")}, GUIImage{Col: s("#7f7f7f")}
	theme.Slider.Body.GUIImage = GUIImage{Rnds: f(1), Col: s("#949494"), BorSz: f(-8), BorCol: s("#808080")}
	theme.Slider.Body.Disabled = GUIImage{Col: s("#323232"), BorCol: s("#464646")}
	theme.Slider.Hnd.GUIImage = GUIImage{Rnds: f(1), Col: s("#ebebeb"), BorSz: f(-8), BorCol: s("#d7d7d7")}
	theme.Slider.Hnd.Focused = GUIImage{Col: s("#ffffff"), BorCol: s("#ebebeb")}
	theme.Slider.Hnd.Clicked = GUIImage{Col: s("#d7d7d7"), BorCol: s("#c3c3c3")}
	theme.Slider.Hnd.Disabled = GUIImage{Col: s("#828282"), BorCol: s("#6e6e6e")}
	theme.Inputbox.Body.GUIImage = GUIImage{Rnds: f(0.3), Col: s("#6c6c6c"), BorSz: f(-8), BorCol: s("#464646")}
	theme.Inputbox.Body.Typing = GUIImage{BorCol: s("#949494")}
	theme.Inputbox.Body.Focused = GUIImage{BorCol: s("#6c6c6c")}
	theme.Inputbox.Body.Disabled = GUIImage{Col: s("#6c6c6c"), BorCol: s("#464646")}
	theme.Inputbox.Value.GUIText = GUIText{LineH: f(50), Margin: s("30 25"), Align: s("0 0.5"), Gap: s("0 0"),
		Col: s("#ffffff"), OutCol: new("#000000"), ShCol: s("#000000"), ShBlur: f(0.15), ShOff: s("1 1")}
	theme.Inputbox.Value.Disabled = GUIText{LineH: f(50), Col: s("#7f7f7f")}
	theme.Inputbox.Placeholder = GUIText{LineH: f(50), Margin: s("30 25"), Align: s("0 0.5"), Gap: s("0 0"),
		Col: s("#464646"), OutCol: s("#000000"), ShCol: s("#00000000"), ShBlur: f(0.15), ShOff: s("1 1")}
	theme.Inputbox.Selection = GUIImage{Rnds: f(0.3), Col: s("#007fff"), BorSz: f(-4), BorCol: s("#28a7ff")}
	theme.Inputbox.Cursor.GUIImage = GUIImage{Rnds: f(1), Col: s("#c3c3c3")}
	theme.Inputbox.Cursor.Width = f(8)
	GUIThemes[0] = theme
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

func s(v string) *string   { return &v }
func f(v float32) *float32 { return &v }
