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
	var borderColor, tint = color.Hex(th.Image.BorCol), color.Hex(th.Image.Col)
	Object(assets.ImageId(th.Image.ImgId), th.Image.Rnds, th.Image.BorSz, borderColor, tint, area, mask, input)
}
func Label(text string, area, mask Area, theme assets.ThemeId, input bool) {
	var t, b = getTheme(theme), getTheme(0)
	handleText(text, area, mask, internal.GuiText{}, t.Label, b.Label, input, false)
}
func Text(text string, area, mask Area, theme assets.ThemeId, input bool) {
	var t, b = getTheme(theme), getTheme(0)
	handleText(text, area, mask, internal.GuiText{}, t.Text, b.Text, input, true)
}

func Scrolls(horizontal, vertical *float32, contentWidth, contentHeight float32, area Area, theme assets.ThemeId) {
	var th, base = getTheme(theme), getTheme(0)
	var tBody, tHnd, bBody, bHnd = th.Scroll.Body, th.Scroll.Handle, base.Scroll.Body, base.Scroll.Handle
	var scrollSpeed = thField(0, tHnd.Speed, bHnd.Speed) / Scale
	var size, contentW, contentH = thField(0, tBody.Size, bBody.Size) * Scale, contentWidth, contentHeight
	var bodyRound, handleRound = thField(0, tBody.Rnds, bBody.Rnds), thField(0, tHnd.Rnds, bHnd.Rnds)
	var bodyCol = thField("", tBody.Col, bBody.Col)
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
	if horizontal != nil && hasHor {
		var horArea = area
		if vertical != nil && hasVer { // make space for vertical slider
			horArea.Width, horArea.X = horArea.Width-size, horArea.X-size/2
		}
		var hor = geometry.NewArea(horArea.X, horArea.Y+horArea.Height/2-size/2, horArea.Width, size)
		var handle = geometry.NewArea(0, hor.Y, (horArea.Width/contentW)*horArea.Width, size)
		var left, right, instant = hor.X - hor.Width/2, hor.X + hor.Width/2, false
		var col = thField("", tHnd.Col, bHnd.Col)
		var roundness = thField(0, tBody.Rnds, bBody.Rnds)
		handle.X = number.Map(*horizontal, 0, 1, left+handle.Width/2, right-handle.Width/2)
		Object(0, bodyRound, roundness, 0, color.Hex(bodyCol), hor, Area{}, true)
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thField("", tHnd.Col, tHnd.Focused.Col, bHnd.Focused.Col)
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize1)
			col = thField("", tHnd.Clicked.Col, tHnd.Col, bHnd.Clicked.Col)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.X = mx // click on scroll body (not handle)
		}
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thField("", tHnd.Focused.Col, tHnd.Col, bHnd.Focused.Col)
		}
		if IsClicked() || instant {
			handle.X += mdx / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize1)
			col = thField("", tHnd.Clicked.Col, tHnd.Col, bHnd.Clicked.Col)
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
		var col = thField("", tHnd.Col, bHnd.Col)
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Object(0, bodyRound, 0, 0, color.Hex(bodyCol), ver, Area{}, true)
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thField("", tHnd.Focused.Col, tHnd.Col, bHnd.Focused.Col)
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize2)
			col = thField("", tHnd.Clicked.Col, tHnd.Col, bHnd.Clicked.Col)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.Y = my // click on scroll body (not handle)
		}
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thField("", tHnd.Focused.Col, tHnd.Col, bHnd.Focused.Col)
		}
		if IsClicked() || instant {
			handle.Y += mdy / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize2)
			col = thField("", tHnd.Clicked.Col, tHnd.Col, bHnd.Clicked.Col)
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
	var t, b = getTheme(theme), getTheme(0)
	var tBody, tVal, bBody, bVal = t.Button.Body, t.Button.Value, b.Button.Body, b.Button.Value
	var roundness = thField(0, tBody.Rnds, bBody.Rnds)
	var imgId, color = thField(0, tBody.ImgId, bBody.ImgId), thField("", tBody.Col, bBody.Col)
	var borSz, borCol = thField(0, tBody.BorSz, bBody.BorSz), thField("", tBody.BorCol, bBody.BorCol)
	var interact internal.GuiText
	mask = scaleMask(mask)

	_, _ = tVal, bVal

	if input {
		handleInput(area, mask, roundness)
	} else {
		imgId = thField(0, tBody.Disabled.ImgId, tBody.ImgId, bBody.Disabled.ImgId, bBody.ImgId)
		color = thField("", tBody.Disabled.Col, tBody.Col, bBody.Disabled.Col, bBody.Col)
		borSz = thField(0, tBody.Disabled.BorSz, tBody.BorSz, bBody.Disabled.BorSz, bBody.BorSz, 0)
		borCol = thField("", tBody.Disabled.BorCol, tBody.BorCol, bBody.Disabled.BorCol, bBody.BorCol)
		interact = thField(internal.GuiText{}, tVal.Disabled, bVal.Disabled)
	}
	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		imgId = thField(0, tBody.Focused.ImgId, tBody.ImgId, bBody.Focused.ImgId, bBody.ImgId)
		color = thField("", tBody.Focused.Col, tBody.Col, bBody.Focused.Col, bBody.Col)
		borSz = thField(0, tBody.Focused.BorSz, tBody.BorSz, bBody.Focused.BorSz, bBody.BorSz, 0)
		borCol = thField("", tBody.Focused.BorCol, tBody.BorCol, bBody.Focused.BorCol, bBody.BorCol)
		interact = thField(internal.GuiText{}, tVal.Focused, bVal.Focused)
	}
	if IsClicked() {
		imgId = thField(0, tBody.Clicked.ImgId, tBody.ImgId, bBody.Clicked.ImgId, bBody.ImgId)
		color = thField("", tBody.Clicked.Col, tBody.Col, bBody.Clicked.Col, bBody.Col)
		borSz = thField(0, tBody.Clicked.BorSz, tBody.BorSz, bBody.Clicked.BorSz, bBody.BorSz, 0)
		borCol = thField("", tBody.Clicked.BorCol, tBody.BorCol, bBody.Clicked.BorCol, bBody.BorCol)
		interact = thField(internal.GuiText{}, tVal.Clicked, bVal.Clicked)
	}
	Object(assets.ImageId(imgId), roundness, borSz, col.Hex(borCol), col.Hex(color), area, mask, false)
	if text != "" {
		handleText(text, area, mask, interact, tVal.GuiText, bVal.GuiText, false, false)
	}
}
func Inputbox(text *string, placeholder string, area, mask Area, theme assets.ThemeId, input bool) {
	if area == (Area{}) {
		return
	}
	const defaultRoundness, defaultBorderSize float32 = 0, -5
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
	var left, right = area.X - area.Width/2, area.X + area.Width/2
	var x = number.Map(*value, 0, 1, left+area.Height/2, right-area.Height/2)
	var hndArea = geometry.NewArea(x, area.Y, area.Height, area.Height)
	var t, b, dragging = getTheme(theme), getTheme(0), false
	var tBody, tHnd, bBody, bHnd = t.Slider.Body, t.Slider.Hnd, b.Slider.Body, b.Slider.Hnd
	var tStep, bStep = t.Slider.Step, b.Slider.Step
	var bodyCol, hndCol = thField("", tBody.Col, bBody.Col), thField("", tHnd.Col, bHnd.Col)
	var bodyImg, hndImg = thField(0, tBody.ImgId, bBody.ImgId), thField(0, tHnd.ImgId, bHnd.ImgId)
	var bodyRnd, hndRnd = thField(0, tBody.Rnds, bBody.Rnds), thField(0, tHnd.Rnds, bHnd.Rnds)
	var bodyBorSz, hndBorSz = thField(0, tBody.BorSz, bBody.BorSz), thField(0, tHnd.BorSz, bHnd.BorSz)
	var bodyBorCol, hndBorCol = thField("", tBody.BorCol, bBody.BorCol), thField("", tHnd.BorCol, bHnd.BorCol)
	mask = scaleMask(mask)

	handleInput(area, mask, bodyRnd)

	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		hndCol = thField("", tHnd.Focused.Col, tHnd.Col, bHnd.Focused.Col)
		hndImg = thField(0, tHnd.Focused.ImgId, tHnd.ImgId, bHnd.Focused.ImgId)
		hndRnd = thField(0, tHnd.Focused.Rnds, tHnd.Rnds, bHnd.Focused.Rnds)
		hndBorSz = thField(0, tHnd.Focused.BorSz, tHnd.BorSz, tHnd.Focused.BorSz)
		hndBorCol = thField("", tHnd.Focused.BorCol, tHnd.BorCol, tHnd.Focused.BorCol)
		bodyCol = thField("", tBody.Focused.Col, tBody.Col, bBody.Focused.Col)
		bodyImg = thField(0, tBody.Focused.ImgId, tBody.ImgId, bBody.Focused.ImgId)
		bodyRnd = thField(0, tBody.Focused.Rnds, tBody.Rnds, bBody.Focused.Rnds)
		bodyBorSz = thField(0, tBody.Focused.BorSz, tBody.BorSz, bBody.Focused.BorSz)
		bodyBorCol = thField("", tBody.Focused.BorCol, tBody.BorCol, bBody.Focused.BorCol)
	}
	if IsClicked() {
		mouse.SetCursor(cursor.Resize1)
		dragging, hndCol = true, thField("", tHnd.Clicked.Col, tHnd.Col, bHnd.Clicked.Col)
		hndImg = thField(0, tHnd.Clicked.ImgId, tHnd.ImgId, bHnd.Clicked.ImgId)
		hndRnd = thField(0, tHnd.Clicked.Rnds, tHnd.Rnds, bHnd.Clicked.Rnds)
		hndBorSz = thField(0, tHnd.Clicked.BorSz, tHnd.BorSz, tHnd.Clicked.BorSz)
		hndBorCol = thField("", tHnd.Clicked.BorCol, tHnd.BorCol, tHnd.Clicked.BorCol)
		bodyCol = thField("", tBody.Clicked.Col, tBody.Col, bBody.Clicked.Col)
		bodyImg = thField(0, tBody.Clicked.ImgId, tBody.ImgId, bBody.Clicked.ImgId)
		bodyRnd = thField(0, tBody.Clicked.Rnds, tBody.Rnds, bBody.Clicked.Rnds)
		bodyBorSz = thField(0, tBody.Clicked.BorSz, tBody.BorSz, bBody.Clicked.BorSz)
		bodyBorCol = thField("", tBody.Clicked.BorCol, tBody.BorCol, bBody.Clicked.BorCol)
	}
	if !input {
		hndCol = thField("", tHnd.Disabled.Col, tHnd.Col, bHnd.Disabled.Col)
		hndImg = thField(0, tHnd.Disabled.ImgId, tHnd.ImgId, bHnd.Disabled.ImgId)
		hndRnd = thField(0, tHnd.Disabled.Rnds, tHnd.Rnds, bHnd.Disabled.Rnds)
		hndBorSz = thField(0, tHnd.Disabled.BorSz, tHnd.BorSz, tHnd.Disabled.BorSz)
		hndBorCol = thField("", tHnd.Disabled.BorCol, tHnd.BorCol, tHnd.Disabled.BorCol)
		bodyCol = thField("", tBody.Disabled.Col, tBody.Col, bBody.Disabled.Col)
		bodyImg = thField(0, tBody.Disabled.ImgId, tBody.ImgId, bBody.Disabled.ImgId)
		bodyRnd = thField(0, tBody.Disabled.Rnds, tBody.Rnds, bBody.Disabled.Rnds)
		bodyBorSz = thField(0, tBody.Disabled.BorSz, tBody.BorSz, bBody.Disabled.BorSz)
		bodyBorCol = thField("", tBody.Disabled.BorCol, tBody.BorCol, bBody.Disabled.BorCol)
	}

	if step > 0 {
		var stepSize = number.Map(step, 0, 1, area.Height/20, area.Height/2)
		var minX, maxX, stepCol = left + area.Height/2, right - area.Height/2, col.Hex(thField("", tStep.Col, bStep.Col))
		var stepImg, stepRnd = thField(0, tStep.ImgId, bStep.ImgId), thField(0, tStep.Rnds, bStep.Rnds)
		var stepBorSz, stepBorCol = thField(0, tStep.BorSz, bStep.BorSz), col.Hex(thField("", tStep.BorCol, bStep.BorCol))
		for t := float32(0.0); t <= 1.0+0.001; t += step {
			var stepArea = geometry.NewArea(number.Map(t, 0, 1, minX, maxX), area.Y, stepSize, stepSize)
			Object(assets.ImageId(stepImg), stepRnd, stepBorSz, stepBorCol, stepCol, stepArea, mask, false)
		}
	}
	Object(assets.ImageId(bodyImg), bodyRnd, bodyBorSz, col.Hex(bodyBorCol), col.Hex(bodyCol), area, mask, false)
	Object(assets.ImageId(hndImg), hndRnd, hndBorSz, col.Hex(hndBorCol), col.Hex(hndCol), hndArea, mask, false)

	if dragging {
		x, _ = view.MousePosition()
	}
	x = number.Limit(x, left+area.Height/2, right-area.Height/2)
	*value = number.Map(x, left+area.Height/2, right-area.Height/2, 0, 1)
	*value = number.Snap(*value, number.Absolute(step))
}

// private ========================================================

var view, obj = graphics.View{}, graphics.Object{}

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

func handleText(text string, area, mask Area, interact, optional, base internal.GuiText, input, isText bool) {
	if area == (Area{}) || text == "" {
		return
	}
	var lineHeight = thField(0, interact.LineHeight, optional.LineHeight, base.LineHeight)
	var fontId = thField(0, interact.FontId, optional.FontId, base.FontId)
	var color = thField("", interact.Color, optional.Color, base.Color)
	var weight = thField(0, interact.Weight, optional.Weight, base.Weight)
	var align = thField("", interact.Align, optional.Align, base.Align)
	var gap = thField("", interact.Gap, optional.Gap, base.Gap)
	var margin = thField("", interact.Margin, optional.Margin, base.Margin)
	var outSz = thField(0, interact.OutlineSz, optional.OutlineSz, base.OutlineSz)
	var outCol = thField("", interact.OutlineCol, optional.OutlineCol, base.OutlineCol)
	var sWeight = thField(0, interact.ShWeight, optional.ShWeight, base.ShWeight)
	var sCol = thField("", interact.ShColor, optional.ShColor, base.ShColor)
	var sOff = thField("", interact.ShOffset, optional.ShOffset, base.ShOffset)
	var sBlur = thField(0, interact.ShBlur, optional.ShBlur, base.ShBlur)

	if !isText {
		lineHeight = area.Height / float32(txt.SplitCount(text, "\n")) * 0.8
	}
	if input {
		handleInput(area, scaleMask(mask), 0)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, 0
	obj.ImageId, obj.Effects.Tint, obj.Mask = 0, palette.White, scaleMask(mask)
	obj.TextFontId, obj.Text, obj.Effects.TextWordWrap = assets.FontId(fontId), text, isText
	obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.Effects.TextWeight = lineHeight, col.Hex(color), weight
	obj.Effects.TextAlignX = txt.ToNumber[float32](txt.SplitAtIndex(align, " ", 0))
	obj.Effects.TextAlignY = txt.ToNumber[float32](txt.SplitAtIndex(align, " ", 1))
	obj.Effects.TextSymbolGap = txt.ToNumber[float32](txt.SplitAtIndex(gap, " ", 0))
	obj.Effects.TextLineGap = txt.ToNumber[float32](txt.SplitAtIndex(gap, " ", 1))
	obj.Effects.TextMarginX = txt.ToNumber[float32](txt.SplitAtIndex(margin, " ", 0))
	obj.Effects.TextMarginY = txt.ToNumber[float32](txt.SplitAtIndex(margin, " ", 1))
	obj.Effects.OutlineSize, obj.Effects.OutlineColor, obj.Effects.TextShadowBlur = outSz, col.Hex(outCol), sBlur
	obj.Effects.TextShadowWeight, obj.Effects.TextShadowColor = sWeight, col.Hex(sCol)
	obj.Effects.TextShadowOffsetX = txt.ToNumber[int8](txt.SplitAtIndex(sOff, " ", 0))
	obj.Effects.TextShadowOffsetY = txt.ToNumber[int8](txt.SplitAtIndex(sOff, " ", 1))
	view.DrawObject(&obj)
}

func scaleMask(mask Area) Area {
	return geometry.NewArea(mask.X*Scale, mask.Y*Scale, mask.Width*Scale, mask.Height*Scale)
}
func getTheme(theme assets.ThemeId) *internal.GuiTheme {
	var th, has = internal.Themes[uint16(theme)]
	if !has {
		th = internal.Themes[0]
	}
	return &th
}
func thField[T comparable](defaultValue, optional T, fallbacks ...T) T {
	if optional == defaultValue {
		for _, f := range fallbacks {
			if f != defaultValue {
				return f
			}
		}
	}
	return optional
}
