package gui

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/internal"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
)

var Scale float32 = 1

// horizontal/vertical 0..1 screen edge percent
//
// width/height 0..1 = screen edge percent, > 1 = absolute screen pixels
func AreaHUD(horizontal, vertical, width, height float32) assets.Area {
	view.Zoom = Scale

	if width >= 0 && width <= 1 {
		var w, _ = view.Size()
		width = w * width
	}
	if height >= 0 && height <= 1 {
		var _, h = view.Size()
		height = h * height
	}

	width, height = width*Scale, height*Scale
	var tlx, tly = view.PointFromEdge(0, 0)
	var brx, bry = view.PointFromEdge(1, 1)
	var x = number.Map(horizontal, 0, 1, tlx+width/2, brx-width/2)
	var y = number.Map(vertical, 0, 1, tly+height/2, bry-height/2)
	return assets.Area{X: x, Y: y, Width: width, Height: height}
}

func Label(text string, area, mask assets.Area) {
	view.Zoom = Scale
	mask.X, mask.Y, mask.Width, mask.Height = mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Effects.TextAlignX, obj.Effects.TextAlignY = 0.5, 0.5
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, text, area.Height*0.8, palette.White
	obj.X, obj.Y, obj.Mask = area.X, area.Y, graphics.Area(mask)
	view.DrawObject(&obj)
}
func Shape(color uint, roundness float32, area, mask assets.Area) {
	view.Zoom = Scale
	mask.X, mask.Y, mask.Width, mask.Height = mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, graphics.Area(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = -10, col.Darken(color, 0.2)
	view.DrawObject(&obj)
}
func Image(imageId assets.ImageId, tint uint, area, mask assets.Area) {
	view.Zoom = Scale
	mask.X, mask.Y, mask.Width, mask.Height = mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = imageId, tint, 0
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, graphics.Area(mask), ""
	view.DrawObject(&obj)
}

func Button(text string, color uint, area, mask assets.Area) {
	if isHovered(area, 1) {
		color = col.Brighten(color, 0.2)
	}

	Shape(color, 0.1, area, mask)
}

// private ========================================================

var view graphics.View
var obj graphics.Object

var wasHovered assets.Area

func isHovered(area assets.Area, roundness float32) bool {
	return geometry.NewRoundedRectangle(area.X, area.Y, area.Width, area.Height, 0, roundness).ContainsPoint(view.MousePosition())
}
