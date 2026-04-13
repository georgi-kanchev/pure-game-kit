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

func tryShowTooltip(w *widget, c *container) {
	if !w.hasTooltip && w != tooltipForWidget {
		return
	}
	var hov = w.isFocused()

	if condition.JustTurnedTrue(hov, w.hoverId) {
		tooltipForWidget = w
		tooltipAt = internal.Runtime
		var tooltipId = w.root.themedField(field.TooltipId, c, w)
		tooltip = w.root.Widgets[tooltipId]

		if tooltip != nil {
			var text = w.root.themedField(field.TooltipText, c, w)
			tooltip.Fields[field.Text] = text

			if text != "" {
				mouse.SetCursor(cursor.Hand)
			}
		}
	}
	if w == tooltipForWidget && condition.JustTurnedTrue(!hov, w.unhoverId) {
		tooltipForWidget = nil
		tooltip = nil
	}
}
func queueTooltip(c *container) {
	if tooltip.textBox == nil {
		tooltip.textBox = &graphics.TextBox{}
	}

	defer func() { tooltipWasVisible = tooltipVisible }()

	var owner = c.root.Containers[tooltipForWidget.OwnerId]
	var hidden = tooltip == nil || tooltipForWidget == nil || internal.Runtime < tooltipAt+0.5 ||
		tooltip.Fields[field.Text] == "" || !tooltipForWidget.isFocused()
	tooltipVisible = !hidden

	if !tooltipWasVisible && tooltipVisible {
		sound.AssetId = defaultValue(c.root.themedField(field.TooltipSound, owner, tooltipForWidget), "~popup")
		sound.Volume = c.root.Volume
		sound.Play()
	}

	if hidden {
		return
	}

	var tCamW, tCamH = c.root.cam.Size()
	var width = dynNum(c, tooltip.Fields[field.Width], 500)
	var margin = parseNum(c.root.themedField(field.TooltipMargin, c, tooltip), 50)
	tooltip.Width, tooltip.Height = width-margin, tCamH

	setupVisualsText(tooltip, true)

	var lines = tooltip.textBox.TextLines()
	var lh = tooltip.textBox.LineHeight
	var textH = float32(len(lines)*int(lh+tooltip.textBox.LineGap)) + lh
	tooltip.textBox.Height = textH
	tooltip.textBox.X = tooltipForWidget.X + tooltipForWidget.Width/2 - tooltip.textBox.Width/2
	tooltip.textBox.Y = tooltipForWidget.Y - textH
	tooltip.textBox.X = number.Limit(tooltip.textBox.X, -tCamW/2, tCamW/2-width)
	tooltip.textBox.Y = number.Limit(tooltip.textBox.Y, -tCamH/2, tCamH/2-textH)
	tooltip.X, tooltip.Y = tooltip.textBox.X, tooltip.textBox.Y
	tooltip.Width, tooltip.Height = width, textH

	tooltip.textBox.X += margin / 2

	if tooltip.Y+tooltip.Height > tooltipForWidget.Y+2 { // margin of error 2 pixels
		tooltip.Y = tooltipForWidget.Y + tooltipForWidget.Height
		tooltip.textBox.Y = tooltip.Y
	}

	setupVisualsTextured(tooltip)
	queueVisuals(tooltip, false, nil)
}
