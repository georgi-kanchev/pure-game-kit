package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Effects struct {
	BlurX, BlurY, Grayscale, Inversion,
	Gamma, Saturation, Contrast, Brightness,
	DepthZ, PixelSize, OutlineSize float32

	OutlineColor, SilhouetteColor uint
}

func NewEffects() *Effects {
	return &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5}
}

//=================================================================
// private

var u = make([]float32, 27) // this is cached and passed to the shader packed to spare some cgo calls

func (e *Effects) updateUniforms(texW, texH int, tileMap *TileMap) {
	u[0], u[1] = float32(texW), float32(texH)
	u[4], u[5], u[6], u[7] = 0.5, 0.5, 0.5, 0.5

	if e != nil {
		var or, og, ob, oa = color.Channels(e.OutlineColor)
		var sr, sg, sb, sa = color.Channels(e.SilhouetteColor)
		u[2], u[3] = e.BlurX, e.BlurY
		u[4], u[5], u[6], u[7], u[8], u[9] = e.Gamma, e.Saturation, e.Contrast, e.Brightness, e.Grayscale, e.Inversion
		u[10], u[11], u[12] = e.PixelSize, e.DepthZ, e.OutlineSize
		u[13], u[14], u[15], u[16] = float32(or)/255, float32(og)/255, float32(ob)/255, float32(oa)/255
		u[17], u[18], u[19], u[20] = float32(sr)/255, float32(sg)/255, float32(sb)/255, float32(sa)/255
	}

	if tileMap != nil {
		var data = internal.TileDatas[tileMap.TileDataId]
		var atlas = internal.TileAtlases[tileMap.TileAtlasId]
		if data != nil && atlas != nil && data.Texture != nil {
			u[21], u[22] = float32(data.Image.Width), float32(data.Image.Height)
			u[23], u[24] = float32(atlas.TileWidth), float32(atlas.TileHeight)

			var loc = rl.GetShaderLocation(internal.Shader, "tileData")
			rl.SetShaderValueTexture(internal.Shader, loc, *data.Texture)
		}
	}

	u[25] = internal.Runtime
	rl.SetShaderValueV(internal.Shader, internal.ShaderLoc, u, rl.ShaderUniformFloat, 27)
}
