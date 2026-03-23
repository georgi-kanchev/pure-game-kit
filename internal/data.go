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
type TileData struct {
	Image   *rl.Image
	Texture *rl.Texture2D
}
type TileSet struct {
	TextureId             string
	TileWidth, TileHeight int
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
var Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed
var ShaderTileMapLoc int32

var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)

var TileDatas = make(map[string]*TileData)
var TileSets = make(map[string]*TileSet)

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
var MouseX, MouseY, MouseDeltaX, MouseDeltaY, Scroll, SmoothScroll float32
var Keys, KeysPrev, Buttons, ButtonsPrev = []int{}, []int{}, []int{}, []int{}
var AnyButtonJustPressed, AnyButtonJustReleased, AnyKeyJustPressed, AnyKeyJustReleased = false, false, false, false

var sineTable [3600]float32

//go:embed shaders/quad.frag
var fragQuad string

//go:embed shaders/default.vert
var vertDefault string

//=================================================================

func AssetData(assetId string) (tex *rl.Texture2D, src rl.Rectangle, rotations int, flip bool) {
	var texture, hasTexture = Textures[assetId]
	src = rl.NewRectangle(0, 0, 0, 0)
	if !hasTexture {
		var rect, hasArea = AtlasRects[assetId]
		if hasArea {
			var atlas, _ = Atlases[rect.AtlasId]
			var tex, _ = Textures[atlas.TextureId]

			texture = tex
			src.X = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
			src.Y = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
			src.Width = float32(atlas.CellWidth * int(rect.CountX))
			src.Height = float32(atlas.CellHeight * int(rect.CountY))
			rotations, flip = rect.Rotations, rect.Flip
		}
	} else {
		src.Width, src.Height = float32(texture.Width), float32(texture.Height)
	}
	tex = texture
	return
}
func EditAssetRects(src, dst *rl.Rectangle, ang float32, rotations int, flip bool) {
	if dst.Width < 0 {
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang+180, -dst.Width)
		src.X += src.Width
		src.Width *= -1
	}
	if dst.Height < 0 {
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang+270, -dst.Height)
		src.Y += src.Height
		src.Height *= -1
	}

	if flip {
		src.Width *= -1
	}
	switch rotations % 4 {
	case 1: // 90
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang, dst.Height)
	case 2: // 180
		src.Height *= -1
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang, dst.Width)
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang+90, dst.Height)
	case 3: // 270
		dst.X, dst.Y = moveAtAngle(dst.X, dst.Y, ang+90, dst.Width)
	}
}

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

	var tileData, hasTileData = TileDatas[assetId]
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
	var pos = rl.GetMousePosition()
	MouseDeltaX, MouseDeltaY = delta.X, delta.Y
	MouseX, MouseY = pos.X, pos.Y
	Scroll = rl.GetMouseWheelMoveV().Y

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += Scroll * scrollAccel * DeltaTime
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
	if Shader.ID == 0 {
		Shader = rl.LoadShaderFromMemory(string(vertDefault), string(fragQuad))
		ShaderLoc = rl.GetLocationUniform(Shader.ID, "u")
		ShaderTileMapLoc = rl.GetLocationUniform(Shader.ID, "tileData")
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
