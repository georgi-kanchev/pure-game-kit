package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
)

func Checkbox(id string, properties ...string) string {
	return newWidget("checkbox", id, properties...)
}

//=================================================================
// private

func checkbox(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var on = root.themedField(field.CheckboxThemeId, owner, widget)
	var off = root.themedField(field.ThemeId, owner, widget)
	var isOff = widget.Fields[field.Value] == ""
	widget.ThemeId = condition.If(isOff, off, on)

	button(cam, root, widget)

	if root.IsButtonJustClicked(widget.Id, cam) {
		var group = root.themedField(field.CheckboxGroup, owner, widget)
		widget.Fields[field.Value] = condition.If(isOff, "v", "")
		var soundId = condition.If(isOff, "~on", "~off")
		sound.AssetId = defaultValue(root.themedField(field.ButtonSoundPress, owner, widget), soundId)
		sound.Volume = root.Volume
		defer sound.Play()

		if group == "" {
			return
		}

		sound.AssetId = "~on"
		for _, w := range root.Widgets {
			var wOwner = root.Containers[w.OwnerId]
			var wGroup = root.themedField(field.CheckboxGroup, wOwner, w)
			if wGroup == group {
				w.Fields[field.Value] = ""
			}
		}

		widget.Fields[field.Value] = "v"
	}
}
