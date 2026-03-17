package graphics

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	txt "pure-game-kit/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (c *Camera) DrawQuads(quads ...*Quad) {
	c.begin()

	var lastEffects *Effects
	for _, s := range quads {
		if s == nil || !c.IsAreaVisible(s.Area()) {
			continue
		}

		var x, y = s.CornerTopLeft() // applying pivot
		var ang = s.Angle
		var src = rl.NewRectangle(0, 0, 1, 1)
		var dst = rl.NewRectangle(x, y, s.Width*s.ScaleX, s.Height*s.ScaleY)
		var effects = condition.If(s.Effects != nil, s.Effects, c.Effects)

		if lastEffects != nil && effects != nil && *lastEffects != *effects {
			batch.Draw() // effects are different & break the batch
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
		}
		batch.QueueQuad(internal.White, src, dst, ang, getColor(s.Tint))
		lastEffects = effects
	}
	batch.Draw()
	c.end()
}
func (c *Camera) DrawSprites(sprites ...*Sprite) {
	c.begin()

	var lastEffects *Effects
	for _, s := range sprites {
		if s == nil || !c.IsAreaVisible(s.Area()) {
			continue
		}

		var texture, src, rotations, flip = internal.AssetData(s.AssetId)
		if texture == nil {
			texture = internal.White
		}

		if s.TextureRepeat {
			src.Width, src.Height = s.Width*s.ScaleX, s.Height*s.ScaleY
		}
		src.X, src.Y = src.X+s.TextureScrollX, src.Y+s.TextureScrollY
		var x, y = s.CornerTopLeft() // applying pivot
		var ang = s.Angle
		var dst = rl.NewRectangle(x, y, s.Width*s.ScaleX, s.Height*s.ScaleY)
		internal.EditAssetRects(&src, &dst, ang, rotations, flip)

		ang += float32(rotations * 90)

		var effects = condition.If(s.Effects != nil, s.Effects, c.Effects)
		if lastEffects != nil && effects != nil && *lastEffects != *effects {
			batch.Draw() // effects are different & break the batch
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
		}
		batch.QueueQuad(texture, src, dst, ang, getColor(s.Tint))
		lastEffects = effects
	}
	batch.Draw()
	c.end()
}
func (c *Camera) DrawBoxes(boxes ...*Box) {
	c.begin()
	defer c.end()

	skipStartEnd = true
	var lastEffects *Effects
	for _, b := range boxes {
		if b == nil || !c.IsAreaVisible(b.Area()) {
			continue
		}

		var w, h = b.Width, b.Height
		var u, r, d, l = b.EdgeBottom, b.EdgeRight, b.EdgeTop, b.EdgeLeft
		var errX, errY float32 = 2, 2
		var col = getColor(b.Tint)
		var assetIds, has = internal.Boxes[b.AssetId]

		if !has {
			c.DrawSprites(&b.Sprite) // fallback to sprite if no 9-slice exists
			continue
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

		var effects = condition.If(b.Effects != nil, b.Effects, c.Effects)
		if lastEffects != nil && effects != nil && *lastEffects != *effects {
			batch.Draw() // break batch only when moving to a box with different effects
			// delay updateUniforms until we get the first part's src.Width below
		}

		var parts = [9]struct {
			x, y, w, h float32
			id         string
		}{
			{0, 0, l, u, assetIds[0]},                                                 // Top Left
			{l, 0, w - l - r, u, assetIds[1]},                                         // Top
			{w - r, 0, r, u, assetIds[2]},                                             // Top Right
			{0, u, l, h - u - d, assetIds[3]},                                         // Left
			{l - errX/2, u - errY/2, w - l - r + errX, h - u - d + errY, assetIds[4]}, // Center
			{w - r, u, r, h - u - d, assetIds[5]},                                     // Right
			{0, h - d, l, d, assetIds[6]},                                             // Bottom Left
			{l, h - d, w - l - r, d, assetIds[7]},                                     // Bottom
			{w - r, h - d, r, d, assetIds[8]},                                         // Bottom Right
		}

		var ang = b.Angle

		for _, p := range parts {
			var texture, src, rotations, flip = internal.AssetData(p.id)
			if texture == nil {
				texture = internal.White
			}

			var globalX, globalY = b.PointToGlobal(p.x, p.y)
			var dst = rl.NewRectangle(globalX, globalY, p.w*b.ScaleX, p.h*b.ScaleY)
			var partAng = ang

			internal.EditAssetRects(&src, &dst, partAng, rotations, flip)
			partAng += float32(rotations * 90)

			if lastEffects != nil && effects != nil && *lastEffects != *effects {
				batch.Draw() // effects are different & break the batch
				effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
			}
			batch.QueueQuad(texture, src, dst, partAng, col)
			lastEffects = effects
		}
	}

	batch.Draw()
	skipStartEnd = false
}
func (c *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	c.begin()
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
		var gapX = t.gapSymbols() * t.ScaleX
		var effects = condition.If(t.Effects != nil, t.Effects, c.Effects)
		effects.updateUniforms(int(font.Texture.Width), int(font.Texture.Height), nil, t, false)

		for _, s := range symbols {
			batch.QueueSymbol(font, s, t.LineHeight*t.ScaleY, gapX)
		}

		batch.Draw()
	}
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
			texture = internal.White
		}

		var x, y = t.CornerTopLeft() // applying pivot
		var src = rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height))
		var dst = rl.NewRectangle(x, y, t.Width*t.ScaleX, t.Height*t.ScaleY)
		var effects = condition.If(t.Effects != nil, t.Effects, c.Effects)

		rl.BeginShaderMode(internal.Shader)
		effects.updateUniforms(int(texture.Width), int(texture.Height), t, nil, false)
		rl.DrawTexturePro(*texture, src, dst, rl.Vector2{}, t.Angle, getColor(t.Tint))
		rl.EndShaderMode()
	}
	c.end()
}
