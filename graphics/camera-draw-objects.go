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
		if n == nil {
			continue
		}

		var x, y, ang, scX, scY = n.TransformToCamera()
		var w, h = n.Width, n.Height
		c.DrawQuad(x, y, w*scX, h*scY, ang, n.Tint)
	}
}
func (c *Camera) DrawSprites(sprites ...*Sprite) {
	c.begin()
	var usingShader = false
	var shouldBatch = len(sprites) > 8
	for _, s := range sprites {
		if s == nil {
			continue
		}

		var x, y, ang, scX, scY = s.TransformToCamera()
		if !c.isAreaVisible(x, y, s.Width*scX, s.Height*scY, ang) {
			continue
		}

		var texture, src, rotations, flip = asset(s.AssetId)
		if texture == nil {
			continue
		}

		var w, h = s.Width, s.Height
		if s.TextureRepeat {
			src.Width, src.Height = w*scX, h*scY
		}
		src.X, src.Y = src.X+s.TextureScrollX, src.Y+s.TextureScrollY

		var dst = rl.Rectangle{X: x, Y: y, Width: w * scX, Height: h * scY}
		editAssetRects(&src, &dst, ang, rotations, flip)

		ang += float32(rotations * 90)

		var effects = condition.If(s.Effects != nil, s.Effects, c.Effects)
		if effects != nil {
			effects.updateUniforms(int(src.Width), int(src.Height))

			if !usingShader && c.Effects == nil {
				rl.BeginShaderMode(internal.Shader)
				rl.EnableDepthTest()
				usingShader = true
			}
		}

		if shouldBatch {
			batch.QueueQuad(texture, src, dst, ang, getColor(s.Tint))
		} else {
			rl.DrawTexturePro(*texture, src, dst, rl.Vector2{}, ang, getColor(s.Tint))
		}
	}
	if shouldBatch {
		batch.Draw()
	}
	if usingShader && c.Effects == nil {
		rl.EndShaderMode()
		rl.DisableDepthTest()
	}
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
	var prevShader = batch.material.Shader
	batch.material.Shader = internal.ShaderText
	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		var x, y, ang, scX, scY = t.TransformToCamera()
		if !c.isAreaVisible(x, y, t.Width*scX, t.Height*scY, ang) {
			continue
		}

		if t.Fast {
			var text = removeTags(txt.Remove(t.Text, "\v"))
			text = condition.If(t.WordWrap, t.TextWrap(text), text)
			var symbol = &symbol{Color: t.Tint, OutlineColor: 255, ShadowColor: 0}
			var pack = packSymbolColor(symbol)
			var col = color.RGBA(pack.R, pack.G, pack.B, pack.A)
			c.DrawTextAdvanced(t.FontId, text, t.X, t.Y, t.LineHeight, t.Thickness, t.gapSymbols(), col)
			continue
		}

		var _, symbols = t.formatSymbols(c)
		var font = t.font()
		var gapX = t.gapSymbols()

		for _, s := range symbols {
			batch.QueueSymbol(font, s, t.LineHeight, gapX)
		}
		// rl.SetShaderValue(internal.ShaderText, internal.ShaderTextLoc, symb, rl.ShaderUniformVec2)
		batch.Draw()

	}
	batch.material.Shader = prevShader
	c.end()
}
