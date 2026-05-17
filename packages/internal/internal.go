package internal

import (
	"math"
	"pure-game-kit/packages/utility/number"
	"strings"

	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ImageData struct {
	Texture rl.Texture2D

	CropX, CropY, CropWidth, CropHeight float32
}

type AtlasRect struct {
	CellX, CellY,
	CountX, CountY float32
	AtlasId   string
	Rotations int
	Flip      bool
}
type Atlas struct {
	TextureId string
	CellWidth, CellHeight,
	Gap int
}
type TileLayer struct {
	Image   *rl.Image
	Texture *rl.Texture2D

	LastDirtyTime   float32
	CellsWithPoints map[int]struct{}
	ObjectPoints    []float32
}
type TileSet struct {
	ImageId               int32
	TileWidth, TileHeight int
	PointsPerTile         map[uint16][]float32
}

//=================================================================

var Textures = make(map[string]rl.Texture2D)
var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)
var Boxes = make(map[string][9]string)

var TileLayers = make(map[string]*TileLayer)
var TileSets = make(map[string]*TileSet)

//=================================================================

var WindowWidth, WindowHeight int
var WindowHovered, WindowFocused, WindowJustResized bool

//=================================================================

var sineTable [3600]float32

//=================================================================

func Init() {
	if Shader.ID == 0 {
		Shader = rl.LoadShaderFromMemory(string(shaderVert), string(shaderFrag))
		ShaderTileDataLoc = rl.GetLocationUniform(Shader.ID, "tileData")
		ShaderLoc = rl.GetLocationUniform(Shader.ID, "u")
	}
	DefaultMatrix = rl.MatrixIdentity()
	DefaultMaterial = rl.LoadMaterialDefault()

	var img = rl.GenImageColor(1, 1, rl.White)
	Images[0] = ImageData{Texture: rl.LoadTextureFromImage(img), CropX: 0.25, CropY: 0.25, CropWidth: 0.5, CropHeight: 0.5}
	rl.UnloadImage(img)

	for i := range 3600 {
		var rad = float64(i) * math.Pi / 1800.0 // convert index to radians (i / 10.0 * Pi / 180.0)
		sineTable[i] = float32(math.Sin(rad))
	}
}
func UpdateWindowData() {
	WindowWidth, WindowHeight = rl.GetScreenWidth(), rl.GetScreenHeight()
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
