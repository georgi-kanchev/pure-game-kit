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

func Object(imageId assets.ImageId, roundness, borderSize float32, borderColor, color uint, area, mask Area, input bool) {
	if area == (Area{}) {
		return
	}
	if input {
		handleInput(area, scaleMask(mask), roundness)
	}
	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, roundness
	obj.ImageId, obj.Effects.Tint, obj.Effects.FillColor, obj.Mask, obj.Text = imageId, palette.White, color, scaleMask(mask), ""
	obj.Effects.BorderSize, obj.Effects.BorderColor = borderSize, borderColor
	if imageId != 0 {
		obj.Effects.Tint, obj.Effects.FillColor = color, 0
	}
	view.DrawObject(&obj)
}

func Image(area, mask Area, theme assets.ThemeId, input bool) {
	var th = getTheme(theme)
	var borderColor, tint = color.Hex(th.Image.BorderColor), color.Hex(th.Image.Color)
	Object(assets.ImageId(th.Image.ImageId), th.Image.Roundness, th.Image.BorderSize, borderColor, tint, area, mask, input)
}
func Label(text string, area, mask Area, theme assets.ThemeId, input bool) {
	handleText(text, area, mask, theme, input, false)
}
func Text(text string, area, mask Area, theme assets.ThemeId, input bool) {
	handleText(text, area, mask, theme, input, true)
}

func Scrolls(horizontal, vertical *float32, contentWidth, contentHeight float32, area Area, theme assets.ThemeId) {
	var th, base = getTheme(theme), getTheme(0)
	var tBody, tHnd, bBody, bHnd = th.Scroll.Body, th.Scroll.Handle, base.Scroll.Body, base.Scroll.Handle
	var scrollSpeed = themeField(0, th.Scroll.Handle.Speed, 40) / Scale
	var size, contentW, contentH = themeField(0, th.Scroll.Body.Size, 12) * Scale, contentWidth, contentHeight
	var bodyRound, handleRound = themeField(0, tBody.Roundness, bBody.Roundness), themeField(0, tHnd.Roundness, bHnd.Roundness)
	var bodyCol = themeField("", tBody.Color, bBody.Color)
	var mx, my = view.MousePosition()
	var mdx, mdy = mouse.CursorDelta()
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)
	var hovered, hasHor, hasVer = area.ContainsPoint(mx, my), contentW > area.Width, contentH > area.Height

	if internal.Frame != lastScrollFrame {
		lastScrollFrame, lastScrollHoveredWidget = internal.Frame, scrollHoveredWidget
		scrollHoveredWidget = 0
	}
	if hovered {
		scrollHoveredWidget = widgetCounter
	}
	if lastScrollHoveredWidget == widgetCounter && mouse.IsButtonJustPressed(button.Middle) {
		scrollDraggedWidget = widgetCounter
	}
	if mouse.IsButtonJustReleased(button.Middle) || !mouse.IsButtonPressed(button.Middle) {
		scrollDraggedWidget = 0
	}

	var dragging = scrollDraggedWidget == widgetCounter && mouse.IsButtonPressed(button.Middle)
	var scrolling = lastScrollHoveredWidget == widgetCounter
	if mouse.IsButtonJustReleased(button.Left) || !mouse.IsButtonPressed(button.Left) {
		scrollHandleDragWidget = 0
		scrollBodyHorDragWidget, scrollBodyVerDragWidget = 0, 0
	}
	if horizontal != nil && hasHor {
		var horArea = area
		if vertical != nil && hasVer { // make space for vertical slider
			horArea.Width, horArea.X = horArea.Width-size, horArea.X-size/2
		}
		var hor = geometry.NewArea(horArea.X, horArea.Y+horArea.Height/2-size/2, horArea.Width, size)
		var handle = geometry.NewArea(0, hor.Y, (horArea.Width/contentW)*horArea.Width, size)
		var left, right, instant = hor.X - hor.Width/2, hor.X + hor.Width/2, false
		var col = themeField("", tHnd.Color, bHnd.Color)
		var roundness = themeField(0, tBody.Roundness, bBody.Roundness, 1)
		handle.X = number.Map(*horizontal, 0, 1, left+handle.Width/2, right-handle.Width/2)
		Object(0, bodyRound, roundness, 0, color.Hex(bodyCol), hor, Area{}, true)
		if scrolling && IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = themeField("", tHnd.Color, tHnd.Focused.Color, bHnd.Focused.Color, bHnd.Color)
		}
		if scrolling && IsClicked() {
			instant = true // use after widget Shape to account for limiting
			scrollBodyHorDragWidget = widgetCounter
			mouse.SetCursor(cursor.Resize1)
			col = themeField("", tHnd.Clicked.Color, tHnd.Color, bHnd.Clicked.Color, bHnd.Color)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.X = mx // click on scroll body (not handle)
		}
		if scrolling && IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = themeField("", tHnd.Focused.Color, tHnd.Color, bHnd.Focused.Color, bHnd.Color)
		}
		if scrolling && IsClicked() {
			scrollHandleDragWidget = widgetCounter // activate handle drag for this axis
		}
		if scrollHandleDragWidget == widgetCounter || scrollBodyHorDragWidget == widgetCounter-1 || (scrolling && instant) {
			handle.X += mdx / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize1)
			col = themeField("", tHnd.Clicked.Color, tHnd.Color, bHnd.Clicked.Color, bHnd.Color)
		}
		if dragging { // middle mouse button dragging on parent box
			handle.X -= mdx / Scale * (hor.Width - handle.Width) / (contentW - horArea.Width)
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
		Object(0, handleRound, 0, 0, color.Hex(col), handle, Area{}, false)
	}
	if vertical != nil && hasVer {
		var ver = geometry.NewArea(area.X+area.Width/2-size/2, area.Y, size, area.Height)
		var handle = geometry.NewArea(ver.X, 0, size, (area.Height/contentH)*area.Height)
		var top, bot, instant = ver.Y - ver.Height/2, ver.Y + ver.Height/2, false
		var col = themeField("", tHnd.Color, bHnd.Color)
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Object(0, bodyRound, 0, 0, color.Hex(bodyCol), ver, Area{}, true)
		if scrolling && IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = themeField("", tHnd.Focused.Color, tHnd.Color, bHnd.Focused.Color, bHnd.Color)
		}
		if scrolling && IsClicked() {
			instant = true // use after widget Shape to account for limiting
			scrollBodyVerDragWidget = widgetCounter
			mouse.SetCursor(cursor.Resize2)
			col = themeField("", tHnd.Clicked.Color, tHnd.Color, bHnd.Clicked.Color, bHnd.Color)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.Y = my // click on scroll body (not handle)
		}
		if scrolling && IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = themeField("", tHnd.Focused.Color, tHnd.Color, bHnd.Focused.Color, bHnd.Color)
		}
		if scrolling && IsClicked() {
			scrollHandleDragWidget = widgetCounter // activate handle drag for this axis
		}
		if scrollHandleDragWidget == widgetCounter || scrollBodyVerDragWidget == widgetCounter-1 || (scrolling && instant) {
			handle.Y += mdy / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize2)
			col = themeField("", tHnd.Clicked.Color, tHnd.Color, bHnd.Clicked.Color, bHnd.Color)
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
		Object(0, handleRound, 0, 0, color.Hex(col), handle, Area{}, false)
	}
}
func Button(text string, area, mask Area, theme assets.ThemeId, input bool) {
	if area == (Area{}) {
		return
	}
	var th, base = getTheme(theme), getTheme(0)
	var tBody, tVal, bBody, bVal = th.Button.Body, th.Button.Value, base.Button.Body, base.Button.Value
	var roundness = themeField(0, tBody.Roundness, bBody.Roundness)
	var imgId, color = themeField(0, tBody.ImageId, bBody.ImageId), themeField("", tBody.Color, bBody.Color)
	var borSz, borCol = themeField(0, tBody.BorderSize, bBody.BorderSize), themeField("", tBody.BorderColor, bBody.BorderColor)
	mask = scaleMask(mask)

	_, _ = tVal, bVal

	if input {
		handleInput(area, mask, roundness)
	} else {
		imgId = themeField(0, tBody.Disabled.ImageId, tBody.ImageId, bBody.Disabled.ImageId, bBody.ImageId)
		color = themeField("", tBody.Disabled.Color, tBody.Color, bBody.Disabled.Color, bBody.Color)
		borSz = themeField(0, tBody.Disabled.BorderSize, tBody.BorderSize, bBody.Disabled.BorderSize, bBody.BorderSize, 0)
		borCol = themeField("", tBody.Disabled.BorderColor, tBody.BorderColor, bBody.Disabled.BorderColor, bBody.BorderColor)
	}
	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		imgId = themeField(0, tBody.Focused.ImageId, tBody.ImageId, bBody.Focused.ImageId, bBody.ImageId)
		color = themeField("", tBody.Focused.Color, tBody.Color, bBody.Focused.Color, bBody.Color)
		borSz = themeField(0, tBody.Focused.BorderSize, tBody.BorderSize, bBody.Focused.BorderSize, bBody.BorderSize, 0)
		borCol = themeField("", tBody.Focused.BorderColor, tBody.BorderColor, bBody.Focused.BorderColor, bBody.BorderColor)
	}
	if IsClicked() {
		imgId = themeField(0, tBody.Clicked.ImageId, tBody.ImageId, bBody.Clicked.ImageId, bBody.ImageId)
		color = themeField("", tBody.Clicked.Color, tBody.Color, bBody.Clicked.Color, bBody.Color)
		borSz = themeField(0, tBody.Clicked.BorderSize, tBody.BorderSize, bBody.Clicked.BorderSize, bBody.BorderSize, 0)
		borCol = themeField("", tBody.Clicked.BorderColor, tBody.BorderColor, bBody.Clicked.BorderColor, bBody.BorderColor)
	}
	Object(assets.ImageId(imgId), roundness, borSz, col.Hex(borCol), col.Hex(color), area, mask, false)
	if text != "" {
		Label(text, area, mask, theme, false)
	}
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
		Object(0, 0, 0, 0, palette.Azure, selectArea, area.Intersect(mask), false)
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
		inputboxDeleteRuneRange(text, a, b) // delete selection
		inputIndexCursor, inputIndexSelection, inputCursorTimer = a, a, 0
	} else {
		if keyboard.IsKeyJustPressed(key.Backspace) || keyboard.IsKeyHeld(key.Backspace, 0.5) {
			inputboxDeleteRuneRange(text, inputIndexCursor, inputIndexCursor-1)
			inputIndexCursor = number.Limit(inputIndexCursor-1, 0, txt.Length(*text))
			inputIndexSelection, inputCursorTimer = inputIndexCursor, 0
		} else if keyboard.IsKeyJustPressed(key.Delete) || keyboard.IsKeyHeld(key.Delete, 0.5) {
			inputboxDeleteRuneRange(text, inputIndexCursor, inputIndexCursor+1)
			inputCursorTimer = 0
		}
	}

	if kb.IsKeyJustPressed(key.LeftArrow) || kb.IsKeyHeld(key.LeftArrow, 0.5) {
		inputCursorTimer = 0
		if a == b || kb.IsKeyPressed(key.LeftShift) || kb.IsKeyPressed(key.RightShift) {
			inputIndexCursor = number.Limit(inputIndexCursor-1, 0, txt.Length(*text))
		} else { // instant jump to start when selected
			inputIndexCursor = a
		}
		inputboxTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.RightArrow) || kb.IsKeyHeld(key.RightArrow, 0.5) {
		inputCursorTimer = 0
		if a == b || kb.IsKeyPressed(key.LeftShift) || kb.IsKeyPressed(key.RightShift) {
			inputIndexCursor = number.Limit(inputIndexCursor+1, 0, txt.Length(*text))
		} else { // instant jump to end  when selected
			inputIndexCursor = b
		}
		inputboxTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.UpArrow) || kb.IsKeyJustPressed(key.Home) {
		inputIndexCursor, inputCursorTimer = 0, 0
		inputboxTryShiftSelect()
	} else if kb.IsKeyJustPressed(key.DownArrow) || kb.IsKeyJustPressed(key.End) {
		inputIndexCursor, inputCursorTimer = txt.Length(*text), 0
		inputboxTryShiftSelect()
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
		Object(0, 1, 0, 0, palette.LightGray, geometry.NewArea(cursorX, obj.Y, Scale*8, obj.Height*0.85), mask, false)
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
	var th, x = getTheme(theme), number.Map(*value, 0, 1, left+area.Height/2, right-area.Height/2)
	var bodyImg = themeField(0, th.Image.ImageId, th.Slider.Body.ImageId, 0)
	var handleImg = themeField(0, th.Image.ImageId, th.Slider.Handle.ImageId, 0)
	mask = scaleMask(mask)

	Object(assets.ImageId(bodyImg), 1, 0, 0, palette.DarkGray, area, mask, input)
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
			Object(0, 1, 0, 0, palette.DarkGray, stepArea, mask, false)
		}
	}
	Object(assets.ImageId(handleImg), 1, 0, 0, color, geometry.NewArea(x, area.Y, area.Height, area.Height), mask, false)

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

func inputboxDeleteRuneRange(text *string, start, end int) {
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
func inputboxTryShiftSelect() {
	if !kb.IsKeyPressed(key.LeftShift) && !kb.IsKeyPressed(key.RightShift) {
		inputIndexSelection = inputIndexCursor
	}
}

func handleText(text string, area, mask Area, theme assets.ThemeId, input, isText bool) {
	if area == (Area{}) || text == "" {
		return
	}
	var th = getTheme(theme)
	var t = th.Label
	var lineHeight = t.LineHeight
	if isText {
		t = th.Text
		lineHeight = t.LineHeight
	} else {
		lineHeight = area.Height / float32(txt.SplitCount(text, "\n")) * 0.8
	}
	if input {
		handleInput(area, scaleMask(mask), 0)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, 0
	obj.ImageId, obj.Effects.Tint, obj.Mask = 0, palette.White, scaleMask(mask)
	obj.TextFontId, obj.Text, obj.Effects.TextWordWrap = assets.FontId(t.FontId), text, isText
	obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.Effects.TextWeight = lineHeight, col.Hex(t.Color), t.Weight
	obj.Effects.TextAlignX = txt.ToNumber[float32](txt.SplitAtIndex(t.Align, " ", 0))
	obj.Effects.TextAlignY = txt.ToNumber[float32](txt.SplitAtIndex(t.Align, " ", 1))
	obj.Effects.TextSymbolGap = txt.ToNumber[float32](txt.SplitAtIndex(t.Gap, " ", 0))
	obj.Effects.TextLineGap = txt.ToNumber[float32](txt.SplitAtIndex(t.Gap, " ", 1))
	obj.Effects.TextMarginX = txt.ToNumber[float32](txt.SplitAtIndex(t.Margin, " ", 0))
	obj.Effects.TextMarginY = txt.ToNumber[float32](txt.SplitAtIndex(t.Margin, " ", 1))
	obj.Effects.OutlineSize, obj.Effects.OutlineColor = float32(t.OutlineSize), col.Hex(t.OutlineColor)
	obj.Effects.TextShadowWeight, obj.Effects.TextShadowColor = t.ShadowWeight, col.Hex(t.ShadowColor)
	obj.Effects.TextShadowOffsetX = txt.ToNumber[int8](txt.SplitAtIndex(t.ShadowOffset, " ", 0))
	obj.Effects.TextShadowOffsetY = txt.ToNumber[int8](txt.SplitAtIndex(t.ShadowOffset, " ", 1))
	obj.Effects.TextShadowBlur = t.ShadowBlur
	view.DrawObject(&obj)
}

func scaleMask(mask Area) Area {
	return geometry.NewArea(mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale)
}
func getTheme(theme assets.ThemeId) internal.GuiTheme {
	var th, has = internal.Themes[uint16(theme)]
	if !has {
		th = internal.Themes[0]
	}
	return th
}
func themeField[T comparable](defaultValue, optional T, fallbacks ...T) T {
	if optional == defaultValue {
		for _, f := range fallbacks {
			if f != defaultValue {
				return f
			}
		}
	}
	return optional
}
