package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/field"
	k "pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	m "pure-kit/engine/input/mouse"
	b "pure-kit/engine/input/mouse/button"
	"pure-kit/engine/input/mouse/cursor"
	"pure-kit/engine/utility/time"
)

func Button(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

//=================================================================
// getters

func (gui *GUI) ButtonOnClickOnce(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedOnce(buttonId, camera)
}
func (gui *GUI) ButtonOnClickAndHold(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedAndHeld(buttonId, camera)
}

//=================================================================
// private

var wPressedOn *widget
var wPressedAt float32
var buttonColor uint

func button(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var prev = widget.ThemeId
	var _, ownerDisabled = owner.Properties[field.Disabled]
	var _, disabled = widget.Properties[field.Disabled]
	var themePress = themedProp(field.ButtonThemeIdPress, root, owner, widget)
	var focus = widget.isFocused(root, cam)

	if focus {
		m.SetCursor(cursor.Hand)

		if disabled || ownerDisabled {
			m.SetCursor(cursor.NotAllowed)
		}

		var themeHover = themedProp(field.ButtonThemeIdHover, root, owner, widget)
		if themeHover != "" {
			widget.ThemeId = themeHover
		}
		tryPress(m.IsButtonPressed(b.Left), m.IsButtonPressedOnce(b.Left), themePress, widget)
	}

	if typingIn == nil { // no hotkeys while typing
		var hotkey = key.FromName(themedProp(field.ButtonHotkey, root, owner, widget))
		tryPress(k.IsKeyPressed(hotkey), k.IsKeyPressedOnce(hotkey), themePress, widget)
	}

	setupVisualsTextured(root, widget)
	setupVisualsText(root, widget, true)
	drawVisuals(cam, root, widget, false, nil)
	buttonColor = parseColor(themedProp(field.Color, root, owner, widget), widget.isDisabled(owner))
	widget.ThemeId = prev
}

func tryPress(press, once bool, themePress string, widget *widget) {
	if press && wPressedOn == widget && themePress != "" {
		widget.ThemeId = themePress
	}
	if once {
		wPressedOn = widget
		wPressedAt = time.RealRuntime()
	}
}
