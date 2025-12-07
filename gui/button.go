package gui

import (
	"pure-game-kit/execution/condition"
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

func (gui *GUI) IsButtonJustClicked(buttonId string, camera *graphics.Camera) bool {
	return gui.root.IsButtonJustClicked(buttonId, camera)
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
	var themePress = root.themedField(field.ButtonThemeIdPress, owner, widget)
	var focus = widget.isFocused(root, cam)

	if focus {
		m.SetCursor(cursor.Hand)

		if widget.isDisabled(owner) {
			m.SetCursor(cursor.NotAllowed)
		}

		var themeHover = root.themedField(field.ButtonThemeIdHover, owner, widget)
		if themeHover != "" {
			widget.ThemeId = themeHover
		}
		tryPress(m.IsButtonPressed(b.Left), m.IsButtonJustPressed(b.Left), btnSounds, themePress, widget, root, owner)
	}

	var hotkeyStr = root.themedField(field.ButtonHotkey, owner, widget)
	if typingIn == nil && hotkeyStr != "" { // no hotkeys while typing
		var hotkey = key.FromName(hotkeyStr)
		tryPress(k.IsKeyPressed(hotkey), k.IsKeyJustPressed(hotkey), btnSounds, themePress, widget, root, owner)

		if btnSounds && root.IsButtonJustClicked(widget.Id, cam) {
			sound.AssetId = defaultValue(root.themedField(field.ButtonSoundPress, owner, widget), "~release")
			sound.Volume = root.Volume
			sound.Play()
		}
	}

	if root.IsButtonJustClicked(widget.Id, cam) { // handling any widgets that this button toggles
		for _, wId := range owner.Widgets {
			var curWidget = root.Widgets[wId]
			var treeParentId = root.themedField(field.ToggleButtonId, owner, curWidget)
			if treeParentId == widget.Id {
				var hidden = curWidget.Fields[field.Hidden]
				curWidget.Fields[field.Hidden] = condition.If(hidden == "", "1", "")
			}
		}
	}

	setupVisualsTextured(root, widget)
	setupVisualsText(root, widget, true)
	drawVisuals(cam, root, widget, false, nil)
	buttonColor = parseColor(root.themedField(field.Color, owner, widget), widget.isDisabled(owner))
	widget.ThemeId = prev
}

func tryPress(press, once, sounds bool, themePress string, widget *widget, root *root, owner *container) {
	if press && wPressedOn == widget && themePress != "" {
		widget.ThemeId = themePress
	}
	if once {
		if sounds {
			sound.AssetId = defaultValue(root.themedField(field.ButtonSoundPress, owner, widget), "~press")
			sound.Volume = root.Volume
			sound.Play()
		}
		wPressedOn = widget
		wPressedAt = time.RealRuntime()
	}
}
