package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	k "pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	m "pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
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
	var _, ownerDisabled = owner.Properties[property.Disabled]
	var _, disabled = widget.Properties[property.Disabled]
	var themePress = themedProp(property.ButtonThemeIdPress, root, owner, widget)
	var focus = widget.isFocused(root, cam)

	if focus {
		m.SetCursor(m.CursorHand)

		if disabled || ownerDisabled {
			m.SetCursor(m.CursorNotAllowed)
		}

		var themeHover = themedProp(property.ButtonThemeIdHover, root, owner, widget)
		if themeHover != "" {
			widget.ThemeId = themeHover
		}
		tryPress(m.IsButtonPressed(m.ButtonLeft), m.IsButtonPressedOnce(m.ButtonLeft), themePress, widget)
	}

	if typingIn == nil { // no hotkeys while typing
		var hotkey = key.FromName(themedProp(property.ButtonHotkey, root, owner, widget))
		tryPress(k.IsKeyPressed(hotkey), k.IsKeyPressedOnce(hotkey), themePress, widget)
	}

	setupVisualsTextured(root, widget)
	setupVisualsText(root, widget, true)
	drawVisuals(cam, root, widget, false, nil)
	buttonColor = parseColor(themedProp(property.Color, root, owner, widget), widget.isDisabled(owner))
	widget.ThemeId = prev
}

func tryPress(press, once bool, themePress string, widget *widget) {
	if press && wPressedOn == widget && themePress != "" {
		widget.ThemeId = themePress
	}
	if once {
		wPressedOn = widget
		wPressedAt = seconds.RealRuntime()
	}
}
