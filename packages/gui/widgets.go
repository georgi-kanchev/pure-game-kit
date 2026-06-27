package gui

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	kb "pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	txt "pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/window"
)

type Area = geometry.Area

var Scale float32 = 1

// horizontal/vertical 0..1 screen edge percent
//
// width/height 0..1 = screen edge percent, > 1 = absolute screen pixels
func AreaHUD(horizontal, vertical, width, height float32) Area {
	if width >= 0 && width <= 1 {
		var w, _ = view.Size()
		width = w * width
	}
	if height >= 0 && height <= 1 {
		var _, h = view.Size()
		height = h * height
	}

	width, height, view.Zoom = width*Scale, height*Scale, Scale
	var tlx, tly = view.PointFromEdge(0, 0)
	var brx, bry = view.PointFromEdge(1, 1)
	var x, y = number.Map(horizontal, 0, 1, tlx+width/2, brx-width/2), number.Map(vertical, 0, 1, tly+height/2, bry-height/2)
	return geometry.NewArea(x, y, width, height)
}

func Shape(imageId assets.ImageId, color uint, area, mask Area, theme assets.ThemeId, input bool) {
	if area == (Area{}) {
		return
	}
	if input {
		handleInput(area, scaleMask(mask), defaultRoundness)
	}
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, defaultRoundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor, obj.Mask, obj.Text = 0, palette.White, color, scaleMask(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = defaultBorderSize, col.Darken(color, 0.25)
	if imageId != 0 {
		obj.Effects.Tint, obj.Effects.FillColor = color, 0
	}
	view.DrawObject(&obj)
}
func Label(text string, area, mask Area, theme assets.ThemeId, input bool) {
	handleText(text, area.Height/float32(txt.SplitCount(text, "\n"))*0.8, 0.5, 0.5, area, mask, theme, input, false)
}
func Text(text string, lineHeight float32, area, mask Area, theme assets.ThemeId, input bool) {
	handleText(text, lineHeight, 0, 0, area, mask, theme, input, true)
}

func Scrolls(horizontal, vertical *float32, contentWidth, contentHeight float32, area Area, theme assets.ThemeId) {
	var scrollSpeed = 40 / Scale
	var size, contentW, contentH = 12 * Scale, contentWidth, contentHeight
	var mx, my = view.MousePosition()
	var mdx, mdy = mouse.CursorDelta()
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)
	var hovered, hasHor, hasVer = area.ContainsPoint(mx, my), contentW > area.Width, contentH > area.Height
	var scrollId = widgetCounter
	if internal.Frame != lastScrollFrame {
		lastScrollFrame, lastScrollHoveredWidget = internal.Frame, scrollHoveredWidget
		scrollHoveredWidget = 0
	}
	if hovered {
		scrollHoveredWidget = scrollId
	}
	if lastScrollHoveredWidget == scrollId && mouse.IsButtonJustPressed(button.Middle) {
		scrollDraggedWidget = scrollId
	}
	if mouse.IsButtonJustReleased(button.Middle) || !mouse.IsButtonPressed(button.Middle) {
		scrollDraggedWidget = 0
	}

	var dragging = scrollDraggedWidget == scrollId && mouse.IsButtonPressed(button.Middle)
	var scrolling = lastScrollHoveredWidget == scrollId
	if horizontal != nil && hasHor {
		var hor = geometry.NewArea(area.X, area.Y+area.Height/2-size/2, area.Width, size)
		var handle = geometry.NewArea(0, hor.Y, (area.Width/contentW)*area.Width, size)
		var left, right, instant = hor.X - hor.Width/2, hor.X + hor.Width/2, false
		var col = palette.LightGray
		handle.X = number.Map(*horizontal, 0, 1, left+handle.Width/2, right-handle.Width/2)
		Shape(0, color.RGBA(0, 0, 0, 127), hor, Area{}, theme, true)
		if IsHovered() {
			mouse.SetCursor(cursor.Hand)
			col = palette.White
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize1)
		}
		Shape(0, col, handle, Area{}, theme, true)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.X = mx // click on scroll body (not handle)
		}
		if IsClicked() || instant {
			handle.X += mdx / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize1)
		}
		if dragging { // middle mouse button dragging on parent box
			handle.X -= mdx / Scale * (hor.Width - handle.Width) / (contentW - area.Width)
			mouse.SetCursor(cursor.Resize1)
		}
		if scrolling {
			if shift || !hasVer { // no vertical - so can be scrolled
				handle.X -= mouse.ScrollY() * scrollSpeed
			} else { // regular scrolling
				handle.X -= mouse.ScrollX() * scrollSpeed
			}
		}
		handle.X = number.Limit(handle.X, left+handle.Width/2, right-handle.Width/2)
		*horizontal = number.Map(handle.X, left+handle.Width/2, right-handle.Width/2, 0, 1)
	}
	if vertical != nil && hasVer {
		var ver = geometry.NewArea(area.X+area.Width/2-size/2, area.Y, size, area.Height)
		var handle = geometry.NewArea(ver.X, 0, size, (area.Height/contentH)*area.Height)
		var top, bot, instant = ver.Y - ver.Height/2, ver.Y + ver.Height/2, false
		var col = palette.LightGray
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Shape(0, color.RGBA(0, 0, 0, 127), ver, Area{}, theme, true)
		if IsHovered() {
			mouse.SetCursor(cursor.Hand)
			col = palette.White
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize2)
		}
		Shape(0, col, handle, Area{}, theme, true)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.Y = my // click on scroll body (not handle)
		}
		if IsClicked() || instant {
			handle.Y += mdy / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize2)
		}
		if dragging { // middle mouse button dragging on parent box
			handle.Y -= mdy / Scale * (ver.Height - handle.Height) / (contentH - area.Height)
			mouse.SetCursor(cursor.Resize2)
			if horizontal != nil && hasHor {
				mouse.SetCursor(cursor.Move)
			}
		}
		if !shift && scrolling { // regular scrolling
			handle.Y -= mouse.ScrollY() * scrollSpeed
		}
		handle.Y = number.Limit(handle.Y, top+handle.Height/2, bot-handle.Height/2)
		*vertical = number.Map(handle.Y, top+handle.Height/2, bot-handle.Height/2, 0, 1)
	}
}
func Button(area, mask Area, theme assets.ThemeId, input bool) {
	if area == (Area{}) {
		return
	}
	const roundness = 0.2
	var baseColor = palette.Gray
	var color = baseColor
	mask = scaleMask(mask)

	if input {
		handleInput(area, mask, roundness)
	}
	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		color = col.Brighten(baseColor, 0.15)
	}
	if IsClicked() {
		color = col.Darken(color, 0.15)
	}
	Shape(0, color, area, mask, theme, false)
}
func Inputbox(text *string, placeholder string, area, mask Area, theme assets.ThemeId, input bool) {
	if area == (Area{}) {
		return
	}
	var marginX float32 = 20 * Scale
	var color = palette.DarkGray
	var borderColor = col.Darken(color, 0.25)
	var mouseInput = mouse.IsAnyButtonJustPressed() || mouse.ScrollX() != 0 || mouse.ScrollY() != 0
	if input {
		handleInput(area, scaleMask(mask), defaultRoundness)
	}

	if IsFocused() {
		mouse.SetCursor(cursor.Input)
	}
	if typingIn == widgetCounter {
		borderColor = palette.Gray
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = area.Width, area.Height, 0, defaultRoundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, color
	obj.X, obj.Y, obj.Mask, obj.Text = area.X, area.Y, scaleMask(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = defaultBorderSize, borderColor
	view.DrawObject(&obj)
	area.Width -= marginX

	if typingIn == widgetCounter && inputIndexCursor != inputIndexSelection {
		var selectArea = geometry.NewArea(ax+(bx-ax)/2, obj.Y, bx-ax, obj.Height*0.85)
		Shape(0, palette.Azure, selectArea, area.Intersect(mask), theme, false)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = 0, 0.5, false
	obj.Width, obj.Height, obj.Effects.FillColor, obj.Roundness = 99999, area.Height, 0, 0
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor = 0, palette.White, 0
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, *text, area.Height*0.8, palette.White
	obj.Effects.TextIsInput, obj.Effects.TextMarginX = true, marginX
	var x = area.X + obj.Width/2 - area.Width/2 - marginX/2
	if typingIn == widgetCounter {
		x += inputScroll
	}
	obj.X, obj.Y, obj.Mask = x, area.Y, scaleMask(area.Intersect(mask))
	if obj.Text == "" {
		obj.Text = placeholder
		obj.Effects.TextColor = col.RGBA(40, 40, 40, 255)
		obj.Effects.TextShadowColor = 0
	}
	view.DrawObject(&obj)

	var a, b = min(inputIndexCursor, inputIndexSelection), max(inputIndexCursor, inputIndexSelection)
	if typingIn == widgetCounter {
		ax, bx = obj.TextCursorPositionAt(a), obj.TextCursorPositionAt(b)
	}

	if IsClicked() {
		var i, closestIndex int
		var x, closestDist float32 = 0, 99999
		var mx, _ = view.MousePosition()
		for {
			x = obj.TextCursorPositionAt(i)
			if number.IsNaN(x) {
				break
			}
			var dist = number.Absolute(mx - x)
			if dist < closestDist {
				closestDist, closestIndex = dist, i
			}
			i++
		}
		inputIndexCursor, inputCursorTimer = closestIndex, 0
		if mouse.IsButtonJustPressed(button.Left) {
			inputIndexSelection = closestIndex
		}
	}

	if IsFocused() && mouseInput {
		inputCursorTimer, typingIn = 0, widgetCounter
	} else if (!IsFocused() && typingIn == widgetCounter && mouseInput) || !window.IsFocused() {
		typingIn, inputIndexSelection = 0, inputIndexCursor

	}
	if typingIn != lastTypingIn { // no longer typing or switching inputbox while typing
		inputScroll = 0
	}
	if typingIn != widgetCounter {
		return
	}

	var inputStr = keyboard.Input()
	if a != b && (len(inputStr) > 0 || keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyJustPressed(key.Delete)) {
		inputDeleteRuneRange(text, a, b) // delete selection
		inputIndexCursor, inputIndexSelection, inputCursorTimer = a, a, 0
	} else {
		if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyHeld(key.Backspace, 0.5) {
			inputDeleteRuneRange(text, inputIndexCursor, inputIndexCursor-1)
			inputIndexCursor = number.Limit(inputIndexCursor-1, 0, txt.Length(*text))
			inputIndexSelection, inputCursorTimer = inputIndexCursor, 0
		} else if keyboard.IsKeyJustPressed(key.Delete) || keyboard.IsKeyHeld(key.Delete, 0.5) {
			inputDeleteRuneRange(text, inputIndexCursor, inputIndexCursor+1)
			inputCursorTimer = 0
		}
	}

	if kb.IsKeyJustPressed(key.LeftArrow) || kb.IsKeyHeld(key.LeftArrow, 0.5) {
		inputCursorTimer = 0
		if a == b {
			inputIndexCursor = number.Limit(inputIndexCursor-1, 0, txt.Length(*text))
		} else {
			inputIndexCursor = a
		}
		inputTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.RightArrow) || kb.IsKeyHeld(key.RightArrow, 0.5) {
		inputCursorTimer = 0
		if a == b {
			inputIndexCursor = number.Limit(inputIndexCursor+1, 0, txt.Length(*text))
		} else {
			inputIndexCursor = b
		}
		inputTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.UpArrow) || kb.IsKeyJustPressed(key.Home) {
		inputIndexCursor, inputCursorTimer = 0, 0
		inputTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.DownArrow) || kb.IsKeyJustPressed(key.End) {
		inputIndexCursor, inputCursorTimer = txt.Length(*text), 0
		inputTryShiftSelect()
	} else if kb.IsComboJustPressed(key.LeftControl, key.A) || kb.IsComboJustPressed(key.RightControl, key.A) {
		inputIndexCursor, inputIndexSelection = txt.Length(*text), 0
	}

	if *text == "" { // cannot select placeholder text
		inputIndexCursor, inputIndexSelection = 0, 0
	}

	var cursorX = obj.TextCursorPositionAt(inputIndexCursor)
	if cursorX > area.X+area.Width/2 {
		inputScroll -= cursorX - (area.X + area.Width/2)
	} else if cursorX < area.X-area.Width/2 {
		inputScroll += (area.X - area.Width/2) - cursorX
	}
	cursorX = number.Limit(cursorX, area.X-area.Width/2, area.X+area.Width/2)
	if inputCursorTimer > 1 {
		inputCursorTimer = 0
	} else if inputCursorTimer < 0.5 {
		Shape(0, palette.LightGray, geometry.NewArea(cursorX, obj.Y, Scale*8, obj.Height*0.8), mask, theme, false)
	}

	if len(inputStr) > 0 {
		var inputStr = string(inputStr)
		*text = txt.Insert(*text, inputStr, inputIndexCursor)
		inputIndexCursor = number.Limit(inputIndexCursor+1, 0, txt.Length(*text))
		inputCursorTimer, inputIndexSelection = 0, inputIndexCursor
	}
}

// Negative step hides the indicators.
func Slider(value *float32, step float32, area, mask Area, theme assets.ThemeId, input bool) {
	var baseColor = palette.Gray
	var color, dragging = baseColor, false
	var left, right = area.X - area.Width/2, area.X + area.Width/2
	var x = number.Map(*value, 0, 1, left+area.Height/2, right-area.Height/2)
	mask = scaleMask(mask)

	Shape(0, palette.DarkGray, area, mask, theme, input)
	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		color = col.Brighten(baseColor, 0.15)
	}
	if IsClicked() {
		color, dragging = col.Darken(color, 0.15), true
		mouse.SetCursor(cursor.Resize1)
	}

	if step > 0 {
		var stepSize = number.Map(step, 0, 1, area.Height/20, area.Height/2)
		var minX, maxX = left + area.Height/2, right - area.Height/2
		for t := float32(0.0); t <= 1.0+0.001; t += step {
			var stepArea = geometry.NewArea(number.Map(t, 0, 1, minX, maxX), area.Y, stepSize, stepSize)
			Shape(0, palette.DarkGray, stepArea, mask, theme, false)
		}
	}
	Shape(0, color, geometry.NewArea(x, area.Y, area.Height, area.Height), mask, theme, false)

	if dragging {
		x, _ = view.MousePosition()
	}
	x = number.Limit(x, left+area.Height/2, right-area.Height/2)
	*value = number.Map(x, left+area.Height/2, right-area.Height/2, 0, 1)
	*value = number.Snap(*value, number.Absolute(step))
}

// private ========================================================

var view, obj = graphics.View{}, graphics.Object{}

const defaultRoundness, defaultBorderSize float32 = 0, -5

func scaleMask(mask Area) Area {
	return geometry.NewArea(mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale)
}
func inputDeleteRuneRange(text *string, start, end int) {
	if text == nil || *text == "" {
		return
	}

	var runes = []rune(*text)
	var totalRunes = len(runes)
	if start > end {
		start, end = end, start
	}
	start, end = max(start, 0), min(end, totalRunes)
	if start >= totalRunes {
		return // invalid range or nothing to delete
	}

	runes = append(runes[:start], runes[end:]...) // delete the range in-place
	*text = string(runes)                         // update the underlying string
}
func inputTryShiftSelect() {
	if !kb.IsKeyPressed(key.LeftShift) && !kb.IsKeyPressed(key.RightShift) {
		inputIndexSelection = inputIndexCursor
	}
}

func handleText(text string, lineHeight, ax, ay float32, area, mask Area, theme assets.ThemeId, input, wordWrap bool) {
	if area == (Area{}) || text == "" {
		return
	}
	if input {
		handleInput(area, scaleMask(mask), 0)
	}
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, 0
	obj.Effects.TextAlignX, obj.Effects.TextAlignY, obj.Effects.TextWordWrap = ax, ay, wordWrap
	obj.ImageId, obj.Effects.Tint, obj.Mask = 0, palette.White, scaleMask(mask)
	obj.TextFontId, obj.Text, obj.Effects.TextLineHeight, obj.Effects.TextColor = 0, text, lineHeight, palette.White
	view.DrawObject(&obj)
}
