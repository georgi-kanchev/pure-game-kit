package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (c *Camera) DrawNodes(nodes ...*Node) {
	for _, n := range nodes {
		if n == nil {
			continue
		}

		var x, y, ang, scX, scY = n.TransformToCamera()
		if !c.isAreaVisible(x, y, n.Width*scX, n.Height*scY, ang) {
			continue
		}

		var w, h = n.Width, n.Height
		c.DrawQuad(x, y, w*scX, h*scY, ang, n.Color)
	}
}
func (c *Camera) DrawSprites(sprites ...*Sprite) {
	c.begin()
	for _, s := range sprites {
		if s == nil {
			continue
		}

		var x, y, ang, scX, scY = s.TransformToCamera()
		if !c.isAreaVisible(x, y, s.Width*scX, s.Height*scY, ang) {
			continue
		}

		var texture, hasTexture = internal.Textures[s.AssetId]
		var texX, texY float32 = 0.0, 0.0
		var texW, texH = 0, 0
		var rotations, flip = 0, false

		if !hasTexture {
			var rect, hasArea = internal.AtlasRects[s.AssetId]
			if hasArea {
				var atlas, _ = internal.Atlases[rect.AtlasId]
				var tex, _ = internal.Textures[atlas.TextureId]

				texture = tex
				texX = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
				texY = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
				texW, texH = atlas.CellWidth*int(rect.CountX), atlas.CellHeight*int(rect.CountY)
				rotations, flip = rect.Rotations, rect.Flip
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

		if texture == nil {
			return
		}

		var w, h = s.Width, s.Height
		if s.TextureRepeat {
			texW, texH = int(w*scX), int(h*scY)
		}
		texX, texY = texX+s.TextureScrollX, texY+s.TextureScrollY

		var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW), Height: float32(texH)}
		var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(w) * scX, Height: float32(h) * scY}

		// raylib doesn't seem to have negative width/height???
		if rectWorld.Width < 0 {
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang+180, -rectWorld.Width)
			rectTexture.Width *= -1
		}
		if rectWorld.Height < 0 {
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang+270, -rectWorld.Height)
			rectTexture.Height *= -1
		}

		if flip {
			rectTexture.Width *= -1
		}
		switch rotations % 4 {
		case 1: // 90
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang, rectWorld.Height)
		case 2: // 180
			rectTexture.Height *= -1
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang, rectWorld.Width)
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang+90, rectWorld.Width)
		case 3: // 270
			rectWorld.X, rectWorld.Y = point.MoveAtAngle(rectWorld.X, rectWorld.Y, ang+90, rectWorld.Width)
		}

		ang += float32(rotations * 90)

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, getColor(s.Color))
	}
	c.end()
}
func (c *Camera) DrawBoxes(boxes ...*Box) {
	c.begin()
	var prevBatch = c.Batch
	c.Batch = true
	defer func() {
		c.Batch = prevBatch
		c.end()
	}()

	for _, s := range boxes {
		if s == nil {
			continue
		}

		var w, h = s.Width, s.Height
		var u, r, d, l = s.EdgeBottom, s.EdgeRight, s.EdgeTop, s.EdgeLeft
		var errX, errY float32 = 2, 2 // this adds margin of error to the middle part (it's behind all other parts)
		var col = s.Color
		var asset, has = internal.Boxes[s.AssetId]

		if !has {
			c.DrawSprites(&s.Sprite)
			return // fallback to sprite rendering if no 9slice asset found
		}

		if w < 0 {
			r *= -1
			l *= -1
			errX *= -1
		}
		if h < 0 {
			u *= -1
			d *= -1
			errY *= -1
		}

		if number.IsBetween(w, -(l + r), l+r, false, false) {
			var total = l + r
			if total != 0 {
				var scale = w / total
				l *= scale
				r *= scale
			}
		}
		if number.IsBetween(h, -(u + d), u+d, false, false) {
			var total = u + d
			if total != 0 {
				var scale = h / total
				u *= scale
				d *= scale
			}
		}

		drawBoxPart(c, &s.Node, l-errX/2, u-errY/2, w-l-r+errX, h-u-d+errY, asset[4], col) // center

		// edges
		drawBoxPart(c, &s.Node, l, 0, w-l-r, u, asset[1], col)   // top
		drawBoxPart(c, &s.Node, 0, u, l, h-u-d, asset[3], col)   // left
		drawBoxPart(c, &s.Node, w-r, u, r, h-u-d, asset[5], col) // right
		drawBoxPart(c, &s.Node, l, h-d, w-l-r, d, asset[7], col) // bottom

		// corners
		drawBoxPart(c, &s.Node, 0, 0, l, u, asset[0], col)     // top left
		drawBoxPart(c, &s.Node, w-r, 0, r, u, asset[2], col)   // top right
		drawBoxPart(c, &s.Node, 0, h-d, l, d, asset[6], col)   // bottom left
		drawBoxPart(c, &s.Node, w-r, h-d, r, d, asset[8], col) // bottom right
	}
}
func (c *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	c.begin()
	var prevBatch = c.Batch
	c.Batch = true
	rl.BeginShaderMode(internal.ShaderText)
	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		var x, y, ang, scX, scY = t.TransformToCamera()
		if !c.isAreaVisible(x, y, t.Width*scX, t.Height*scY, ang) {
			continue
		}

		var _, symbols = t.formatSymbols()
		var lastThickness = t.Thickness
		var assetTag = string(t.EmbeddedAssetsTag)
		var thickSmooth = []float32{number.Limit(t.Thickness, 0, 0.999), t.Smoothness * t.LineHeight / 5}

		rl.SetShaderValue(internal.ShaderText, internal.ShaderUniformLoc, thickSmooth, rl.ShaderUniformVec2)
		for _, s := range symbols {
			var camX, camY = t.PointToCamera(c, s.X, s.Y)
			var pos = rl.Vector2{X: camX, Y: camY}

			if s.Thickness != lastThickness {
				thickSmooth[0] = s.Thickness
				rl.SetShaderValue(internal.ShaderText, internal.ShaderUniformLoc, thickSmooth, rl.ShaderUniformVec2)
				lastThickness = s.Thickness
			}

			if s.Value == assetTag && s.AssetId != "" {
				var w, h = internal.AssetSize(s.AssetId)
				var sprite = NewSprite(s.AssetId, camX, camY)
				var aspect = float32(h / w)

				sprite.Height = t.LineHeight
				sprite.Width = sprite.Height * aspect
				sprite.PivotX, sprite.PivotY = 0, 0
				sprite.Angle = s.Angle
				sprite.Color = s.Color

				rl.EndShaderMode()
				c.update()
				c.DrawSprites(sprite)
				rl.BeginShaderMode(internal.ShaderText)
				rl.SetShaderValue(internal.ShaderText, internal.ShaderUniformLoc, thickSmooth, rl.ShaderUniformVec2)
				continue
			}

			if s.Value != assetTag {
				rl.DrawTextPro(*s.Font, s.Value, pos, rl.Vector2{}, s.Angle, s.Height, 0, getColor(s.Color))
			}
		}
	}
	rl.EndShaderMode()
	c.Batch = prevBatch
	c.end()
}
