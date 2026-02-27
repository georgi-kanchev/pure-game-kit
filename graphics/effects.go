package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/random"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Effects struct {
	BlurX, BlurY, Grayscale, Inversion,
	Gamma, Saturation, Contrast, Brightness,
	DepthZ, PixelSize, OutlineSize float32

	OutlineColor, SilhouetteColor uint

	TileColumns, TileRows, TileWidth, TileHeight byte

	hash uint32
}

func NewEffects() *Effects {
	return &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5}
}

//=================================================================
// private

var Tex rl.Texture2D

var u = make([]float32, 26) // this is cached and passed to the shader packed to spare some cgo calls

func (e *Effects) updateUniforms(texW, texH int) {
	tileDataLoc := rl.GetShaderLocation(internal.Shader, "tileData")
	rl.SetShaderValueTexture(internal.Shader, tileDataLoc, Tex)

	var hash = random.Hash(e)
	if e.hash == hash {
		return // no change in values, no need to update shader
	}
	e.hash = hash

	var or, og, ob, oa = color.Channels(e.OutlineColor)
	var sr, sg, sb, sa = color.Channels(e.SilhouetteColor)
	u[0], u[1] = float32(texW), float32(texH)
	u[2], u[3] = e.BlurX, e.BlurY
	u[4], u[5], u[6], u[7], u[8], u[9] = e.Gamma, e.Saturation, e.Contrast, e.Brightness, e.Grayscale, e.Inversion
	u[10], u[11], u[12] = e.PixelSize, e.DepthZ, e.OutlineSize
	u[13], u[14], u[15], u[16] = float32(or)/255, float32(og)/255, float32(ob)/255, float32(oa)/255
	u[17], u[18], u[19], u[20] = float32(sr)/255, float32(sg)/255, float32(sb)/255, float32(sa)/255
	u[21], u[22] = float32(e.TileColumns), float32(e.TileRows)
	u[23], u[24] = float32(e.TileWidth), float32(e.TileHeight)
	u[25] = 256
	rl.SetShaderValueV(internal.Shader, internal.ShaderLoc, u, rl.ShaderUniformFloat, 26)
}

func IDToColor(id uint32) rl.Color {
	return rl.Color{
		R: uint8(id & 0xFF),         // Lower 8 bits
		G: uint8((id >> 8) & 0xFF),  // Next 8 bits
		B: uint8((id >> 16) & 0xFF), // Next 8 bits
		A: uint8((id >> 24) & 0xFF), // Upper 8 bits
	}
}
