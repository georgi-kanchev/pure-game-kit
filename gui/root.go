package gui

import (
	"encoding/xml"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	k "pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	m "pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
)

type root struct {
	XmlName       xml.Name     `xml:"GUI"`
	XmlContainers []*container `xml:"Container"`
	XmlScale      float32      `xml:"scale,attr"`
	XmlVolume     float32      `xml:"volume,attr"`

	Volume       float32
	ContainerIds []string
	Themes       map[string]*theme
	Containers   map[string]*container
	Widgets      map[string]*widget
}

func (root *root) IsButtonJustClicked(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	if !exists {
		return false
	}

	var owner = root.Containers[widget.OwnerId]
	var hotkey = key.FromName(root.themedField(field.ButtonHotkey, owner, widget))
	var focus = widget.isFocused(root, camera) && wPressedOn == widget
	var input = k.IsKeyJustPressed(hotkey) || (focus && m.IsButtonJustReleased(b.Left))

	return exists && input
}
func (root *root) IsButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	if !exists {
		return false
	}

	var focus = widget.isFocused(root, camera)
	var owner = root.Containers[widget.OwnerId]
	var hotkey = key.FromName(root.themedField(field.ButtonHotkey, owner, widget))
	var first = k.IsKeyJustPressed(hotkey) || (focus && m.IsButtonJustPressed(b.Left))
	var tick = time.RealRuntime() > wPressedAt+0.5
	var inputHold = k.IsKeyPressed(hotkey) || (focus && wPressedOn == widget && m.IsButtonPressed(b.Left))
	var hold = inputHold && condition.TrueEvery(0.1, text.New(";;hold-", buttonId, "-", hotkey)) && tick

	return first || hold
}
