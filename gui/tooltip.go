package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

func Tooltip(id string, properties ...string) string {
	return newWidget("tooltip", id, properties...)
}

//=================================================================
// private

func tryShowTooltip(widget *widget, root *root, c *container, cam *graphics.Camera) {
	var hov = widget.isFocused(root, cam)

	if condition.JustTurnedTrue(hov, ";;hoverrr-"+widget.Id) {
		tooltipForWidget = widget
		tooltipAt = internal.Runtime
		var tooltipId = root.themedField(field.TooltipId, c, widget)
		tooltip = root.Widgets[tooltipId]

		if tooltip != nil {
			var text = root.themedField(field.TooltipText, c, widget)
			tooltip.Fields[field.Text] = text

			if text != "" {
				mouse.SetCursor(cursor.Hand)
			}
		}
	}
	if widget == tooltipForWidget && condition.JustTurnedTrue(!hov, ";;unhoverrr-"+widget.Id) {
		tooltipForWidget = nil
		tooltip = nil
	}
}
func drawTooltip(root *root, c *container, cam *graphics.Camera) {
	if tooltip.textBox == nil {
		tooltip.textBox = &graphics.TextBox{}
	}

	defer func() { tooltipWasVisible = tooltipVisible }()

	var owner = root.Containers[tooltipForWidget.OwnerId]
	var hidden = tooltip == nil || tooltipForWidget == nil || internal.Runtime < tooltipAt+0.5 ||
		tooltip.Fields[field.Text] == "" || !tooltipForWidget.isFocused(root, cam)
	tooltipVisible = !hidden

	if !tooltipWasVisible && tooltipVisible {
		sound.AssetId = defaultValue(root.themedField(field.TooltipSound, owner, tooltipForWidget), "~popup")
		sound.Volume = root.Volume
		sound.Play()
	}

	if hidden {
		return
	}

	var camW, camH = cam.Size()
	var width = parseNum(dyn(c, tooltip.Fields[field.Width], "500"), 500)
	var margin = parseNum(root.themedField(field.TooltipMargin, c, tooltip), 50)
	tooltip.Width, tooltip.Height = width-margin, camH

	setupVisualsText(root, tooltip, true)

	var lines = tooltip.textBox.TextLines(cam)
	var lh = tooltip.textBox.LineHeight
	var textH = float32(len(lines)*int(lh+tooltip.textBox.LineGap)) + lh
	tooltip.textBox.Height = textH
	tooltip.textBox.X = tooltipForWidget.X + tooltipForWidget.Width/2 - tooltip.textBox.Width/2
	tooltip.textBox.Y = tooltipForWidget.Y - textH
	tooltip.textBox.X = number.Limit(tooltip.textBox.X, -camW/2, camW/2-width)
	tooltip.textBox.Y = number.Limit(tooltip.textBox.Y, -camH/2, camH/2-textH)
	tooltip.X, tooltip.Y = tooltip.textBox.X, tooltip.textBox.Y
	tooltip.Width, tooltip.Height = width, textH

	tooltip.textBox.X += margin / 2

	if tooltip.Y+tooltip.Height > tooltipForWidget.Y+2 { // margin of error 2 pixels
		tooltip.Y = tooltipForWidget.Y + tooltipForWidget.Height
		tooltip.textBox.Y = tooltip.Y
	}

	setupVisualsTextured(root, tooltip)
	drawVisuals(cam, root, tooltip, false, nil)
}
