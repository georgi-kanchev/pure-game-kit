package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Effects struct {
	Gamma, Saturation, Contrast, Brightness, Grayscale, Inversion float32 // Ranged -1..1

	BlurX, BlurY, PixelSize, OutlineSize float32

	DepthZ float32 // Requires semi-transparent pixels to be drawn last. Fully opaque pixels work in any sorting.

	OutlineColor, SilhouetteColor uint
}

func NewEffects() *Effects {
	return &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5}
}

//=================================================================
// private

var u = make([]float32, 32) // this is cached and passed to the shader packed to spare some cgo calls

func (e *Effects) updateUniforms(texW, texH int, tileMap *TileMap, textBox *TextBox, force bool) {
	clear(u)
	u[0], u[1] = float32(texW), float32(texH)
	u[4], u[5], u[6], u[7] = 0.5, 0.5, 0.5, 0.5
	u[25] = internal.Runtime

	var dirty = false

	if e != nil {
		var or, og, ob, oa = color.Channels(e.OutlineColor)
		var sr, sg, sb, sa = color.Channels(e.SilhouetteColor)
		u[2], u[3] = e.BlurX, e.BlurY
		u[4], u[5], u[6], u[7], u[8], u[9] = e.Gamma, e.Saturation, e.Contrast, e.Brightness, e.Grayscale, e.Inversion
		u[10], u[11], u[12] = e.PixelSize, e.DepthZ, e.OutlineSize
		u[13], u[14], u[15], u[16] = float32(or)/255, float32(og)/255, float32(ob)/255, float32(oa)/255
		u[17], u[18], u[19], u[20] = float32(sr)/255, float32(sg)/255, float32(sb)/255, float32(sa)/255

		if u[4] != 0.5 || u[5] != 0.5 || u[6] != 0.5 || u[7] != 0.5 || u[8] != 0 || u[9] != 0 {
			u[26] = 1.0 // do calculations for color adjust
		}
		dirty = true
	}

	if tileMap != nil {
		var data = internal.TileDatas[tileMap.TileDataId]
		var atlas = internal.TileSets[tileMap.TileSetId]
		if data != nil && atlas != nil && data.Texture != nil {
			u[21], u[22] = float32(data.Image.Width), float32(data.Image.Height)
			u[23], u[24] = float32(atlas.TileWidth), float32(atlas.TileHeight)

			rl.DrawRenderBatchActive()        // flush raylib's internal batch to mess texture slots
			rl.ActiveTextureSlot(1)           // switch to slot 1
			rl.EnableTexture(data.Texture.ID) // bind data texture there
			rl.SetShaderValueTexture(internal.Shader, internal.ShaderTileMapLoc, *data.Texture)
		}
		dirty = true
	}

	if textBox != nil {
		u[27] = 1.0 // do calculations for sdf text
		u[28], u[29] = textBox.ShadowOffsetX/200, textBox.ShadowOffsetY/200
		dirty = true
	}

	if dirty || force {
		rl.SetShaderValueV(internal.Shader, internal.ShaderLoc, u, rl.ShaderUniformFloat, 32)
	}
}
