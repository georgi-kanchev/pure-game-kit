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

func Object(imageId assets.ImageId, roundness, borderSize float32, borderColor, color uint, area, mask Area, enabled bool) {
	if area == (Area{}) {
		return
	}
	if enabled {
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

func Image(area, mask Area, theme assets.GUIThemeId, enabled bool) {
	var t = getTheme(theme).Image
	var imgId, rnds = assets.ImageId(thNum(t.ImgId)), thNum(t.Rnds)
	var borSz, borCol = thNum(t.BorSz), col.TagHex(thStr(t.BorCol))
	Object(imgId, rnds, borSz, borCol, col.TagHex(thStr(t.Col)), area, mask, enabled)
}
func Label(text string, area, mask Area, theme assets.GUIThemeId, enabled bool) {
	handleText(text, area, mask, internal.GUIText{}, getTheme(theme).Label, enabled, false, false)
}
func Text(text string, area, mask Area, theme assets.GUIThemeId, enabled bool) {
	handleText(text, area, mask, internal.GUIText{}, getTheme(theme).Text, enabled, true, false)
}

func Scrolls(horizontal, vertical *float32, contentWidth, contentHeight float32, area Area, theme assets.GUIThemeId) {
	var t = getTheme(theme)
	var body, hnd = t.Scroll.Body, t.Scroll.Handle
	var scrollSpeed = thNum(hnd.Speed) / Scale
	var size, contentW, contentH = thNum(body.Size) * Scale, contentWidth, contentHeight
	var bodyRound, handleRound = thNum(body.Rnds), thNum(hnd.Rnds)
	var bodyCol = thStr(body.Col)
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
		var col = thStr(hnd.Col)
		var roundness = thNum(body.Rnds)
		handle.X = number.Map(*horizontal, 0, 1, left+handle.Width/2, right-handle.Width/2)
		Object(0, bodyRound, roundness, 0, color.TagHex(bodyCol), hor, Area{}, true)
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thStr(hnd.Col, hnd.Focused.Col)
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize1)
			col = thStr(hnd.Clicked.Col, hnd.Col)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.X = mx // click on scroll body (not handle)
		}
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thStr(hnd.Focused.Col, hnd.Col)
		}
		if IsClicked() || instant {
			handle.X += mdx / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize1)
			col = thStr(hnd.Clicked.Col, hnd.Col)
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
		Object(0, handleRound, 0, 0, color.TagHex(col), handle, Area{}, false)
	}
	if vertical != nil && hasVer {
		var ver = geometry.NewArea(area.X+area.Width/2-size/2, area.Y, size, area.Height)
		var handle = geometry.NewArea(ver.X, 0, size, (area.Height/contentH)*area.Height)
		var top, bot, instant = ver.Y - ver.Height/2, ver.Y + ver.Height/2, false
		var col = thStr(hnd.Col)
		handle.Y = number.Map(*vertical, 0, 1, top+handle.Height/2, bot-handle.Height/2)
		Object(0, bodyRound, 0, 0, color.TagHex(bodyCol), ver, Area{}, true)
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thStr(hnd.Focused.Col, hnd.Col)
		}
		if IsClicked() {
			instant = true // use after widget Shape to account for limiting
			mouse.SetCursor(cursor.Resize2)
			col = thStr(hnd.Clicked.Col, hnd.Col)
		}

		handleInput(handle, Area{}, handleRound)
		if instant && mouse.IsButtonJustPressed(button.Left) {
			handle.Y = my // click on scroll body (not handle)
		}
		if IsFocused() {
			mouse.SetCursor(cursor.Hand)
			col = thStr(hnd.Focused.Col, hnd.Col)
		}
		if IsClicked() || instant {
			handle.Y += mdy / Scale // dragging handle or scroll body after instant click
			mouse.SetCursor(cursor.Resize2)
			col = thStr(hnd.Clicked.Col, hnd.Col)
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
		Object(0, handleRound, 0, 0, color.TagHex(col), handle, Area{}, false)
	}
}
func Button(text string, area, mask Area, theme assets.GUIThemeId, enabled bool) {
	if area == (Area{}) {
		return
	}
	var t = getTheme(theme)
	var body, val = t.Button.Body, t.Button.Value
	var roundness = thNum(body.Rnds)
	var imgId, color = thNum(body.ImgId), thStr(body.Col)
	var borSz, borCol = thNum(body.BorSz), thStr(body.BorCol)
	var interact internal.GUIText
	mask = scaleMask(mask)

	if enabled {
		handleInput(area, mask, roundness)
	} else {
		imgId = thNum(body.Disabled.ImgId, body.ImgId)
		color = thStr(body.Disabled.Col, body.Col)
		borSz = thNum(body.Disabled.BorSz, body.BorSz)
		borCol = thStr(body.Disabled.BorCol, body.BorCol)
		interact = val.Disabled
	}
	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		imgId = thNum(body.Focused.ImgId, body.ImgId)
		color = thStr(body.Focused.Col, body.Col)
		borSz = thNum(body.Focused.BorSz, body.BorSz)
		borCol = thStr(body.Focused.BorCol, body.BorCol)
		interact = val.Focused
	}
	if IsClicked() {
		imgId = thNum(body.Clicked.ImgId, body.ImgId)
		color = thStr(body.Clicked.Col, body.Col)
		borSz = thNum(body.Clicked.BorSz, body.BorSz)
		borCol = thStr(body.Clicked.BorCol, body.BorCol)
		interact = val.Clicked
	}
	Object(assets.ImageId(imgId), roundness, borSz, col.TagHex(borCol), col.TagHex(color), area, mask, false)
	if text != "" {
		handleText(text, area, mask, interact, val.GUIText, false, false, false)
	}
}
func Inputbox(text *string, placeholder string, area, mask Area, theme assets.GUIThemeId, enabled bool) {
	if area == (Area{}) {
		return
	}
	var t, selectionCursorHeight = getTheme(theme), float32(0.85)
	var body, val = t.Inputbox.Body, t.Inputbox.Value
	var sel, cur = t.Inputbox.Selection, t.Inputbox.Cursor
	var plh = t.Inputbox.Placeholder
	var bodyRnds, bodyImg = thNum(body.Rnds), thNum(body.ImgId)
	var bodyBorSz, bodyBorCol = thNum(body.BorSz), thStr(body.BorCol)
	var bodyCol, margin, inter = thStr(body.Col), thStr(val.Margin), internal.GUIText{}
	var mouseInput = mouse.IsAnyButtonJustPressed() || mouse.ScrollX() != 0 || mouse.ScrollY() != 0
	if enabled {
		handleInput(area, scaleMask(mask), bodyRnds)
	} else {
		bodyImg = thNum(body.Disabled.ImgId, body.ImgId)
		bodyBorSz = thNum(body.Disabled.BorSz)
		bodyBorCol = thStr(body.Disabled.BorCol)
		bodyCol = thStr(body.Disabled.Col, body.Col)
		margin = thStr(val.Disabled.Margin, val.Margin)
		inter = val.Disabled
	}

	if IsFocused() {
		mouse.SetCursor(cursor.Input)
		bodyImg = thNum(body.Focused.ImgId, body.ImgId)
		bodyBorSz = thNum(body.Focused.BorSz)
		bodyBorCol = thStr(body.Focused.BorCol)
		bodyCol = thStr(body.Focused.Col, body.Col)
		margin = thStr(val.Focused.Margin, val.Margin)
		inter = val.Focused
	}
	if typingIn == widgetCounter {
		bodyImg = thNum(body.Typing.ImgId, body.ImgId)
		bodyBorSz = thNum(body.Typing.BorSz)
		bodyBorCol = thStr(body.Typing.BorCol)
		bodyCol = thStr(body.Typing.Col, body.Col)
		margin = thStr(val.Typing.Margin, val.Margin)
		inter = val.Typing
	}

	Object(assets.ImageId(bodyImg), bodyRnds, bodyBorSz, col.TagHex(bodyBorCol), col.TagHex(bodyCol), area, scaleMask(mask), false)

	if typingIn == widgetCounter && inputIndexCursor != inputIndexSelection { //#0 typing border
		var selRnds, selImg = thNum(sel.Rnds), assets.ImageId(thNum(sel.ImgId))
		var selBorSz, selBorCol = thNum(sel.BorSz), col.TagHex(thStr(sel.BorCol))
		var selCol = col.TagHex(thStr(sel.Col))
		var selArea = geometry.NewArea(ax+(bx-ax)/2, obj.Y, bx-ax, obj.Height*selectionCursorHeight)
		Object(selImg, selRnds, selBorSz, selBorCol, selCol, selArea, area.Intersect(mask), false)
	} //#

	const valueWidth = 99999 //#1 scroll + text render
	var x = area.X + valueWidth/2 - area.Width/2
	if typingIn == widgetCounter {
		x += inputScroll
	} //#

	var valueArea = geometry.NewArea(x, area.Y, valueWidth, area.Height) //#2 text render
	if *text == "" {
		handleText(placeholder, valueArea, area.Intersect(mask), internal.GUIText{}, plh, false, false, false)
	} else {
		inter.Margin = &margin
		handleText(*text, valueArea, area.Intersect(mask), inter, val.GUIText, false, false, true)
	} //#

	var a, b = min(inputIndexCursor, inputIndexSelection), max(inputIndexCursor, inputIndexSelection) //#3 cursor
	if typingIn == widgetCounter {
		ax, bx = obj.TextCursorPositionAt(a), obj.TextCursorPositionAt(b)
	}

	if IsClicked() {
		var i, closestIndex, x, closestDist = 0, 0, float32(0), float32(valueWidth)
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
	} //#

	if IsFocused() && mouseInput { //#4
		inputCursorTimer, typingIn = 0, widgetCounter
	} else if (!IsFocused() && typingIn == widgetCounter && mouseInput) || !window.IsFocused() {
		typingIn, inputIndexSelection = 0, inputIndexCursor
	}
	if typingIn != lastTypingIn { // no longer typing or switching inputbox while typing
		inputScroll = 0
	}
	if typingIn != widgetCounter {
		return //=================================================================
	} //#

	var inputStr = keyboard.Input() //#5 keyboard input capture & shortcuts
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
	} //#

	if *text == "" { // cannot select placeholder text
		inputIndexCursor, inputIndexSelection = 0, 0
	}

	var cursorX = obj.TextCursorPositionAt(inputIndexCursor) //#6 cursor
	if cursorX > area.X+area.Width/2 {
		inputScroll -= cursorX - (area.X + area.Width/2)
	} else if cursorX < area.X-area.Width/2 {
		inputScroll += (area.X - area.Width/2) - cursorX
	}
	cursorX = number.Limit(cursorX, area.X-area.Width/2, area.X+area.Width/2)
	if inputCursorTimer > 1 {
		inputCursorTimer = 0
	} else if inputCursorTimer < 0.5 {
		var curRnds, curImg = thNum(cur.Rnds), assets.ImageId(thNum(cur.ImgId))
		var curBorSz, curBorCol = thNum(cur.BorSz), col.TagHex(thStr(cur.BorCol))
		var curCol, curWidth = col.TagHex(thStr(cur.Col)), thNum(cur.Width)
		var curArea = geometry.NewArea(cursorX, obj.Y, Scale*curWidth, obj.Height*selectionCursorHeight)
		Object(curImg, curRnds, curBorSz, curBorCol, curCol, curArea, mask, false)
	} //#

	if len(inputStr) > 0 { //#7 typing
		var inputStr = string(inputStr)
		*text = txt.Insert(*text, inputStr, inputIndexCursor)
		inputIndexCursor = number.Limit(inputIndexCursor+1, 0, txt.Length(*text))
		inputCursorTimer, inputIndexSelection = 0, inputIndexCursor
	} //#
}

// Negative step hides the indicators.
func Slider(value *float32, step float32, area, mask Area, theme assets.GUIThemeId, enabled bool) {
	var left, right = area.X - area.Width/2, area.X + area.Width/2
	var x = number.Map(*value, 0, 1, left+area.Height/2, right-area.Height/2)
	var hndArea = geometry.NewArea(x, area.Y, area.Height, area.Height)
	var t, dragging = getTheme(theme), false
	var body, hnd = t.Slider.Body, t.Slider.Hnd
	var tStep = t.Slider.Step
	var bodyCol, hndCol = thStr(body.Col), thStr(hnd.Col)
	var bodyImg, hndImg = thNum(body.ImgId), thNum(hnd.ImgId)
	var bodyRnd, hndRnd = thNum(body.Rnds), thNum(hnd.Rnds)
	var bodyBorSz, hndBorSz = thNum(body.BorSz), thNum(hnd.BorSz)
	var bodyBorCol, hndBorCol = thStr(body.BorCol), thStr(hnd.BorCol)
	mask = scaleMask(mask)

	handleInput(area, mask, bodyRnd)

	if IsFocused() {
		mouse.SetCursor(cursor.Hand)
		hndCol = thStr(hnd.Focused.Col, hnd.Col)
		hndImg = thNum(hnd.Focused.ImgId, hnd.ImgId)
		hndRnd = thNum(hnd.Focused.Rnds, hnd.Rnds)
		hndBorSz = thNum(hnd.Focused.BorSz, hnd.BorSz)
		hndBorCol = thStr(hnd.Focused.BorCol, hnd.BorCol)
		bodyCol = thStr(body.Focused.Col, body.Col)
		bodyImg = thNum(body.Focused.ImgId, body.ImgId)
		bodyRnd = thNum(body.Focused.Rnds, body.Rnds)
		bodyBorSz = thNum(body.Focused.BorSz, body.BorSz)
		bodyBorCol = thStr(body.Focused.BorCol, body.BorCol)
	}
	if IsClicked() {
		mouse.SetCursor(cursor.Resize1)
		dragging, hndCol = true, thStr(hnd.Clicked.Col, hnd.Col)
		hndImg = thNum(hnd.Clicked.ImgId, hnd.ImgId)
		hndRnd = thNum(hnd.Clicked.Rnds, hnd.Rnds)
		hndBorSz = thNum(hnd.Clicked.BorSz, hnd.BorSz)
		hndBorCol = thStr(hnd.Clicked.BorCol, hnd.BorCol)
		bodyCol = thStr(body.Clicked.Col, body.Col)
		bodyImg = thNum(body.Clicked.ImgId, body.ImgId)
		bodyRnd = thNum(body.Clicked.Rnds, body.Rnds)
		bodyBorSz = thNum(body.Clicked.BorSz, body.BorSz)
		bodyBorCol = thStr(body.Clicked.BorCol, body.BorCol)
	}
	if !enabled {
		hndCol = thStr(hnd.Disabled.Col, hnd.Col)
		hndImg = thNum(hnd.Disabled.ImgId, hnd.ImgId)
		hndRnd = thNum(hnd.Disabled.Rnds, hnd.Rnds)
		hndBorSz = thNum(hnd.Disabled.BorSz, hnd.BorSz)
		hndBorCol = thStr(hnd.Disabled.BorCol, hnd.BorCol)
		bodyCol = thStr(body.Disabled.Col, body.Col)
		bodyImg = thNum(body.Disabled.ImgId, body.ImgId)
		bodyRnd = thNum(body.Disabled.Rnds, body.Rnds)
		bodyBorSz = thNum(body.Disabled.BorSz, body.BorSz)
		bodyBorCol = thStr(body.Disabled.BorCol, body.BorCol)
	}

	if step > 0 {
		var stepSize = number.Map(step, 0, 1, area.Height/20, area.Height/2)
		var minX, maxX, stepCol = left + area.Height/2, right - area.Height/2, col.TagHex(thStr(tStep.Col))
		var stepImg, stepRnd = thNum(tStep.ImgId), thNum(tStep.Rnds)
		var stepBorSz, stepBorCol = thNum(tStep.BorSz), col.TagHex(thStr(tStep.BorCol))
		for t := float32(0.0); t <= 1.0+0.001; t += step {
			var stepArea = geometry.NewArea(number.Map(t, 0, 1, minX, maxX), area.Y, stepSize, stepSize)
			Object(assets.ImageId(stepImg), stepRnd, stepBorSz, stepBorCol, stepCol, stepArea, mask, false)
		}
	}
	Object(assets.ImageId(bodyImg), bodyRnd, bodyBorSz, col.TagHex(bodyBorCol), col.TagHex(bodyCol), area, mask, false)
	Object(assets.ImageId(hndImg), hndRnd, hndBorSz, col.TagHex(hndBorCol), col.TagHex(hndCol), hndArea, mask, false)

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

func handleText(text string, area, mask Area, inter, opt internal.GUIText, enabled, wordWrap, isInputbox bool) {
	if area == (Area{}) || text == "" {
		return
	}
	var lineH = thNum(inter.LineH, opt.LineH)
	var fontId, color = thNum(inter.FontId, opt.FontId), thStr(inter.Col, opt.Col)
	var wgt, align = thNum(inter.Wgt, opt.Wgt), thStr(inter.Align, opt.Align, new("0.5 0.5"))
	var gap, mar = thStr(inter.Gap, opt.Gap, new("0 0")), thStr(inter.Margin, opt.Margin, new("0 0"))
	var outSz, outCol = thNum(inter.OutSz, opt.OutSz), thStr(inter.OutCol, opt.OutCol)
	var sWgt, sBlur = thNum(inter.ShWgt, opt.ShWgt), thNum(inter.ShBlur, opt.ShBlur)
	var sCol, sOff = thStr(inter.ShCol, opt.ShCol), thStr(inter.ShOff, opt.ShOff)
	var marX = txt.ToNumber[float32](txt.SplitAtIndex(mar, " ", 0))
	var marY = txt.ToNumber[float32](txt.SplitAtIndex(mar, " ", 1))
	area.Width -= marX
	area.Height -= marY

	if enabled {
		handleInput(area, scaleMask(mask), 0)
	}

	obj.Effects = graphics.Effects(internal.DefaultEffects)
	obj.X, obj.Y, obj.Width, obj.Height, obj.Roundness = area.X, area.Y, area.Width, area.Height, 0
	obj.Effects.TextIsInput, obj.ImageId, obj.Effects.Tint, obj.Mask = isInputbox, 0, palette.White, scaleMask(mask)
	obj.TextFontId, obj.Text, obj.Effects.TextWordWrap = assets.FontId(fontId), text, wordWrap
	obj.Effects.TextLineHeight, obj.Effects.TextColor, obj.Effects.TextWeight = lineH, col.TagHex(color), wgt
	obj.Effects.TextAlignX = txt.ToNumber[float32](txt.SplitAtIndex(align, " ", 0))
	obj.Effects.TextAlignY = txt.ToNumber[float32](txt.SplitAtIndex(align, " ", 1))
	obj.Effects.TextSymbolGap = txt.ToNumber[float32](txt.SplitAtIndex(gap, " ", 0))
	obj.Effects.TextLineGap = txt.ToNumber[float32](txt.SplitAtIndex(gap, " ", 1))
	obj.Effects.OutlineSize, obj.Effects.OutlineColor, obj.Effects.TextShadowBlur = outSz, col.TagHex(outCol), sBlur
	obj.Effects.TextShadowWeight, obj.Effects.TextShadowColor = sWgt, col.TagHex(sCol)
	obj.Effects.TextShadowOffsetX = txt.ToNumber[float32](txt.SplitAtIndex(sOff, " ", 0))
	obj.Effects.TextShadowOffsetY = txt.ToNumber[float32](txt.SplitAtIndex(sOff, " ", 1))
	view.DrawObject(&obj)
}

func scaleMask(mask Area) Area { return mask }
func getTheme(theme assets.GUIThemeId) internal.GUITheme {
	var th, has = internal.GUIThemes[uint16(theme)]
	if !has {
		th = internal.GUIThemes[0]
	}
	return th
}

func thNum[T number.Number](optional *T, fallbacks ...*T) T {
	if optional == nil {
		for _, f := range fallbacks {
			if f != nil {
				return *f
			}
		}
	}
	if optional == nil {
		return 0
	}
	return *optional
}
func thStr(optional *string, fallbacks ...*string) string {
	if optional == nil {
		for _, f := range fallbacks {
			if f != nil {
				return *f
			}
		}
	}
	if optional == nil {
		return ""
	}
	return *optional
}
