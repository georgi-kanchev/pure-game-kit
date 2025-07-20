package graphics

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (camera *Camera) DrawNodes(nodes ...*Node) {
	camera.begin()
	for _, n := range nodes {
		if n == nil {
			continue
		}

		var x, y, ang, scX, scY = n.ToCamera()
		var w, h = n.Width, n.Height
		var rec = rl.Rectangle{X: x, Y: y, Width: w * scX, Height: h * scY}

		rl.DrawRectanglePro(rec, rl.Vector2{}, ang, rl.GetColor(n.Color))
	}
	camera.end()
}
func (camera *Camera) DrawSprites(sprites ...*Sprite) {
	camera.begin()
	for _, s := range sprites {
		if s == nil {
			continue
		}

		var texture, hasTexture = internal.Textures[s.AssetId]
		var texX, texY float32 = 0.0, 0.0
		var repX, repY = s.RepeatX, s.RepeatY
		var x, y, ang, scX, scY = s.ToCamera()
		var texW, texH = 0, 0

		if !hasTexture {
			var rect, hasArea = internal.AtlasRects[s.AssetId]
			if hasArea {
				var atlas, _ = internal.Atlases[rect.AtlasId]
				var tex, _ = internal.Textures[atlas.TextureId]

				texture = tex
				texX = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
				texY = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
				texW, texH = atlas.CellWidth, atlas.CellHeight
			} else {
				var font, hasFont = internal.Fonts[s.AssetId]
				if !hasFont {
					continue
				}
				texture = &font.Texture
				texW, texH = int(texture.Width), int(texture.Height)
			}
		} else {
			texW, texH = int(texture.Width), int(texture.Height)
		}

		var w, h = s.Width, s.Height
		var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW) * repX, Height: float32(texH) * repY}
		var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(w) * scX, Height: float32(h) * scY}

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, rl.GetColor(s.Color))
	}
	camera.end()
}
func (camera *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	camera.begin()

	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		beginShader(t)
		const colorTag = "`"
		const assetTag = "@"
		var wrapped = t.WrapValue(t.Value)
		var lines = strings.Split(wrapped, "\n")
		var _, _, ang, _, _ = t.ToCamera()
		var curX, curY float32 = 0, 0
		var font = t.font()
		var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
		var color = rl.GetColor(t.Color)
		var colorIndex, assetIndex = 0, 0

		for l, line := range lines {
			var tagless = strings.ReplaceAll(line, colorTag, "")
			tagless = strings.ReplaceAll(tagless, assetTag, "")
			var lineSize = rl.MeasureTextEx(*font, tagless, t.LineHeight, t.gapSymbols())
			var lineLength = text.Length(line)
			var skipRender = false

			curX = (t.Width - lineSize.X) * t.AlignmentX
			curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*t.AlignmentY

			// hide text outside the box left, top & bottom
			if curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height {
				skipRender = true // no need for right cuz text wraps there
			}

			for c := range lineLength {
				var char = string(line[c])

				if char == colorTag {
					if colorIndex < len(t.EmbeddedColors) {
						color = rl.GetColor(t.EmbeddedColors[colorIndex])
						colorIndex++
						continue
					}
					color = rl.GetColor(t.Color)
					continue
				}

				if char == assetTag && assetIndex < len(t.EmbeddedAssetIds) {
					if !skipRender {
						var assetId = t.EmbeddedAssetIds[assetIndex]
						var w, h = internal.AssetSize(assetId)
						var camX, camY = t.PointToCamera(camera, curX, curY)
						var sprite = NewSprite(assetId, camX, camY)
						var aspect = float32(h / w)

						sprite.Height = t.LineHeight
						sprite.Width = sprite.Height * aspect
						sprite.PivotX, sprite.PivotY = 0, 0
						sprite.Angle = ang
						sprite.Color = uint(rl.ColorToInt(color))

						if curX+float32(sprite.Width) < t.Width && curY+t.LineHeight-1 < t.Height {
							endShader()
							camera.end()
							camera.DrawSprites(&sprite)
							camera.begin()
							beginShader(t)
						}
					}
					assetIndex++
					continue
				}

				if !skipRender {
					var charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0)
					var camX, camY = t.PointToCamera(camera, curX, curY)
					var pos = rl.Vector2{X: camX, Y: camY}
					rl.DrawTextPro(*font, char, pos, rl.Vector2{}, ang, t.LineHeight, 0, color)
					curX += charSize.X + t.gapSymbols()
				}
			}
		}
		endShader()
	}

	camera.end()
}

func (camera *Camera) DrawNineSlices(nineSlices ...*NineSlice) {

}

// #region private

func beginShader(t *TextBox) {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		var smoothness = []float32{t.Smoothness}
		var thickness = []float32{t.Thickness}
		thickness[0] = number.Limit(thickness[0], 0, 0.999)
		smoothness[0] *= t.LineHeight / 5
		rl.BeginShaderMode(sh)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), smoothness, rl.ShaderUniformFloat)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "thickness"), thickness, rl.ShaderUniformFloat)
	}
}
func endShader() {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		rl.EndShaderMode()
	}
}

// #endregion
