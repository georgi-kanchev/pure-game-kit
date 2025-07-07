package graphics

import (
	"math"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (camera *Camera) DrawColor(color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(h), rl.GetColor(color))
}
func (camera *Camera) DrawFrame(size int, color uint) {
	camera.update()

	var x, y, w, h = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	rl.DrawRectangle(int32(x), int32(y), int32(w), int32(size), rl.GetColor(color))        // upper
	rl.DrawRectangle(int32(x+w-size), int32(y), int32(size), int32(h), rl.GetColor(color)) // right
	rl.DrawRectangle(int32(x), int32(y+h-size), int32(w), int32(size), rl.GetColor(color)) // lower
	rl.DrawRectangle(int32(x), int32(y), int32(size), int32(h), rl.GetColor(color))        // left
}
func (camera *Camera) DrawGrid(thickness, spacing float32, color uint) {
	camera.begin()

	var sx, sy, sw, sh = camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight
	ulx, uly := camera.PointFromScreen(sx, sy)
	urx, ury := camera.PointFromScreen(sx+sw, sy)
	lrx, lry := camera.PointFromScreen(sx+sw, sy+sh)
	llx, lly := camera.PointFromScreen(sx, sy+sh)

	xs := []float32{ulx, urx, llx, lrx}
	ys := []float32{uly, ury, lly, lry}

	minX, maxX := xs[0], xs[0]
	minY, maxY := ys[0], ys[0]

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

	left := float32(math.Floor(float64(minX/spacing))) * spacing
	right := float32(math.Ceil(float64(maxX/spacing))) * spacing
	top := float32(math.Floor(float64(minY/spacing))) * spacing
	bottom := float32(math.Ceil(float64(maxY/spacing))) * spacing

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
func (camera *Camera) DrawRectangle(x, y, width, height, angle float32, color uint) {
	camera.begin()
	var rect = rl.Rectangle{X: x, Y: y, Width: width, Height: height}
	rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawCircle(x, y, radius float32, color uint) {
	camera.begin()
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
	camera.end()
}
func (camera *Camera) DrawNodes(nodes ...*Node) {
	camera.begin()
	for _, node := range nodes {
		if node == nil {
			continue
		}

		var texture, fullTexture = internal.Textures[node.AssetId]
		var texX, texY float32 = 0.0, 0.0
		var repX, repY = node.RepeatX, node.RepeatY
		var x, y, ang, scX, scY = node.ToGlobal()

		if !fullTexture {
			var rect, has = internal.AtlasRects[node.AssetId]
			if !has {
				continue
			}

			var atlas, has2 = internal.Atlases[rect.AtlasId]
			if !has2 {
				continue
			}

			var tex, has3 = internal.Textures[atlas.TextureId]
			if !has3 {
				continue
			}

			texture = tex
			texX = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
			texY = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
		}

		var texW, texH = node.Size()
		var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW) * repX, Height: float32(texH) * repY}
		var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(texW) * scX, Height: float32(texH) * scY}

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, rl.GetColor(node.Tint))
	}
	camera.end()
}
