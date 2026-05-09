package graphics

import (
	"pure-game-kit/packages/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// view =================================================================

var rlCam = rl.Camera2D{}
var drawText = NewTextBox("", 0, 0)
var drawTexture = NewSprite(0, 0, 0)

var skipStartAndEnd bool

const placeholderCharAsset = '@'

func (v *View) area() (x, y, w, h float32) {
	if v.Area == (Area{}) {
		var ww, wh = window.Size()
		return 0, 0, float32(ww), float32(wh)
	}
	return v.Area.X, v.Area.Y, v.Area.Width, v.Area.Height
}

// call before draw to update view but use screen space instead of view space
func (v *View) update() {
	tryRecreateWindow()

	var sx, sy, sw, sh = v.area()
	rlCam.Target.X = float32(v.X)
	rlCam.Target.Y = float32(v.Y)
	rlCam.Rotation = float32(v.Angle)
	rlCam.Zoom = float32(v.Zoom)
	rlCam.Offset.X = sx + sw/2
	rlCam.Offset.Y = sy + sh/2
}

// call before draw to update view and use view space
func (v *View) begin() {
	v.update()
	if skipStartAndEnd {
		return
	}

	rl.BeginMode2D(rlCam)

	if v.Area != (Area{}) {
		var mx, my, mw, mh = v.area()
		rl.BeginScissorMode(int32(mx), int32(my), int32(mw), int32(mh))
	}

	rl.EnableDepthTest()
	v.Effects.updateUniforms(1, 1, nil, nil, true)

	if v.Blend != 0 {
		rl.BeginBlendMode(rl.BlendMode(v.Blend))
	}
}

// call after draw to get back to using screen space
func (v *View) end() {
	if skipStartAndEnd {
		return
	}

	if v.Blend != 0 {
		rl.EndBlendMode()
	}

	rl.DisableDepthTest()
	if v.Area != (Area{}) {
		rl.EndScissorMode()
	}
	rl.EndMode2D()
}
