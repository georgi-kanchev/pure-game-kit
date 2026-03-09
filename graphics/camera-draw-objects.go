package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (c *Camera) DrawNodes(nodes ...*Node) {
	for _, n := range nodes {
		if n != nil {
			var x, y = n.CornerTopLeft() // apply pivot
			c.DrawQuad(x, y, n.Width*n.ScaleX, n.Height*n.ScaleY, n.Angle, n.Tint)
		}
	}
}
func (c *Camera) DrawSprites(sprites ...*Sprite) {
	c.begin()

	var prevShader = batch.material.Shader
	batch.material.Shader = internal.Shader

	var lastEffects *Effects
	var initUniforms = false

	for _, s := range sprites {
		if s == nil || !c.IsAreaVisible(s.Area()) {
			continue
		}

		var texture, src, rotations, flip = asset(s.AssetId)
		if texture == nil {
			continue
		}

		if s.TextureRepeat {
			src.Width, src.Height = s.Width*s.ScaleX, s.Height*s.ScaleY
		}
		src.X, src.Y = src.X+s.TextureScrollX, src.Y+s.TextureScrollY
		var x, y = s.CornerTopLeft() // applying pivot
		var ang = s.Angle
		var dst = rl.NewRectangle(x, y, s.Width*s.ScaleX, s.Height*s.ScaleY)
		editAssetRects(&src, &dst, ang, rotations, flip)

		ang += float32(rotations * 90)

		var effects = condition.If(s.Effects != nil, s.Effects, c.Effects)

		if !initUniforms {
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil)
			initUniforms = true
		}

		if lastEffects != nil && effects != nil && *lastEffects != *effects {
			rl.EnableDepthTest()
			batch.Draw() // effects are different & break the batch
			rl.DisableDepthTest()
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil)
		}
		batch.QueueQuad(texture, src, dst, ang, getColor(s.Tint))
		lastEffects = effects
	}

	rl.EnableDepthTest()
	batch.Draw()
	rl.DisableDepthTest()

	batch.material.Shader = prevShader

	c.end()
}
func (c *Camera) DrawBoxes(boxes ...*Box) {
	c.begin()
	defer c.end()

	for _, s := range boxes {
		if s == nil {
			continue
		}

		var w, h = s.Width, s.Height
		var u, r, d, l = s.EdgeBottom, s.EdgeRight, s.EdgeTop, s.EdgeLeft
		var errX, errY float32 = 2, 2 // this adds margin of error to the middle part (it's behind all other parts)
		var col = s.Tint
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

		reusableSprite.Effects = s.Effects
		drawBoxPart(&s.Sprite, c, l-errX/2, u-errY/2, w-l-r+errX, h-u-d+errY, asset[4], col) // center

		drawBoxPart(&s.Sprite, c, l, 0, w-l-r, u, asset[1], col)   // top
		drawBoxPart(&s.Sprite, c, 0, u, l, h-u-d, asset[3], col)   // left
		drawBoxPart(&s.Sprite, c, w-r, u, r, h-u-d, asset[5], col) // right
		drawBoxPart(&s.Sprite, c, l, h-d, w-l-r, d, asset[7], col) // bottom

		drawBoxPart(&s.Sprite, c, 0, 0, l, u, asset[0], col)     // top left
		drawBoxPart(&s.Sprite, c, w-r, 0, r, u, asset[2], col)   // top right
		drawBoxPart(&s.Sprite, c, 0, h-d, l, d, asset[6], col)   // bottom left
		drawBoxPart(&s.Sprite, c, w-r, h-d, r, d, asset[8], col) // bottom right
	}
}
func (c *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	c.begin()
	var prevShader = batch.material.Shader
	batch.material.Shader = internal.Shader
	for _, t := range textBoxes {
		if t == nil || !c.IsAreaVisible(t.Area()) {
			continue
		}

		if t.Fast {
			var text = removeTags(txt.Remove(t.Text, "\v"))
			text = condition.If(t.WordWrap, t.TextWrap(text), text)
			defaultTextPack.Color = t.Tint
			var pack = packSymbolColor(defaultTextPack)
			var col = color.RGBA(pack.R, pack.G, pack.B, pack.A)
			c.DrawTextAdvanced(t.FontId, text, t.X, t.Y, t.LineHeight, t.gapSymbols(), t.gapLines(), col)
			continue
		}

		var _, symbols = t.formatSymbols()
		var font = t.font()
		var gapX = t.gapSymbols()

		for _, s := range symbols {
			batch.QueueSymbol(font, s, t.LineHeight, gapX)
		}

		var effects = condition.If(t.Effects != nil, t.Effects, c.Effects)
		effects.updateUniforms(int(font.Texture.Width), int(font.Texture.Height), nil, t)
		batch.Draw()

	}
	batch.material.Shader = prevShader
	c.end()
}
func (c *Camera) DrawTileMaps(tileMaps ...*TileMap) {
	c.begin()
	for _, t := range tileMaps {
		if t == nil || !c.IsAreaVisible(t.Area()) {
			continue
		}

		var atlas = internal.TileSets[t.TileSetId]
		var data = internal.TileDatas[t.TileDataId]
		if atlas == nil && data == nil {
			continue
		}

		var texture = internal.Textures[atlas.TextureId]
		if texture == nil {
			continue
		}

		var x, y = t.CornerTopLeft() // applying pivot
		var src = rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height))
		var dst = rl.NewRectangle(x, y, t.Width*t.ScaleX, t.Height*t.ScaleY)

		rl.BeginShaderMode(internal.Shader)
		rl.EnableDepthTest()
		t.Effects.updateUniforms(int(texture.Width), int(texture.Height), t, nil)
		rl.DrawTexturePro(*texture, src, dst, rl.Vector2{}, t.Angle, getColor(t.Tint))
		rl.EndShaderMode()
		rl.DisableDepthTest()
	}
	c.end()
}
