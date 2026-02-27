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
			var x, y = n.PointToGlobal(0, 0) // apply pivot
			c.DrawQuad(x, y, n.Width*n.ScaleX, n.Height*n.ScaleY, n.Angle, n.Tint)
		}
	}
}
func (c *Camera) DrawSprites(sprites ...*Sprite) {
	if Tex.ID == 0 {
		// 1. Generate data on the CPU to ensure pixel-perfect memory
		img := rl.GenImageColor(256, 256, rl.Blank)
		// 2. Write exact bytes directly to RAM
		rl.ImageDrawPixel(img, 0, 0, rl.NewColor(0, 0, 0, 1))
		rl.ImageDrawPixel(img, 255, 3, rl.NewColor(0, 0, 0, 255))

		// 3. Upload raw data directly to the GPU texture
		Tex = rl.LoadTextureFromImage(img)
		rl.SetTextureFilter(Tex, rl.FilterPoint)

		// 4. Free the CPU memory
		rl.UnloadImage(img)
	}

	c.begin()
	var usingShader = false
	var shouldBatch = len(sprites) > 8
	for _, s := range sprites {
		if s == nil {
			continue
		}

		if !c.IsAreaVisible(s.Area()) {
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
		var dst = rl.Rectangle{X: x, Y: y, Width: s.Width * s.ScaleX, Height: s.Height * s.ScaleY}
		editAssetRects(&src, &dst, ang, rotations, flip)

		ang += float32(rotations * 90)

		var effects = condition.If(s.Effects != nil, s.Effects, c.Effects)
		if effects != nil {
			effects.updateUniforms(int(src.Width), int(src.Height))

			if !usingShader && c.Effects == nil {
				rl.BeginShaderMode(internal.Shader)
				effects.updateUniforms(int(src.Width), int(src.Height))
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
	batch.material.Shader = internal.ShaderText
	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		if !c.IsAreaVisible(t.Area()) {
			continue
		}

		if t.Fast {
			var text = removeTags(txt.Remove(t.Text, "\v"))
			text = condition.If(t.WordWrap, t.TextWrap(text), text)
			defaultTextPack.Color = t.Tint
			var pack = packSymbolColor(defaultTextPack)
			var col = color.RGBA(pack.R, pack.G, pack.B, pack.A)
			c.DrawTextAdvanced(t.FontId, text, t.X, t.Y, t.LineHeight, t.gapSymbols(), col)
			continue
		}

		var _, symbols = t.formatSymbols()
		var font = t.font()
		var gapX = t.gapSymbols()

		for _, s := range symbols {
			batch.QueueSymbol(font, s, t.LineHeight, gapX)
		}
		var shadowOffset = []float32{t.ShadowOffsetX / 200, t.ShadowOffsetY / 200}
		rl.SetShaderValue(internal.ShaderText, internal.ShaderTextShOffLoc, shadowOffset, rl.ShaderUniformVec2)
		batch.Draw()

	}
	batch.material.Shader = prevShader
	c.end()
}
