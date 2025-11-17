package internal

import (
	"math"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"strings"

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

var Textures = make(map[string]*rl.Texture2D)
var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)
var TiledTilesets = make(map[string]*Tileset)
var TiledMaps = make(map[string]*Map)
var TiledProjects = make(map[string]*Project)

var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)
var Boxes = make(map[string][9]string)

var Fonts = make(map[string]*rl.Font)
var ShaderText = rl.Shader{}

//=================================================================

var Cursor int
var MouseDeltaX, MouseDeltaY, SmoothScroll float32
var Input = ""
var Keys = []int{}
var KeysPrev = []int{}
var Buttons = []int{}
var AnyButtonPressedOnce = false
var AnyButtonReleasedOnce = false

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

//=================================================================
// private

var prevCursor int

// timers from engine/execution/flow
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

// keys & buttons + scroll from engine/input/keyboard & mouse
func updateInput() {
	AnyButtonPressedOnce = false
	AnyButtonReleasedOnce = false
	for i := range 7 {
		if rl.IsMouseButtonPressed(rl.MouseButton(i)) {
			Buttons = append(Buttons, i)
			AnyButtonPressedOnce = true
		}
		if rl.IsMouseButtonReleased(rl.MouseButton(i)) {
			Buttons = collection.Remove(Buttons, i)
			AnyButtonReleasedOnce = true
		}
	}

	if prevCursor != Cursor {
		rl.SetMouseCursor(int32(Cursor))
	}
	prevCursor = Cursor

	var delta = rl.GetMouseDelta()
	MouseDeltaX, MouseDeltaY = delta.X, delta.Y

	const scrollAccel, scrollDecay = 600.0, 8.0
	SmoothScroll += rl.GetMouseWheelMoveV().Y * scrollAccel * Delta
	SmoothScroll *= float32(math.Exp(-float64(scrollDecay) * float64(Delta)))

	//=================================================================

	Input = ""
	var char = rl.GetCharPressed()
	for char > 0 {
		Input += string(char)
		char = rl.GetCharPressed()
	}

	KeysPrev = collection.Clone(Keys)
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

func checkKeyRange(from, to int) {
	for i := from; i < to+1; i++ {
		if rl.IsKeyPressed(int32(i)) {
			Keys = append(Keys, i)
		}
		if rl.IsKeyReleased(int32(i)) {
			Keys = collection.Remove(Keys, i)
		}
	}
}

func audioDuration(frameCount uint32, stream *rl.AudioStream) (seconds, milliseconds int) {
	seconds = int(float32(frameCount) / float32(stream.SampleRate))
	milliseconds = int(math.Mod(float64(seconds), 1.0) * 1000)
	return
}
