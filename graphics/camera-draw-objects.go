package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (camera *Camera) DrawNodes(nodes ...*Node) {
	for _, n := range nodes {
		if n == nil {
			continue
		}

		var x, y, ang, scX, scY = n.TransformToCamera()
		if !camera.isAreaVisible(x, y, n.Width*scX, n.Height*scY, ang) {
			continue
		}

		var w, h = n.Width, n.Height
		camera.DrawQuad(x, y, w*scX, h*scY, ang, n.Color)
	}
}
func (camera *Camera) DrawSprites(sprites ...*Sprite) {
	camera.begin()
	for _, s := range sprites {
		if s == nil {
			continue
		}

		var x, y, ang, scX, scY = s.TransformToCamera()
		if !camera.isAreaVisible(x, y, s.Width*scX, s.Height*scY, ang) {
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

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, rl.GetColor(s.Color))
	}
	camera.end()
}
func (camera *Camera) DrawBoxes(boxes ...*Box) {
	camera.begin()
	camera.Batch = true
	for _, s := range boxes {
		if s == nil {
			continue
		}

		var w, h = s.Width, s.Height
		var u, r, d, l = s.EdgeBottom, s.EdgeRight, s.EdgeTop, s.EdgeLeft
		var errX, errY float32 = 2, 2 // this adds margin of error to the middle part (it's behind all other parts)
		var c = s.Color
		var asset, has = internal.Boxes[s.AssetId]

		if !has {
			camera.DrawSprites(&s.Sprite)
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

		drawBoxPart(camera, &s.Node, l-errX/2, u-errY/2, w-l-r+errX, h-u-d+errY, asset[4], c) // center

		// edges
		drawBoxPart(camera, &s.Node, l, 0, w-l-r, u, asset[1], c)   // top
		drawBoxPart(camera, &s.Node, 0, u, l, h-u-d, asset[3], c)   // left
		drawBoxPart(camera, &s.Node, w-r, u, r, h-u-d, asset[5], c) // right
		drawBoxPart(camera, &s.Node, l, h-d, w-l-r, d, asset[7], c) // bottom

		// corners
		drawBoxPart(camera, &s.Node, 0, 0, l, u, asset[0], c)     // top left
		drawBoxPart(camera, &s.Node, w-r, 0, r, u, asset[2], c)   // top right
		drawBoxPart(camera, &s.Node, 0, h-d, l, d, asset[6], c)   // bottom left
		drawBoxPart(camera, &s.Node, w-r, h-d, r, d, asset[8], c) // bottom right
	}
	camera.Batch = false
	camera.end()
}
func (camera *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	camera.begin()
	camera.Batch = true
	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		var x, y, ang, scX, scY = t.TransformToCamera()
		if !camera.isAreaVisible(x, y, t.Width*scX, t.Height*scY, ang) {
			continue
		}

		var _, symbols = t.formatSymbols()
		var lastThickness = t.Thickness
		var assetTag = string(t.EmbeddedAssetsTag)

		beginShader(t, t.Thickness)
		for _, s := range symbols {
			var camX, camY = t.PointToCamera(camera, s.X, s.Y)
			var pos = rl.Vector2{X: camX, Y: camY}

			if s.Thickness != lastThickness {
				endShader()
				beginShader(t, s.Thickness)
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
				sprite.Color = uint(rl.ColorToInt(s.Color))

				endShader()
				camera.update()
				camera.DrawSprites(sprite)
				beginShader(t, s.Thickness)
				continue
			}

			if s.Value != assetTag {
				rl.DrawTextPro(*s.Font, s.Value, pos, rl.Vector2{}, s.Angle, s.Height, 0, s.Color)
			}
		}
		endShader()
	}
	camera.Batch = false
	camera.end()
}
