package gui

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
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
func AreaHUD(horizontal, vertical, width, height float32) geometry.Area {
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
	return geometry.Area{X: x, Y: y, Width: width, Height: height}
}

func Label(text string, area, mask geometry.Area) {
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
	obj.X, obj.Y, obj.Mask = area.X, area.Y, mask
	// obj.Text = txt.New(WidgetCounter)

	view.DrawObject(&obj)
}
func Shape(color uint, roundness float32, area, mask geometry.Area) {
	mask = scaleMask(mask)
	update(area, mask, roundness)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, mask, ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = -10, col.Darken(color, 0.25)

	// obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0, 0, false
	// obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, txt.New(WidgetCounter), 20, palette.White

	view.DrawObject(&obj)
}
func Image(imageId assets.ImageId, tint uint, area, mask geometry.Area) {
	mask = scaleMask(mask)
	update(area, mask, 0)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = imageId, tint, 0
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, mask, ""
	view.DrawObject(&obj)
}

func Scrolls(layoutId assets.LayoutId, boxId int, horizontal, vertical *float32) {
	var layout = internal.Layouts[uint32(layoutId)]
	if layout == nil {
		return
	}
	var box = layout.Boxes[boxId]
	var size, area, contentW, contentH = 10 * Scale, layoutId.Box(boxId), box.ContentWidth, box.ContentHeight
	var mx, my = view.MousePosition()
	var mdx, mdy = view.MouseDelta()
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)
	var hovered = area.ContainsPoint(mx, my)
	if hovered && mouse.IsButtonJustPressed(button.Middle) {
		scrollDraggedLayoutId = layoutId
		scrollDraggedBoxId = boxId
	}
	if mouse.IsButtonJustReleased(button.Middle) || !mouse.IsButtonPressed(button.Middle) {
		scrollDraggedLayoutId = 0
		scrollDraggedBoxId = 0
	}
	var dragging = scrollDraggedLayoutId == layoutId && scrollDraggedBoxId == boxId && mouse.IsButtonPressed(button.Middle)
	if horizontal != nil && contentW > area.Width {
		var hor = geometry.Area{X: area.X, Y: area.Y + area.Height/2 - size/2, Width: area.Width, Height: size}
		var handle = geometry.Area{Y: hor.Y, Width: (area.Width / contentW) * area.Width, Height: size}
		var left, right, instant = hor.X - hor.Width/2, hor.X + hor.Width/2, false
		handle.X = number.Map(*horizontal, 0, 1, left+handle.Width/2, right-handle.Width/2)
		Shape(color.RGBA(0, 0, 0, 127), 0, hor, geometry.Area{})
		if IsHovered() {
			mouse.SetCursor(cursor.Hand)
		}
		if IsClicked() {
			instant = true
		}
		Shape(palette.White, 1, handle, geometry.Area{})
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.X = mx
		}
		if IsClicked() || instant {
			handle.X += mdx
		}
		if dragging {
			handle.X -= mdx * (hor.Width - handle.Width) / (contentW - area.Width)
		}
		if shift && hovered {
			handle.X -= mouse.ScrollY() * scrollSpeed
		} else if !shift && hovered {
			handle.X -= mouse.ScrollX() * scrollSpeed
		}
		handle.X = number.Limit(handle.X, left+handle.Width/2, right-handle.Width/2)
		*horizontal = number.Map(handle.X, left+handle.Width/2, right-handle.Width/2, 0, 1)
	}
	if vertical != nil && contentH > area.Height {
		var ver = geometry.Area{X: area.X + area.Width/2 - size/2, Y: area.Y, Width: size, Height: area.Height}
		var handle = geometry.Area{X: ver.X, Width: size, Height: (area.Height / contentH) * area.Height}
		var top, bot, instant = ver.Y - ver.Height/2, ver.Y + ver.Height/2, false
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Shape(color.RGBA(0, 0, 0, 127), 0, ver, geometry.Area{})
		if IsHovered() {
			mouse.SetCursor(cursor.Hand)
		}
		if IsClicked() {
			instant = true
		}
		Shape(palette.White, 1, handle, geometry.Area{})
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.Y = my
		}
		if IsClicked() || instant {
			handle.Y += mdy
		}
		if dragging {
			handle.Y -= mdy * (ver.Height - handle.Height) / (contentH - area.Height)
		}
		if !shift && hovered {
			handle.Y -= mouse.ScrollY() * scrollSpeed
		}
		handle.Y = number.Limit(handle.Y, top+handle.Height/2, bot-handle.Height/2)
		*vertical = number.Map(handle.Y, top+handle.Height/2, bot-handle.Height/2, 0, 1)
	}
}
func Button(text string, area, mask geometry.Area) {
	const roundness = 0.2
	var baseColor = palette.Gray
	var color = baseColor
	mask = scaleMask(mask)

	update(area, mask, roundness)

	if IsFocused() {
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
func Inputbox(text *string, area, mask geometry.Area) {
	const roundness = 0.2
	var color = palette.DarkGray
	mask = scaleMask(mask)
	update(area, mask, roundness)

	if IsFocused() {
		mouse.SetCursor(cursor.Input)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, mask, ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = -10, col.Darken(color, 0.25)
	view.DrawObject(&obj)

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0, 0.5, false
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, *text, area.Height*0.8, palette.White
	obj.X, obj.Y, obj.Mask = area.X, area.Y, mask
	obj.Effects.TextMarginX = 30
	view.DrawObject(&obj)
}

//=================================================================

func IsHovered() bool       { return nowHovered == widgetCounter }
func IsFocused() bool       { return nowFocused == widgetCounter }
func IsJustFocused() bool   { return nowFocused == widgetCounter && lastFocused != widgetCounter }
func IsJustUnfocused() bool { return lastFocused == widgetCounter && nowFocused != widgetCounter }

func IsClicked() bool       { return clickedWidget == widgetCounter }
func IsJustClicked() bool   { return justClickedWidget == widgetCounter }

func IsJustDragged() bool { return IsFocused() && mouse.IsButtonJustPressed(button.Left) }
func IsJustDropped() bool {
	return lastClickedWidget == widgetCounter && !IsClicked() && mouse.IsButtonJustReleased(button.Left)
}
func IsJustDroppedUpon() bool {
	return drag != (geometry.Area{}) && IsHovered() && mouse.IsButtonJustReleased(button.Left)
}
func Drag() geometry.Area {
	if IsJustDragged() {
		drag = widgetArea
	}
	if IsClicked() {
		var mx, my = mouse.CursorDelta()
		drag.X, drag.Y = drag.X+mx/Scale, drag.Y+my/Scale
	}
	if mouse.IsButtonJustReleased(button.Left) {
		drag = geometry.Area{}
	}
	return drag
}

// private ========================================================

var view graphics.View
var obj graphics.Object
var lastUpdateOnFrame uint64
var widgetCounter = 0  // resets every frame, each widget increases it, used for id
var skipUpdate = false // used for internal calls to the widget functions only for drawing (no input)
const scrollSpeed float32 = 20

var nowHovered int  // the latest widget under the mouse on the current frame
var lastHovered int // the hovered widget from the previous frame
var nowFocused int  // the focused widget for interaction on the current frame
var lastFocused int // the focused widget from the previous frame

var lastClickedWidget int   // the widget that was clicked last frame
var clickedWidget int       // the widget that was initially clicked and is being held
var justClickedWidget int   // the widget that completed a full press-and-release cycle this frame
var scrollDraggedLayoutId assets.LayoutId
var scrollDraggedBoxId int
var widgetArea, drag geometry.Area

var inputCursorIndex int

func scaleMask(mask geometry.Area) geometry.Area {
	return geometry.Area{X: mask.X * Scale, Y: mask.Y * Scale, Width: mask.Width * Scale, Height: mask.Height * Scale}
}
func update(area, mask geometry.Area, roundness float32) {
	if skipUpdate {
		return
	}

	if internal.Frame != lastUpdateOnFrame { // frame reset, runs exactly once on the first widget of a new frame
		lastUpdateOnFrame = internal.Frame
		lastFocused = nowFocused

		if nowHovered == lastHovered {
			nowFocused = nowHovered // widget won input last frame AND won input this frame
		} else {
			nowFocused = 0 // focus broken or shifting
		}

		lastClickedWidget = clickedWidget
		justClickedWidget = 0
		if mouse.IsButtonJustPressed(button.Left) {
			clickedWidget = nowFocused //  lock the active widget to whatever is currently hovered
		} else if mouse.IsButtonJustReleased(button.Left) {
			if clickedWidget != 0 && clickedWidget == nowFocused { // same widget we started clicking on?
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
	var maskHor = mx >= mask.X-mask.Width/2 && mx <= mask.X+mask.Width/2
	var maskVer = my >= mask.Y-mask.Height/2 && my <= mask.Y+mask.Height/2
	var maskCheck = mask == (geometry.Area{}) || (mask != (geometry.Area{}) && maskHor && maskVer)
	if shape.ContainsPoint(mx, my) && maskCheck {
		nowHovered = widgetCounter // top-most logic: later widgets naturally overwrite earlier widgets
	}
}
