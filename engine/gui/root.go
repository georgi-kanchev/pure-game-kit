package gui

import (
	"encoding/xml"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	k "pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	m "pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/utility/text"
)

type root struct {
	XmlName       xml.Name     `xml:"GUI"`
	XmlContainers []*container `xml:"Container"`
	XmlScale      string       `xml:"scale,attr"`

	ContainerIds []string
	Themes       map[string]*theme
	Containers   map[string]*container
	Widgets      map[string]*widget
}

func (root *root) ButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	var owner = root.Containers[widget.OwnerId]
	var hotkey = key.FromName(themedProp(property.ButtonHotkey, root, owner, widget))
	var focus = widget.isFocused(root, camera) && wPressedOn == widget
	var input = k.IsKeyPressedOnce(hotkey) || (focus && m.IsButtonReleasedOnce(m.ButtonLeft))

	return exists && input
}
func (root *root) ButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	if !exists {
		return false
	}

	var focus = widget.isFocused(root, camera)
	var owner = root.Containers[widget.OwnerId]
	var hotkey = key.FromName(themedProp(property.ButtonHotkey, root, owner, widget))
	var first = k.IsKeyPressedOnce(hotkey) || (focus && m.IsButtonPressedOnce(m.ButtonLeft))
	var tick = seconds.RealRuntime() > wPressedAt+0.5
	var inputHold = k.IsKeyPressed(hotkey) || (focus && wPressedOn == widget && m.IsButtonPressed(m.ButtonLeft))
	var hold = inputHold && condition.TrueEvery(0.1, text.New(";;hold-", buttonId, "-", hotkey)) && tick

	return first || hold
}
