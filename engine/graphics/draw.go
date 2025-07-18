package graphics

import (
	"math"
	"pure-kit/engine/internal"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/point"
	"pure-kit/engine/utility/text"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (camera *Camera) DrawColor(color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.GetColor(color))
}
func (camera *Camera) DrawScreenFrame(thickness int, color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(thickness), rl.GetColor(color))             // upper
	rl.DrawRectangle(int32(x+w-thickness), int32(y), int32(thickness), int32(h), rl.GetColor(color)) // right
	rl.DrawRectangle(int32(x), int32(y+h-thickness), int32(w), int32(thickness), rl.GetColor(color)) // lower
	rl.DrawRectangle(int32(x), int32(y), int32(thickness), int32(h), rl.GetColor(color))             // left
}
func (camera *Camera) DrawGrid(thickness, spacing float32, color uint) {
	camera.begin()

	var sx, sy, sw, sh = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	var ulx, uly = camera.PointFromScreen(sx, sy)
	var urx, ury = camera.PointFromScreen(sx+sw, sy)
	var lrx, lry = camera.PointFromScreen(sx+sw, sy+sh)
	var llx, lly = camera.PointFromScreen(sx, sy+sh)
	var xs = []float32{ulx, urx, llx, lrx}
	var ys = []float32{uly, ury, lly, lry}
	var minX, maxX = xs[0], xs[0]
	var minY, maxY = ys[0], ys[0]

	for i := 1; i < 4; i++ {
		if xs[i] < minX {
			minX = xs[i]
		}
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] < minY {
			minY = ys[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}

	var left = float32(math.Floor(float64(minX/spacing))) * spacing
	var right = float32(math.Ceil(float64(maxX/spacing))) * spacing
	var top = float32(math.Floor(float64(minY/spacing))) * spacing
	var bottom = float32(math.Ceil(float64(maxY/spacing))) * spacing

	// vertical
	for x := left; x <= right; x += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(x), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(x, top, x, bottom, myThickness, color)
	}

	// horizontal
	for y := top; y <= bottom; y += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(y), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(left, y, right, y, myThickness, color)
	}

	// x
	if top <= 0 && bottom >= 0 {
		camera.DrawLine(left, 0, right, 0, thickness*3, color)
	}

	// y
	if left <= 0 && right >= 0 {
		camera.DrawLine(0, top, 0, bottom, thickness*3, color)
	}

	camera.end()
}

func (camera *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	camera.begin()
	rl.DrawLineEx(rl.Vector2{X: ax, Y: ay}, rl.Vector2{X: bx, Y: by}, thickness, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawLinesPath(thickness float32, color uint, points ...[2]float32) {
	camera.begin()
	for i := 1; i < len(points); i++ {
		rl.DrawLineEx(
			rl.Vector2{X: points[i-1][0], Y: points[i-1][1]},
			rl.Vector2{X: points[i][0], Y: points[i][1]}, thickness, rl.GetColor(color))
	}
	camera.end()
}
func (camera *Camera) DrawCircle(x, y, radius float32, color uint) {
	camera.begin()
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawFrame(x, y, width, height, angle, thickness float32, color uint) {
	if thickness == 0 {
		return
	}

	var c = rl.GetColor(color)
	var o = rl.Vector2{X: 0, Y: 0}

	camera.begin()
	if thickness < 0 {
		thickness *= -1
		var trx, try = point.MoveAt(x, y, angle, width-thickness)
		var brx, bry = point.MoveAt(x, y, angle+90, height-thickness)
		rl.DrawRectanglePro(rl.Rectangle{X: x, Y: y, Width: width, Height: thickness}, o, angle, c)
		rl.DrawRectanglePro(rl.Rectangle{X: trx, Y: try, Width: thickness, Height: height}, o, angle, c)
		rl.DrawRectanglePro(rl.Rectangle{X: brx, Y: bry, Width: width, Height: thickness}, o, angle, c)
		rl.DrawRectanglePro(rl.Rectangle{X: x, Y: y, Width: thickness, Height: height}, o, angle, c)
		return
	}

	var x1, y1 = point.MoveAt(x, y, angle-90, thickness)
	var tlx, tly = point.MoveAt(x1, y1, angle-180, thickness)
	var trx, try = point.MoveAt(x1, y1, angle, width)
	var blx, bly = point.MoveAt(tlx, tly, angle+90, height+thickness)
	rl.DrawRectanglePro(rl.Rectangle{X: tlx, Y: tly, Width: width + thickness*2, Height: thickness}, o, angle, c)
	rl.DrawRectanglePro(rl.Rectangle{X: trx, Y: try, Width: thickness, Height: height + thickness*2}, o, angle, c)
	rl.DrawRectanglePro(rl.Rectangle{X: blx, Y: bly, Width: width + thickness*2, Height: thickness}, o, angle, c)
	rl.DrawRectanglePro(rl.Rectangle{X: tlx, Y: tly, Width: thickness, Height: height + thickness*2}, o, angle, c)

	camera.end()
}

func (camera *Camera) DrawRectangle(x, y, width, height, angle float32, color uint) {
	camera.begin()
	var rect = rl.Rectangle{X: x, Y: y, Width: width, Height: height}
	rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawNodes(nodes ...*Node) {
	camera.begin()
	for _, n := range nodes {
		if n == nil {
			continue
		}

		var x, y, ang, scX, scY = n.ToCamera()
		var w, h = n.Width, n.Height
		var rec = rl.Rectangle{X: x, Y: y, Width: w * scX, Height: h * scY}

		rl.DrawRectanglePro(rec, rl.Vector2{}, ang, n.color())
	}
	camera.end()
}

func (camera *Camera) DrawTexture(textureId string, x, y, width, height, angle float32, color uint) {
	camera.begin()
	var texture, _ = internal.Textures[textureId]
	var texX, texY float32 = 0.0, 0.0
	var w, h = width, height
	var texW, texH = texture.Width, texture.Height
	var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW), Height: float32(texH)}
	var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(w), Height: float32(h)}

	rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, 0, rl.GetColor(color))
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

		if !hasTexture {
			var rect, hasArea = internal.AtlasRects[s.AssetId]
			if hasArea {
				var atlas, _ = internal.Atlases[rect.AtlasId]
				var tex, _ = internal.Textures[atlas.TextureId]

				texture = tex
				texX = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
				texY = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
			} else {
				var font, hasFont = internal.Fonts[s.AssetId]
				if !hasFont {
					continue
				}
				texture = &font.Texture
			}

		}

		var w, h = s.Width, s.Height
		var texW, texH = texture.Width, texture.Height
		var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW) * repX, Height: float32(texH) * repY}
		var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(w) * scX, Height: float32(h) * scY}

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, s.color())
	}
	camera.end()
}

func (camera *Camera) DrawText(fontId, value string, x, y float32, color uint) {
	camera.begin()

	var sh = internal.ShaderText
	var pos = rl.Vector2{X: x, Y: y}
	var smoothness = []float32{0}
	var thickness = []float32{0.5}
	var font, has = internal.Fonts[fontId]

	if !has {
		var defaultFont = rl.GetFontDefault()
		font = &defaultFont
	}

	if sh.ID != 0 {
		rl.BeginShaderMode(sh)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), smoothness, rl.ShaderUniformFloat)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "thickness"), thickness, rl.ShaderUniformFloat)
	}

	rl.DrawTextPro(*font, value, pos, rl.Vector2{}, 0, float32(font.BaseSize), 0, rl.GetColor(color))

	if sh.ID != 0 {
		rl.EndShaderMode()
	}

	camera.end()
}
func (camera *Camera) DrawTextBoxes(textBoxes ...*TextBox) {
	camera.begin()

	var sh = internal.ShaderText
	for _, t := range textBoxes {
		if t == nil {
			continue
		}

		var x, y, a, _, _ = t.ToCamera()
		var font = t.font()
		var pos = rl.Vector2{X: x, Y: y}
		var smoothness = []float32{t.Smoothness}
		var thickness = []float32{t.Thickness}
		var valueLength = text.Length(t.Value)
		thickness[0] = number.Limit(thickness[0], 0, 0.999)
		smoothness[0] *= t.LineHeight / 5

		if sh.ID != 0 {
			rl.BeginShaderMode(sh)
			rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "smoothness"), smoothness, rl.ShaderUniformFloat)
			rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "thickness"), thickness, rl.ShaderUniformFloat)
		}

		for i := 0; i < valueLength; i++ {

		}

		rl.DrawTextPro(*font, t.Value, pos, rl.Vector2{}, a, t.LineHeight, t.gapSymbols(), t.color())

		if sh.ID != 0 {
			rl.EndShaderMode()
		}
	}

	camera.end()
}

func (camera *Camera) DrawNineSlices(nineSlices ...*NineSlice) {

}

// #region private
func GetOrDefault(values []float32, index int, defaultValue float32) float32 {
	if index >= len(values) {
		return defaultValue
	}
	return values[index]
}

// #endregion
