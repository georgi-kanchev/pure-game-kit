package internal

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AtlasRect struct {
	CellX, CellY, CountX, CountY float32
	AtlasId                      string
}
type Atlas struct {
	TextureId                  string
	CellWidth, CellHeight, Gap int
}

type Sequence struct {
	Steps        []Step
	CurrentIndex int
}

type StateMachine struct {
	States       []func()
	CurrentIndex int
}

type Step interface {
	Continue() bool
}

var Textures = make(map[string]*rl.Texture2D)
var AtlasRects = make(map[string]AtlasRect)
var Atlases = make(map[string]Atlas)
var TiledData = make(map[string]TiledMap)

var Fonts = make(map[string]*rl.Font)
var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)
var ShaderText = rl.Shader{}

var Flows = make(map[string]*Sequence)
var FlowSignals = []string{}
var States = make(map[string]*StateMachine)

func AssetSize(assetId string) (width, height int) {
	var texture, hasTexture = Textures[assetId]
	width, height = 0, 0

	if hasTexture {
		return int(texture.Width), int(texture.Height)
	}

	var rect, hasArea = AtlasRects[assetId]
	if hasArea {
		var atlas = Atlases[rect.AtlasId]
		return atlas.CellWidth * int(rect.CountX), atlas.CellHeight * int(rect.CountY)
	}

	var font, hasFont = Fonts[assetId]
	if hasFont {
		return int(font.Texture.Width), int(font.Texture.Height)
	}

	return
}
