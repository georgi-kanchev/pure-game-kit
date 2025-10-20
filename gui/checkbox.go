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
	var on = themedProp(field.CheckboxThemeId, root, owner, widget)
	var off = themedProp(field.ThemeId, root, owner, widget)
	var isOff = widget.Properties[field.Value] == ""
	widget.ThemeId = condition.If(isOff, off, on)

	button(cam, root, widget)

	if root.IsButtonClickedOnce(widget.Id, cam) {
		var group = themedProp(field.CheckboxGroup, root, owner, widget)
		widget.Properties[field.Value] = condition.If(isOff, "v", "")
		var soundId = condition.If(isOff, "~on", "~off")
		sound.AssetId = defaultValue(themedProp(field.ButtonSoundPress, root, owner, widget), soundId)
		sound.Volume = root.Volume
		defer sound.Play()

		if group == "" {
			return
		}

		sound.AssetId = "~on"
		for _, w := range root.Widgets {
			var wOwner = root.Containers[w.OwnerId]
			var wGroup = themedProp(field.CheckboxGroup, root, wOwner, w)
			if wGroup == group {
				w.Properties[field.Value] = ""
			}
		}

		widget.Properties[field.Value] = "v"
	}
}
