package internal

import (
	"math"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"strings"

	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

var WindowReady = false

//=================================================================

var White *rl.Texture2D
var Textures = make(map[string]*rl.Texture2D)
var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)
var Boxes = make(map[string][9]string)

var Fonts = make(map[string]*rl.Font)

var MatrixDefault rl.Matrix
var ShaderText, Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed

//go:embed shaders/text.frag
var fragText string

//go:embed shaders/sprite.frag
var fragSprite string

//go:embed shaders/default.vert
var vertDefault string

var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)

var TiledTilesets = make(map[string]*Tileset)
var TiledMaps = make(map[string]*Map)
var TiledProjects = make(map[string]*Project)
var TiledWorlds = make(map[string][2]float32) // used to store map offsets in the world when reloading maps

var Screens []interface {
	OnLoad()
	OnEnter()
	OnUpdate()
	OnExit()
}
var CurrentScreen int

//=================================================================

var Cursor int
var Input = ""
var MouseDeltaX, MouseDeltaY, SmoothScroll float32
var Keys, KeysPrev, Buttons, ButtonsPrev = []int{}, []int{}, []int{}, []int{}
var AnyButtonJustPressed, AnyButtonJustReleased, AnyKeyJustPressed, AnyKeyJustReleased = false, false, false, false

var sineTable [3600]float32

//=================================================================

func AssetSize(assetId string) (width, height int) {
	var texture, hasTexture = Textures[assetId]
	width, height = -1, -1

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
			w = number.Biggest(curW, h)
			h = number.Biggest(curH, h)
		}
		return w, h
	}

	var font, hasFont = Fonts[assetId]
	if hasFont {
		return int(font.Texture.Width), int(font.Texture.Height)
	}

	var tileset, hasTileset = TiledTilesets[assetId]
	if hasTileset {
		return int(tileset.Columns), int(tileset.TileCount / tileset.Columns)
	}

	var tiledMap, hasMap = TiledMaps[assetId]
	if hasMap {
		return int(tiledMap.Width), int(tiledMap.Height)
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

//=================================================================
// private

var prevCursor int

func updateTimers() {
	for k, v := range CallAfter {
		if Runtime > k {
			for _, f := range v {
				f()
				delete(CallAfter, k)
			}
		}
	}
	for k, v := range CallFor {
		for _, f := range v {
			f(number.Biggest(k-Runtime, 0))
		}
		if Runtime > k {
			delete(CallFor, k)
		}
	}
}
func updateInput() {
	AnyButtonJustPressed = false
	AnyButtonJustReleased = false
	ButtonsPrev = collection.Clone(Buttons)

	for i := range 7 {
		if rl.IsMouseButtonPressed(rl.MouseButton(i)) {
			Buttons = append(Buttons, i)
			AnyButtonJustPressed = true
		}
		if rl.IsMouseButtonReleased(rl.MouseButton(i)) {
			Buttons = collection.Remove(Buttons, i)
			AnyButtonJustReleased = true
		}
	}

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
	}
	prevCursor = Cursor

	var delta = rl.GetMouseDelta()
	MouseDeltaX, MouseDeltaY = delta.X, delta.Y

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += rl.GetMouseWheelMoveV().Y * scrollAccel * DeltaTime
	SmoothScroll *= number.Exponential(-scrollDecay * DeltaTime)
	if SmoothScroll != 0 && number.IsWithin(SmoothScroll, 0, 0.0001) {
		SmoothScroll = 0
	}

	if AnyButtonJustPressed || AnyButtonJustReleased || AnyKeyJustPressed || AnyKeyJustReleased {
		SmoothScroll = 0
	}

	//=================================================================

	AnyKeyJustPressed = false
	AnyKeyJustReleased = false
	KeysPrev = collection.Clone(Keys)
	Input = ""

	var char = rl.GetCharPressed()
	for char > 0 {
		Input += string(char)
		char = rl.GetCharPressed()
	}

	checkKeyRange(32, 96)
	checkKeyRange(256, 349)

	if !rl.IsWindowFocused() {
		Keys = []int{}
		Buttons = []int{}
	}
}
func updateMusic() {
	for _, v := range Music {
		rl.UpdateMusicStream(*v)
	}
}
func updateAnimatedTiles() {
	for _, tileset := range TiledTilesets {
		for _, tile := range tileset.AnimatedTiles {
			tile.Update()
		}
	}
}
func updateScreens() {
	if CurrentScreen >= 0 && CurrentScreen < len(Screens) {
		Screens[CurrentScreen].OnUpdate()
	}
}

func checkKeyRange(from, to int) {
	for i := from; i < to+1; i++ {
		if rl.IsKeyPressed(int32(i)) {
			Keys = append(Keys, i)
			AnyKeyJustPressed = true
		}
		if rl.IsKeyReleased(int32(i)) {
			Keys = collection.Remove(Keys, i)
			AnyKeyJustReleased = true
		}
	}
}

func audioDuration(frameCount uint32, stream *rl.AudioStream) (seconds, milliseconds int) {
	seconds = int(float32(frameCount) / float32(stream.SampleRate))
	milliseconds = int(math.Mod(float64(seconds), 1.0) * 1000)
	return
}

func initData() {
	if ShaderText.ID == 0 {
		ShaderText = rl.LoadShaderFromMemory("", fragText)
		// ShaderTextLoc = rl.GetLocationUniform(ShaderText.ID, "thickSmooth")
	}
	if Shader.ID == 0 {
		Shader = rl.LoadShaderFromMemory(string(vertDefault), string(fragSprite))
		ShaderLoc = rl.GetLocationUniform(Shader.ID, "u")
	}
	MatrixDefault = rl.MatrixIdentity()

	var img = rl.GenImageColor(1, 1, rl.White)
	var tex = rl.LoadTextureFromImage(img)
	White = &tex
	rl.UnloadImage(img)

	for i := range 3600 {
		var rad = float64(i) * math.Pi / 1800.0 // convert index to radians (i / 10.0 * Pi / 180.0)
		sineTable[i] = float32(math.Sin(rad))
	}
}
