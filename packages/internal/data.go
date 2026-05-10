package internal

import (
	"encoding/xml"
	"math"
	"pure-game-kit/packages/utility/number"
	"strings"

	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FontData struct {
	AtlasId int32    // see assets.ImageId
	XMLName xml.Name `xml:"font"`

	Info struct {
		Face     string `xml:"face,attr"`
		Size     int    `xml:"size,attr"`
		Bold     int    `xml:"bold,attr"`
		Italic   int    `xml:"italic,attr"`
		Charset  string `xml:"charset,attr"`
		Unicode  int    `xml:"unicode,attr"`
		StretchH int    `xml:"stretchH,attr"`
		Smooth   int    `xml:"smooth,attr"`
		AA       int    `xml:"aa,attr"`
		Padding  string `xml:"padding,attr"`
		Spacing  string `xml:"spacing,attr"`
		Outline  int    `xml:"outline,attr"`
	} `xml:"info"`

	Common struct {
		LineHeight int `xml:"lineHeight,attr"`
		Base       int `xml:"base,attr"`
		ScaleW     int `xml:"scaleW,attr"`
		ScaleH     int `xml:"scaleH,attr"`
		Pages      int `xml:"pages,attr"`
		Packed     int `xml:"packed,attr"`
		AlphaChnl  int `xml:"alphaChnl,attr"`
		RedChnl    int `xml:"redChnl,attr"`
		GreenChnl  int `xml:"greenChnl,attr"`
		BlueChnl   int `xml:"blueChnl,attr"`
	} `xml:"common"`

	Pages []struct {
		ID   int    `xml:"id,attr"`
		File string `xml:"file,attr"`
	} `xml:"pages>page"`

	DistanceField struct {
		FieldType     string `xml:"fieldType,attr"`
		DistanceRange int    `xml:"distanceRange,attr"`
	} `xml:"distanceField"`

	Chars struct {
		Count int `xml:"count,attr"`
		Chars []struct {
			ID       int    `xml:"id,attr"`
			Index    int    `xml:"index,attr"`
			Char     string `xml:"char,attr"`
			Width    int    `xml:"width,attr"`
			Height   int    `xml:"height,attr"`
			XOffset  int    `xml:"xoffset,attr"`
			YOffset  int    `xml:"yoffset,attr"`
			XAdvance int    `xml:"xadvance,attr"`
			Chnl     int    `xml:"chnl,attr"`
			X        int    `xml:"x,attr"`
			Y        int    `xml:"y,attr"`
			Page     int    `xml:"page,attr"`
		} `xml:"char"`
	} `xml:"chars"`

	Kernings struct {
		Count    int `xml:"count,attr"`
		Kernings []struct {
			First  int `xml:"first,attr"`
			Second int `xml:"second,attr"`
			Amount int `xml:"amount,attr"`
		} `xml:"kerning"`
	} `xml:"kernings"`
}
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

var White1x1 rl.Texture2D
var Textures = make(map[string]rl.Texture2D)
var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)
var Boxes = make(map[string][9]string)

var Fonts = make(map[string]rl.Font)

var Images = make(map[int32]ImageData) // negative = crops; 0 = White1x1; positive = full images
var Fonts2 = make(map[byte]FontData)   // 0 = default
var Font2NextId byte
var NextImageId int16
var NextImageCropId int16

var MatrixDefault rl.Matrix
var Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed
var ShaderTileDataLoc int32

var Sounds = make(map[string]rl.Sound)
var Music = make(map[string]rl.Music)

var TileLayers = make(map[string]*TileLayer)
var TileSets = make(map[string]*TileSet)

var Screens []interface {
	OnLoad()
	OnEnter()
	OnUpdate()
	OnExit()
}
var CurrentScreen int

//=================================================================

var WindowWidth, WindowHeight int
var WindowHovered, WindowFocused, WindowJustResized bool

//=================================================================

var sineTable [3600]float32

//go:embed shaders/quad.frag
var fragQuad string

//go:embed shaders/default.vert
var vertDefault string

//=================================================================

func AssetSize(assetId string) (width, height int) {
	var texture, hasTexture = Textures[assetId]
	width, height = 0, 0

	var tileSet, hasTileSet = TileSets[assetId]
	if hasTileSet && hasTexture {
		return int(texture.Width) / tileSet.TileWidth, int(texture.Height) / tileSet.TileHeight
	}

	if hasTexture {
		return int(texture.Width), int(texture.Height)
	}

	var rect, hasArea = AtlasRects[assetId]
	if hasArea {
		var atlas = Atlases[rect.AtlasId]
		return atlas.CellWidth * int(rect.CountX), atlas.CellHeight * int(rect.CountY)
	}

	var box, hasBox = Boxes[assetId]
	if hasBox {
		var w, h = 0, 0
		for _, id := range box {
			if id == "" {
				continue
			}
			var curW, curH = AssetSize(id)
			w = number.Maximum(curW, h)
			h = number.Maximum(curH, h)
		}
		return w, h
	}

	var font, hasFont = Fonts[assetId]
	if hasFont {
		return int(font.Texture.Width), int(font.Texture.Height)
	}

	var tileData, hasTileData = TileLayers[assetId]
	if hasTileData {
		return int(tileData.Image.Width), int(tileData.Image.Height)
	}

	var sound, hasSound = Sounds[assetId]
	if hasSound {
		return audioDuration(sound.FrameCount, &sound.Stream)
	}
	var music, hasMusic = Music[assetId]
	if hasMusic {
		return audioDuration(music.FrameCount, &music.Stream)
	}

	return
}

func IsLoaded(assetId string) bool {
	var w, h = AssetSize(assetId)
	return w != -1 && h != -1
}

func Path(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// optimized for speed via lookup table
func SinCos(degrees float32) (sin, cos float32) {
	idx := int(degrees * 10)           // convert to index (0.1 degree precision)
	idx = ((idx % 3600) + 3600) % 3600 // and wrap 0-3599

	// sine is direct lookup, cosine is sine shifted by 90 degrees (900 indices)
	return sineTable[idx], sineTable[(idx+900)%3600]
}

func Init() {
	if isInit {
		return
	}
	isInit = true

	if Shader.ID == 0 {
		Shader = rl.LoadShaderFromMemory(string(vertDefault), string(fragQuad))
		ShaderTileDataLoc = rl.GetLocationUniform(Shader.ID, "tileData")
		ShaderLoc = rl.GetLocationUniform(Shader.ID, "u")
	}
	MatrixDefault = rl.MatrixIdentity()

	var img = rl.GenImageColor(1, 1, rl.White)
	White1x1 = rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)

	for i := range 3600 {
		var rad = float64(i) * math.Pi / 1800.0 // convert index to radians (i / 10.0 * Pi / 180.0)
		sineTable[i] = float32(math.Sin(rad))
	}
}
func Update() {
	updateWindowData()
	updateTimeData()
	updateMusic()
	updateScreens()
}

// private ========================================================

var prevCursor int
var isInit bool

func updateWindowData() {
	WindowWidth, WindowHeight = rl.GetScreenWidth(), rl.GetScreenHeight()
	WindowHovered, WindowFocused, WindowJustResized = rl.IsCursorOnScreen(), rl.IsWindowFocused(), rl.IsWindowResized()
}
func updateMusic() {
	for _, v := range Music {
		rl.UpdateMusicStream(v)
	}
}
func updateScreens() {
	if CurrentScreen >= 0 && CurrentScreen < len(Screens) {
		Screens[CurrentScreen].OnUpdate()
	}
}

func audioDuration(frameCount uint32, stream *rl.AudioStream) (seconds, milliseconds int) {
	seconds = int(float32(frameCount) / float32(stream.SampleRate))
	milliseconds = int(math.Mod(float64(seconds), 1.0) * 1000)
	return
}

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
