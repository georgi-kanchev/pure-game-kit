package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	btn "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	txt "pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func InputField(id string, properties ...string) string {
	return newWidget("inputField", id, properties...)
}

//=================================================================

func (g *GUI) InputFieldTyping() (inputFieldId string) {
	if typingIn != nil {
		return typingIn.Id
	}
	return ""
}
func (g *GUI) InputFieldStopTyping() {
	typingIn = nil
	scrollX = 0
}

//=================================================================
// private

const cursorWidth float32 = 2

func setupText(margin float32, w *widget, skipEmpty bool) {
	setupVisualsText(w, skipEmpty)
	w.textBox.AlignmentX, w.textBox.AlignmentY = 0, 0
	w.textBox.Width = 99999
	w.textBox.Height = w.Height - margin
	w.textBox.LineHeight = w.Height - margin
	var scroll = condition.If(typingIn == w, scrollX, 0)
	w.textBox.X = w.textBox.X + margin - scroll
	w.textBox.Y += margin / 2
	w.textBox.Text = txt.Remove(w.textBox.Text, "\n")
}
func inputField(w *widget) {
	maskText = true
	if w.textBox == nil {
		w.textBox = &graphics.TextBox{}
	}

	var owner = w.root.Containers[w.OwnerId]
	var margin = parseNum(w.root.themedField(field.InputFieldMargin, owner, w), 10)

	if keyboard.IsComboJustPressed(key.LeftControl, key.A) ||
		keyboard.IsComboJustPressed(key.RightControl, key.A) {
		indexSelect = 0
		indexCursor = len(symbolXs) - 1
	}

	setupText(margin, w, true)
	setupVisualsTextured(w)

	if w.isFocused() {
		mouse.SetCursor(cursor.Input)
	}

	var anyInput = mouse.IsAnyButtonJustPressed() || mouse.Scroll() != 0
	var focused = w.isFocused()
	var meTyping = typingIn == w // each input field should disable its own typing
	var text = txt.Remove(w.root.themedField(field.Text, owner, w), "\n")
	text = internal.RemoveTags(text)

	if meTyping && ((anyInput && !focused) || !window.IsHovered() || keyboard.IsKeyJustPressed(key.Escape)) {
		typingIn = nil
		scrollX = 0
	}
	if mouse.IsButtonJustPressed(btn.Left) && focused {
		if typingIn != w {
			scrollX = 0
		}

		typingIn = w
		w.textBox.Text = text
		calculateXs(w)
	}
	if typingIn == w {
		text = tryInput(text, margin)
		text = tryRemove(text, margin)
		if !tryCycleSelection(text) {
			tryMoveCursor(text, margin)
		}
		tryFocusNextField()

		scrollX = condition.If(txt.Length(text) == 0, 0, scrollX)
		cursorTime += internal.DeltaTime
	}

	var isPlaceholder = false
	if text == "" {
		var placeholder = w.root.themedField(field.InputFieldPlaceholder, owner, w)
		placeholder = txt.Remove(defaultValue(placeholder, "Type..."), "\n")
		setupText(margin, w, false) // don't skip when empty!
		w.textBox.Text = placeholder
		isPlaceholder = true
	}

	queue(margin, w, isPlaceholder)
	cursorTime = condition.If(cursorTime > 1, 0, cursorTime)
	maskText = false
}
func queue(margin float32, w *widget, isPlaceholder bool) {
	if typingIn == w {
		if w.highlight == nil {
			w.highlight = graphics.NewNinePatch("", 0, 0)
		}
		w.highlight.X, w.highlight.Y = w.X-0.5, w.Y-0.5
		w.highlight.Width, w.highlight.Height = w.Width+0.5, w.Height+0.5
		w.highlight.Tint = palette.Azure
		w.highlight.PivotX, w.highlight.PivotY = 0, 0
		w.highlight.Mask = w.root.Containers[w.OwnerId].mask
		w.root.boxes = append(w.root.boxes, w.highlight)
	}

	textMargin = margin
	queueVisuals(w, isPlaceholder, func() {
		if indexCursor == indexSelect || typingIn != w || len(symbolXs) == 0 || indexCursor >= len(symbolXs) {
			return
		}
		var ax = symbolXs[indexCursor] - scrollX
		var bx = symbolXs[indexSelect] - scrollX

		if ax > bx {
			ax, bx = bx, ax
		}

		if w.handle == nil {
			w.handle = graphics.NewSprite("", 0, 0)
		}

		w.handle.X, w.handle.Y = ax, w.textBox.Y+margin/2
		w.handle.Width, w.handle.Height = bx-ax, w.textBox.Height-margin
		w.handle.Tint, w.handle.Mask = palette.Azure, w.textBox.Mask
		w.handle.PivotX, w.handle.PivotY = 0, 0
		w.root.sprites = append(w.root.sprites, w.handle)
	})

	if typingIn == w && cursorTime < 0.5 {
		var x, y = cursorX(margin, w), w.textBox.Y + margin/2
		var cw, ch = cursorWidth, w.textBox.Height - margin

		if w.cursor1 == nil || w.cursor2 == nil {
			w.cursor1 = graphics.NewSprite("", 0, 0)
			w.cursor2 = graphics.NewSprite("", 0, 0)
		}

		w.cursor1.X, w.cursor1.Y = x-cw/2, y-cw/2
		w.cursor1.Width, w.cursor1.Height = cw+cw, ch+cw
		w.cursor1.Tint, w.cursor1.Mask = palette.Azure, w.textBox.Mask
		w.cursor1.PivotX, w.cursor1.PivotY = 0, 0

		w.cursor2.X, w.cursor2.Y = x, y
		w.cursor2.Width, w.cursor2.Height = cw, ch
		w.cursor2.Tint, w.cursor2.Mask = palette.Black, w.textBox.Mask
		w.cursor2.PivotX, w.cursor2.PivotY = 0, 0
		w.root.spritesAbove = append(w.root.spritesAbove, w.cursor1, w.cursor2)
	}
}
