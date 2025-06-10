package render

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
	camera.start()

	// Get all 4 world-space corners of the screen
	ulx, uly := camera.CornerUpperLeft(0, 0)
	urx, ury := camera.CornerUpperRight(0, 0)
	llx, lly := camera.CornerLowerLeft(0, 0)
	lrx, lry := camera.CornerLowerRight(0, 0)

	// Compute axis-aligned bounding box (AABB) from the rotated corners
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

	// Snap bounds to grid spacing
	left := float32(math.Floor(float64(minX/spacing))) * spacing
	right := float32(math.Ceil(float64(maxX/spacing))) * spacing
	top := float32(math.Floor(float64(minY/spacing))) * spacing
	bottom := float32(math.Ceil(float64(maxY/spacing))) * spacing

	// Draw vertical lines
	for x := left; x <= right; x += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(x), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(x, top, x, bottom, myThickness, color)
	}

	// Draw horizontal lines
	for y := top; y <= bottom; y += spacing {
		var myThickness = thickness
		if float32(math.Mod(float64(y), float64(spacing)*10)) == 0 {
			myThickness *= 2
		}

		camera.DrawLine(left, y, right, y, myThickness, color)
	}

	// Draw X axis
	if top <= 0 && bottom >= 0 {
		camera.DrawLine(left, 0, right, 0, thickness*3, color)
	}

	// Draw Y axis
	if left <= 0 && right >= 0 {
		camera.DrawLine(0, top, 0, bottom, thickness*3, color)
	}

	camera.stop()
}

func (camera *Camera) DrawLine(ax, ay, bx, by, thickness float32, color uint) {
	camera.start()
	rl.DrawLineEx(rl.Vector2{X: ax, Y: ay}, rl.Vector2{X: bx, Y: by}, thickness, rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawLinesPath(thickness float32, color uint, points ...[2]float32) {
	camera.start()
	for i := 1; i < len(points); i++ {
		rl.DrawLineEx(
			rl.Vector2{X: points[i-1][0], Y: points[i-1][1]},
			rl.Vector2{X: points[i][0], Y: points[i][1]}, thickness, rl.GetColor(color))
	}
	camera.stop()
}
func (camera *Camera) DrawRectangle(x, y, width, height, angle float32, color uint) {
	camera.start()
	var rect = rl.Rectangle{X: x, Y: y, Width: width, Height: height}
	rl.DrawRectanglePro(rect, rl.Vector2{X: 0, Y: 0}, angle, rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawCircle(x, y, radius float32, color uint) {
	camera.start()
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
	camera.stop()
}
func (camera *Camera) DrawNodes(nodes ...*Node) {
	camera.start()
	for _, node := range nodes {
		var texture, fullTexture = internal.Textures[node.AssetID]
		var texX, texY float32 = 0.0, 0.0
		var repX, repY = node.RepeatX, node.RepeatY
		var x, y, ang, scX, scY = node.Global()

		if !fullTexture {
			var rect, has = internal.AtlasRects[node.AssetID]
			var atlas = rect.Atlas

			if !has {
				continue
			}

			texture = atlas.Texture
			texX, texY = rect.CellX*float32(atlas.CellWidth), rect.CellY*float32(atlas.CellHeight)
		}

		var texW, texH = node.Size()
		var rectTexture = rl.Rectangle{X: texX, Y: texY, Width: float32(texW) * repX, Height: float32(texH) * repY}
		var rectWorld = rl.Rectangle{X: x, Y: y, Width: float32(texW) * scX, Height: float32(texH) * scY}

		rl.DrawTexturePro(*texture, rectTexture, rectWorld, rl.Vector2{}, ang, rl.White)
	}
	camera.stop()
}
