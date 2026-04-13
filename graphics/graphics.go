// Strictly tied to the window, drawing on it through a camera and converting between the two coordinate systems.
// The camera's drawing consists of two categories: primitives and objects. While using the assets for drawing,
// the graphical objects are still very lightweight and exist independently of them.
// The concept map of the package looks like this:
//   - Camera - draws objects
//   - Quad - no asset, flat color
//   - ├ Sprite - texture asset
//   - ├ NinePatch - box asset
//   - ├ TextBox - font asset
//   - └ TileMap - tile set asset + tile data asset
package graphics

import (
	"image/color"
	"pure-game-kit/internal"
	col "pure-game-kit/utility/color"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Area struct{ X, Y, Width, Height float32 }

//=================================================================
// private

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}

func getColor(value uint) color.RGBA {
	var r, g, b, a = col.Channels(value)
	return color.RGBA{R: r, G: g, B: b, A: a}
}
func packSymbolColor(s symbol) rl.Color {
	var packLayer = func(c rl.Color) uint8 {
		var r = (c.R >> 6) & 0x03
		var g = (c.G >> 6) & 0x03
		var b = (c.B >> 6) & 0x03
		var a = (c.A >> 6) & 0x03
		return (r << 6) | (g << 4) | (b << 2) | a
	}

	var thick, out, sh, shSmooth byte = s.Weight, s.OutlineWeight, s.ShadowWeight, s.ShadowBlur
	var r = packLayer(getColor(s.Color))
	var g = packLayer(getColor(s.OutlineColor))
	var b = packLayer(getColor(s.ShadowColor))
	var a = ((thick & 0x03) << 6) | ((out & 0x03) << 4) | ((sh & 0x03) << 2) | (shSmooth & 0x03)
	return rl.NewColor(r, g, b, a)
}
