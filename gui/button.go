package gui

import (
	"pure-game-kit/execution/condition"
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

func (g *GUI) IsButtonJustClicked(buttonId string) bool {
	return clickedId == buttonId
}
func (g *GUI) IsButtonClickedAndHeld(buttonId string) bool {
	return clickedAndHeldId == buttonId
}

//=================================================================
// private

func button(w *widget) {
	var owner = w.root.Containers[w.OwnerId]
	var prev = w.ThemeId
	var themePress = w.root.themedField(field.ButtonThemeIdPress, owner, w)
	var focus = w.isFocused()

	if focus {
		m.SetCursor(cursor.Hand)

		if w.isDisabled(owner) {
			m.SetCursor(cursor.NotAllowed)
		}

		var themeHover = w.root.themedField(field.ButtonThemeIdHover, owner, w)
		if themeHover != "" {
			w.ThemeId = themeHover
		}
		var press, justPress = m.IsButtonPressed(b.Left), m.IsButtonJustPressed(b.Left)
		tryPress(press, justPress, btnSounds, themePress, w, owner)
	}

	var hotkeyStr = w.root.themedField(field.ButtonHotkey, owner, w)
	if typingIn == nil {
		if hotkeyStr != "" { // no hotkeys while typing
			var hotkey = key.FromName(hotkeyStr)
			tryPress(k.IsKeyPressed(hotkey), k.IsKeyJustPressed(hotkey), btnSounds, themePress, w, owner)
		}
		if btnSounds && w.root.IsButtonJustClicked(w.Id) {
			sound.AssetId = defaultValue(w.root.themedField(field.ButtonSoundPress, owner, w), "~release")
			sound.Volume = w.root.Volume
			sound.Play()
		}
	}

	if w.root.IsButtonJustClicked(w.Id) { // handling any widgets that this button toggles
		for _, wId := range owner.Widgets {
			var curWidget = w.root.Widgets[wId]
			var toggleParentId = w.root.themedField(field.ToggleButtonId, owner, curWidget)
			if toggleParentId == w.Id {
				var hidden = curWidget.Fields[field.Hidden]
				var newHidden = condition.If(hidden == "", "1", "")
				curWidget.Fields[field.Hidden] = newHidden
			}
		}
	}

	setupVisualsTextured(w)
	setupVisualsText(w, false)
	queueVisuals(w, false, nil)
	buttonColor = parseColor(w.root.themedField(field.Color, owner, w), w.isDisabled(owner))
	w.ThemeId = prev
}

func tryPress(press, once, sounds bool, themePress string, widget *widget, owner *container) {
	if press && widget.root.wPressedOn == widget && themePress != "" {
		widget.ThemeId = themePress
	}
	if once {
		if sounds {
			sound.AssetId = defaultValue(widget.root.themedField(field.ButtonSoundPress, owner, widget), "~press")
			sound.Volume = widget.root.Volume
			sound.Play()
		}
		widget.root.wPressedOn = widget
		widget.root.wPressedAt = time.Running()
	}
}
