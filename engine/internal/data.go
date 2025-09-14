package internal

import (
	"pure-kit/engine/utility/collection"
	"pure-kit/engine/utility/number"

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
var Boxes = make(map[string][9]string)

var Fonts = make(map[string]*rl.Font)
var Sounds = make(map[string]*rl.Sound)
var Music = make(map[string]*rl.Music)
var ShaderText = rl.Shader{}

var Flows = make(map[string]*Sequence)
var FlowSignals = []string{}
var States = make(map[string]*StateMachine)

var TiledTilesets = make(map[string]*Tileset)
var TiledMaps = make(map[string]*Map)

var Cursor int
var Input = ""
var Keys = []int{}
var KeysPrev = []int{}
var Buttons = []int{}
var AnyButtonPressedOnce = false
var AnyButtonReleasedOnce = false

var WindowReady = false

var prevCursor int

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

	var box, hasBox = Boxes[assetId]
	if hasBox {
		w, h := 0, 0
		for _, id := range box {
			if id == "" {
				continue
			}
			var curW, curH = AssetSize(id)
			if curW > w {
				w = curW
			}
			if curH > h {
				h = curH
			}
		}
		return w, h
	}

	var font, hasFont = Fonts[assetId]
	if hasFont {
		return int(font.Texture.Width), int(font.Texture.Height)
	}

	return
}

//=================================================================
// private

// state machines from engine/execution/states
func updateStates() {
	for _, v := range States {
		if v.CurrentIndex >= 0 && v.CurrentIndex < len(v.States) {
			v.States[v.CurrentIndex]()
		}
	}
}

// flows from engine/execution/flow
func updateFlows() {
	for _, v := range Flows {
		var prev = v.CurrentIndex // this checks if we changed index inside the step itself, skip increment if so
		var keepGoing = v.CurrentIndex >= 0 && v.CurrentIndex < len(v.Steps) && v.Steps[v.CurrentIndex].Continue()
		if keepGoing && prev == v.CurrentIndex {
			v.CurrentIndex++
		}
	}
}

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

// keys & buttons from engine/input/keyboard & mouse
func updateKeysAndButtons() {
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
