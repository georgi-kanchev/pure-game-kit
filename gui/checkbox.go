package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui/field"
)

func Checkbox(id string, properties ...string) string {
	return newWidget("checkbox", id, properties...)
}

//=================================================================
// private

func checkbox(widget *widget) {
	var owner = widget.root.Containers[widget.OwnerId]
	var on = widget.root.themedField(field.CheckboxThemeId, owner, widget)
	var off = widget.root.themedField(field.ThemeId, owner, widget)
	var isOff = widget.Fields[field.Value] == ""
	widget.ThemeId = condition.If(isOff, off, on)

	button(widget)

	if !widget.root.IsButtonJustClicked(widget.Id) {
		return
	}

	var group = widget.root.themedField(field.CheckboxGroup, owner, widget)
	widget.Fields[field.Value] = condition.If(isOff, "1", "")
	var soundId = condition.If(isOff, "~on", "~off")
	sound.AssetId = defaultValue(widget.root.themedField(field.ButtonSoundPress, owner, widget), soundId)
	sound.Volume = widget.root.Volume
	defer sound.Play()

	if group != "" {
		sound.AssetId = "~on"
		for _, w := range widget.root.Widgets {
			if w.Id == widget.Id {
				continue
			}

			var wOwner = widget.root.Containers[w.OwnerId]
			var wGroup = widget.root.themedField(field.CheckboxGroup, wOwner, w)
			if wGroup == group {
				w.Fields[field.Value] = ""
				w.tryToggleChildrenVisible(owner)
			}
		}

		widget.Fields[field.Value] = "1"
	}

	widget.tryToggleChildrenVisible(owner)
}

func (w *widget) tryToggleChildrenVisible(owner *container) {
	for _, wId := range owner.Widgets {
		var curWidget = w.root.Widgets[wId]
		var toggleParentId = w.root.themedField(field.ToggleButtonId, owner, curWidget)
		if toggleParentId == w.Id {
			var newHidden = condition.If(w.Fields[field.Value] == "", "1", "")
			curWidget.Fields[field.Hidden] = newHidden
		}
	}
}
