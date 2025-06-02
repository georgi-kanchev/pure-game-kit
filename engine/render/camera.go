package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	RenderX, RenderY          int
	RenderWidth, RenderHeight int
	X, Y, Angle, Zoom         float32
	Color                     uint
}

func (camera *Camera) DrawRectangle(x, y, width, height float32, color uint) {
	// rl.BeginMode2D(cam)
	rl.DrawRectangle(int32(x), int32(y), int32(width), int32(height), rl.GetColor(color))
	rl.EndMode2D()
}
func Circle(x, y, radius float32, color uint) {
	rl.DrawCircle(int32(x), int32(y), radius, rl.GetColor(color))
}
func TileMap() {

}

// region private

var cam = rl.Camera2D{}

// endregion
