package gui

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
)

var Scale float32 = 1

// horizontal/vertical 0..1 screen edge percent
//
// width/height 0..1 = screen edge percent, > 1 = absolute screen pixels
func AreaHUD(horizontal, vertical, width, height float32) assets.Area {
	if width >= 0 && width <= 1 {
		var w, _ = view.Size()
		width = w * width
	}
	if height >= 0 && height <= 1 {
		var _, h = view.Size()
		height = h * height
	}

	view.Zoom = Scale
	width, height = width*Scale, height*Scale
	var tlx, tly = view.PointFromEdge(0, 0)
	var brx, bry = view.PointFromEdge(1, 1)
	var x, y = number.Map(horizontal, 0, 1, tlx+width/2, brx-width/2), number.Map(vertical, 0, 1, tly+height/2, bry-height/2)
	return assets.Area{X: x, Y: y, Width: width, Height: height}
}

func Label(text string, area, mask assets.Area) {
	mask = scaleMask(mask)
	update(area, mask, 0)

	if text == "" {
		return
	}
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0.5, 0.5, false
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, text, area.Height*0.8, palette.White
	obj.X, obj.Y, obj.Mask = area.X, area.Y, graphics.Area(mask)
	// obj.Text = txt.New(WidgetCounter)

	view.DrawObject(&obj)
}
func Shape(color uint, roundness float32, area, mask assets.Area) {
	mask = scaleMask(mask)
	update(area, mask, roundness)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, graphics.Area(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = -10, col.Darken(color, 0.25)

	// obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0, 0, false
	// obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, txt.New(WidgetCounter), 20, palette.White

	view.DrawObject(&obj)
}
func Image(imageId assets.ImageId, tint uint, area, mask assets.Area) {
	mask = scaleMask(mask)
	update(area, mask, 0)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = imageId, tint, 0
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, graphics.Area(mask), ""
	view.DrawObject(&obj)
}

func Scrolls(layoutId assets.LayoutId, boxId int, horizontal, vertical *float32) {
	var layout = internal.Layouts[uint32(layoutId)]
	if layout == nil {
		return
	}
	var box = layout.Boxes[boxId]
	var _, contentH = box.ContentWidth, box.ContentHeight
	var area = layoutId.Box(boxId)
	var size = 10 * Scale
	var mx, my = view.MousePosition()
	var _, mdy = view.MouseDelta()
	var hoveredX = mx > area.X-area.Width/2 && mx < area.X+area.Width/2
	var hoveredY = my > area.Y-area.Height/2 && my < area.Y+area.Height/2
	if contentH > area.Height {
		var ver = assets.Area{X: area.X + area.Width/2 - size/2, Y: area.Y, Width: size, Height: area.Height}
		var handle = assets.Area{X: ver.X, Width: size, Height: (area.Height / contentH) * area.Height}
		var top, bot = ver.Y - ver.Height/2, ver.Y + ver.Height/2
		var instant = false
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Shape(color.RGBA(0, 0, 0, 127), 0, ver, assets.Area{})
		if nowHovered == widgetCounter {
			mouse.SetCursor(cursor.Hand)
		}
		if IsClicked() {
			instant = true
		}
		Shape(palette.White, 1, handle, assets.Area{})
		if IsClicked() {
			handle.Y += mdy
		}
		if instant {
			handle.Y = my
		}
		if hoveredX && hoveredY {
			handle.Y -= mouse.ScrollSmooth()
		}
		handle.Y = number.Limit(handle.Y, top+handle.Height/2, bot-handle.Height/2)
		*vertical = number.Map(handle.Y, top+handle.Height/2, bot-handle.Height/2, 0, 1)
	}
}
func Button(text string, area, mask assets.Area) {
	const roundness = 0.2
	var baseColor = palette.Gray
	var color = baseColor
	mask = scaleMask(mask)

	update(area, mask, roundness)

	if IsHovered() {
		mouse.SetCursor(cursor.Hand)
		color = col.Brighten(baseColor, 0.15)
	}
	if IsClicked() {
		color = col.Darken(color, 0.15)
	}

	skipUpdate = true
	Shape(color, roundness, area, mask)
	Label(text, area, mask)
	skipUpdate = false
}
func Inputbox(text *string, area, mask assets.Area) {
	const roundness = 0.2
	var color = palette.DarkGray
	mask = scaleMask(mask)
	update(area, mask, roundness)

	if IsHovered() {
		mouse.SetCursor(cursor.Input)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, graphics.Area(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = -10, col.Darken(color, 0.25)
	view.DrawObject(&obj)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0, 0.5, false
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, *text, area.Height*0.8, palette.White
	obj.X, obj.Y, obj.Mask = area.X, area.Y, graphics.Area(mask)
	obj.Effects.TextMarginX = 30
	view.DrawObject(&obj)
}

//=================================================================

func IsHovered() bool       { return nowActive == widgetCounter }
func IsJustHovered() bool   { return nowActive == widgetCounter && lastActive != widgetCounter }
func IsJustUnhovered() bool { return lastActive == widgetCounter && nowActive != widgetCounter }

func IsClicked() bool     { return clickedWidget == widgetCounter }
func IsJustClicked() bool { return justClickedWidget == widgetCounter }

func IsJustDragged() bool { return IsHovered() && mouse.IsButtonJustPressed(button.Left) }
func IsJustDropped() bool {
	return lastClickedWidget == widgetCounter && !IsClicked() && mouse.IsButtonJustReleased(button.Left)
}
func IsJustDroppedUpon() bool {
	return drag != (assets.Area{}) && nowHovered == widgetCounter && mouse.IsButtonJustReleased(button.Left)
}
func Drag() assets.Area {
	if IsJustDragged() {
		drag = widgetArea
	}
	if IsClicked() {
		var mx, my = mouse.CursorDelta()
		drag.X, drag.Y = drag.X+mx/Scale, drag.Y+my/Scale
	}
	if mouse.IsButtonJustReleased(button.Left) {
		drag = assets.Area{}
	}
	return drag
}

// private ========================================================

var view graphics.View
var obj graphics.Object
var lastUpdateOnFrame uint64
var widgetCounter = 0  // resets every frame, each widget increases it, used for id
var skipUpdate = false // used for internal calls to the widget functions only for drawing (no input)

var nowHovered int  // the latest widget under the mouse on the current frame
var lastHovered int // the hovered widget from the previous frame
var nowActive int   // the active widget for interaction on the current frame
var lastActive int  // the active widget from the previous frame

var lastClickedWidget int // the widget that was clicked last frame
var clickedWidget int     // the widget that was initially clicked and is being held
var justClickedWidget int // the widget that completed a full press-and-release cycle this frame
var widgetArea, drag assets.Area

var inputCursorIndex int

func scaleMask(mask assets.Area) assets.Area {
	return assets.Area{X: mask.X * Scale, Y: mask.Y * Scale, Width: mask.Width * Scale, Height: mask.Height * Scale}
}
func update(area, mask assets.Area, roundness float32) {
	if skipUpdate {
		return
	}

	if internal.Frame != lastUpdateOnFrame { // frame reset, runs exactly once on the first widget of a new frame
		lastUpdateOnFrame = internal.Frame
		lastActive = nowActive

		if nowHovered == lastHovered {
			nowActive = nowHovered // widget won input last frame AND won input this frame
		} else {
			nowActive = 0 // focus broken or shifting
		}

		lastClickedWidget = clickedWidget
		justClickedWidget = 0
		if mouse.IsButtonJustPressed(button.Left) {
			clickedWidget = nowActive //  lock the active widget to whatever is currently hovered
		} else if mouse.IsButtonJustReleased(button.Left) {
			if clickedWidget != 0 && clickedWidget == nowActive { // same widget we started clicking on?
				justClickedWidget = clickedWidget
			}
			clickedWidget = 0 // clear the lock
		} else if !mouse.IsButtonPressed(button.Left) {
			clickedWidget = 0 // if the button not held, ensure nothing is active
		}

		lastHovered = nowHovered
		nowHovered = 0
		widgetCounter = 0

		view.Zoom = Scale
		mouse.SetCursor(cursor.Arrow)
	}

	widgetCounter++
	widgetArea = area

	var mx, my = view.MousePosition()
	var shape = geometry.NewRoundedRectangle(area.X, area.Y, area.Width, area.Height, 0, roundness)
	var maskHor = mx >= mask.X && mx <= mask.X+mask.Width
	var maskVer = my >= mask.Y && my <= mask.Y+mask.Height
	var maskCheck = mask == (assets.Area{}) || (mask != (assets.Area{}) && maskHor && maskVer)
	if shape.ContainsPoint(mx, my) && maskCheck {
		nowHovered = widgetCounter // top-most logic: later widgets naturally overwrite earlier widgets
	}
}
