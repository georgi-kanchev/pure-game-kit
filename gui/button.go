package gui

import (
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	k "pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	m "pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/time"
)

func Button(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

//=================================================================
// getters

func (gui *GUI) IsButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	return gui.root.IsButtonClickedOnce(buttonId, camera)
}
func (gui *GUI) IsButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	return gui.root.IsButtonClickedAndHeld(buttonId, camera)
}

//=================================================================
// private

var wPressedOn *widget
var wPressedAt float32
var buttonColor uint
var btnSounds = true

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
		tryPress(m.IsButtonPressed(b.Left), m.IsButtonPressedOnce(b.Left), btnSounds, themePress, widget, root, owner)
	}

	if typingIn == nil { // no hotkeys while typing
		var hotkey = key.FromName(themedProp(field.ButtonHotkey, root, owner, widget))
		tryPress(k.IsKeyPressed(hotkey), k.IsKeyPressedOnce(hotkey), btnSounds, themePress, widget, root, owner)

		if btnSounds && root.IsButtonClickedOnce(widget.Id, cam) {
			sound.AssetId = defaultValue(themedProp(field.ButtonSoundPress, root, owner, widget), "~release")
			sound.Volume = root.Volume
			sound.Play()
		}
	}

	setupVisualsTextured(root, widget)
	setupVisualsText(root, widget, true)
	drawVisuals(cam, root, widget, false, nil)
	buttonColor = parseColor(themedProp(field.Color, root, owner, widget), widget.isDisabled(owner))
	widget.ThemeId = prev
}

func tryPress(press, once, sounds bool, themePress string, widget *widget, root *root, owner *container) {
	if press && wPressedOn == widget && themePress != "" {
		widget.ThemeId = themePress
	}
	if once {
		if sounds {
			sound.AssetId = defaultValue(themedProp(field.ButtonSoundPress, root, owner, widget), "~press")
			sound.Volume = root.Volume
			sound.Play()
		}
		wPressedOn = widget
		wPressedAt = time.RealRuntime()
	}
}
