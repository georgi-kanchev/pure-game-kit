package graphics

import (
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/point"
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

		var x, y, ang, scX, scY = n.TransformToCamera()
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
		var x, y, ang, scX, scY = s.TransformToCamera()
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

		// raylib doesn't seem to have negative width/height???
		if rectWorld.Width < 0 {
			var px, py = point.MoveAt(rectWorld.X, rectWorld.Y, ang+180, -rectWorld.Width)
			rectWorld.X = px
			rectWorld.Y = py
			rectTexture.Width *= -1
		}
		if rectWorld.Height < 0 {
			var px, py = point.MoveAt(rectWorld.X, rectWorld.Y, ang+180, -rectWorld.Width)
			rectWorld.X = px
			rectWorld.Y = py
			rectTexture.Height *= -1
		}

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, rl.GetColor(s.Color))
	}
	camera.end()
}
func (camera *Camera) DrawNineSlices(nineSlices ...*NineSlice) {
	for _, s := range nineSlices {
		if s == nil {
			continue
		}

		var w, h = s.Width, s.Height
		var fx, fy = s.SliceFlipX, s.SliceFlipY
		var u, r, d, l = s.SliceSizes[0], s.SliceSizes[1], s.SliceSizes[2], s.SliceSizes[3]
		drawSlice(camera, &s.Node, l, u, w-l-r, h-u-d, false, false, s.AssetId)
		drawSlice(camera, &s.Node, 0, 0, l, u, fx[0], fy[0], s.SliceIds[0])
		drawSlice(camera, &s.Node, l, 0, w-l-r, u, fx[1], fy[1], s.SliceIds[1])
		drawSlice(camera, &s.Node, w-r, 0, r, u, fx[2], fy[2], s.SliceIds[2])
		drawSlice(camera, &s.Node, w-r, u, r, h-u-d, fx[3], fy[3], s.SliceIds[3])
		drawSlice(camera, &s.Node, w-r, h-d, u, r, fx[4], fy[4], s.SliceIds[4])
		drawSlice(camera, &s.Node, l, h-d, w-l-r, u, fx[5], fy[5], s.SliceIds[5])
		drawSlice(camera, &s.Node, 0, h-d, l, d, fx[6], fy[6], s.SliceIds[6])
		drawSlice(camera, &s.Node, 0, u, l, h-u-d, fx[7], fy[7], s.SliceIds[7])
	}
}
func (camera *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	camera.begin()

	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		beginShader(t, t.Thickness)
		var assetTag = string(t.EmbeddedAssetsTag)
		var colorTag = string(t.EmbeddedColorsTag)
		var thickTag = string(t.EmbeddedThicknessTag)
		var wrapped = t.WrapValue(t.Value)
		var lines = strings.Split(wrapped, "\n")
		var _, _, ang, _, _ = t.TransformToCamera()
		var curX, curY float32 = 0, 0
		var font = t.font()
		var textHeight = (t.LineHeight+t.gapLines())*float32(len(lines)) - t.gapLines()
		var curColor = rl.GetColor(t.Color)
		var curThick = t.Thickness
		var alignX, alignY = number.Limit(t.AlignmentX, 0, 1), number.Limit(t.AlignmentY, 0, 1)
		var colorIndex, assetIndex, thickIndex = 0, 0, 0
		// although some chars are invisible, they still need to be iterated cuz of colorIndex and assetIndex

		for l, line := range lines {
			var tagless = strings.ReplaceAll(line, colorTag, "")
			tagless = strings.ReplaceAll(tagless, thickTag, "")
			var lineSize = rl.MeasureTextEx(*font, tagless, t.LineHeight, t.gapSymbols())
			var lineLength = text.Length(line)
			var skipRender = false // it's not 'continue' to avoid skipping the offset calculations

			curX = (t.Width - lineSize.X) * alignX
			curY = float32(l)*(t.LineHeight+t.gapLines()) + (t.Height-textHeight)*alignY

			// hide text outside the box left, top & bottom
			if curX < 0 || curY < 0 || curY+t.LineHeight-1 > t.Height {
				skipRender = true // no need for right cuz text wraps there
			}

			for c := range lineLength {
				var char = string(line[c])
				var charSize = rl.MeasureTextEx(*font, char, t.LineHeight, 0)

				if line[c] == '\r' {
					continue // use as zerospace character or skip altogether
				}

				if curX+charSize.X > t.Width {
					skipRender = true
				}

				var lastChar = string(line[number.LimitInt(c-1, 0, lineLength-1)])

				if char == colorTag {
					if colorIndex < len(t.EmbeddedColors) {
						curColor = rl.GetColor(t.EmbeddedColors[colorIndex])
						colorIndex++
						continue
					}
					curColor = rl.GetColor(t.Color)
					continue
				}

				if char == thickTag {
					if thickIndex < len(t.EmbeddedThicknesses) {
						curThick = t.EmbeddedThicknesses[thickIndex]
						thickIndex++
						endShader()
						beginShader(t, curThick)
						continue

					}
					curThick = t.Thickness
					endShader()
					beginShader(t, curThick)
					continue
				}

				if char == assetTag && (lastChar != assetTag || c == 0) && assetIndex < len(t.EmbeddedAssetIds) {
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
						sprite.Color = uint(rl.ColorToInt(curColor))

						endShader()
						camera.end()
						camera.DrawSprites(&sprite)
						camera.begin()
						beginShader(t, curThick)
					}
					assetIndex++ // skipping render shouldn't affect assetIndexes
				}

				if char == assetTag {
					curX += charSize.X + t.gapSymbols()
					continue
				}

				if !skipRender {
					var camX, camY = t.PointToCamera(camera, curX, curY)
					var pos = rl.Vector2{X: camX, Y: camY}
					rl.DrawTextPro(*font, char, pos, rl.Vector2{}, ang, t.LineHeight, 0, curColor)
					curX += charSize.X + t.gapSymbols()
				}
			}
		}
		endShader()
	}

	camera.end()
}

// #region private

var reusableSprite = NewSprite("", 0, 0)

func drawSlice(camera *Camera, parent *Node, x, y, w, h float32, flipX, flipY bool, id string) {
	reusableSprite.AssetId = id
	reusableSprite.X, reusableSprite.Y = x, y
	reusableSprite.Parent = parent
	reusableSprite.Width, reusableSprite.Height = w, h
	reusableSprite.ScaleX, reusableSprite.ScaleY = 1, 1

	if flipX {
		reusableSprite.ScaleX = -1
		reusableSprite.X += w
	}
	if flipY {
		reusableSprite.ScaleY = -1
		reusableSprite.Y += h
	}

	camera.DrawSprites(&reusableSprite)
}

func beginShader(t *TextBox, thick float32) {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		var smoothness = []float32{t.Smoothness * t.LineHeight / 5}
		rl.BeginShaderMode(sh)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), smoothness, rl.ShaderUniformFloat)
		setShaderThick(thick)
	}
}
func setShaderThick(thick float32) {
	var sh = internal.ShaderText

	if sh.ID != 0 {
		var thickness = []float32{thick}
		thickness[0] = number.Limit(thickness[0], 0, 0.999)
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
