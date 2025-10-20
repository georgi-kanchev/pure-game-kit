package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
)

func Tooltip(id string, properties ...string) string {
	return newWidget("tooltip", id, properties...)
}

//=================================================================
// private

var tooltip *widget
var tooltipForWidget *widget
var tooltipAt float32
var tooltipVisible, tooltipWasVisible = false, false

func tryShowTooltip(widget *widget, root *root, c *container, cam *graphics.Camera) {
	var hov = widget.isFocused(root, cam)

	if condition.TrueOnce(hov, ";;hoverrr-"+widget.Id) {
		tooltipForWidget = widget
		tooltipAt = time.RealRuntime()
		var tooltipId = themedProp(field.TooltipId, root, c, widget)
		tooltip = root.Widgets[tooltipId]

		if tooltip != nil {
			var text = themedProp(field.TooltipText, root, c, widget)
			tooltip.Properties[field.Text] = text

			if text != "" {
				mouse.SetCursor(cursor.Hand)
			}
		}
	}
	if widget == tooltipForWidget && condition.TrueOnce(!hov, ";;unhoverrr-"+widget.Id) {
		tooltipForWidget = nil
		tooltip = nil
	}
}
func drawTooltip(root *root, c *container, cam *graphics.Camera) {
	defer func() { tooltipWasVisible = tooltipVisible }()

	var owner = root.Containers[tooltipForWidget.OwnerId]
	var hidden = tooltip == nil || tooltipForWidget == nil || time.RealRuntime() < tooltipAt+0.5 ||
		tooltip.Properties[field.Text] == "" || !tooltipForWidget.isFocused(root, cam)
	tooltipVisible = !hidden

	if !tooltipWasVisible && tooltipVisible {
		sound.AssetId = defaultValue(themedProp(field.TooltipSound, root, owner, tooltipForWidget), "~popup")
		sound.Volume = root.Volume
		sound.Play()
	}

	if hidden {
		return
	}

	var camW, camH = cam.Size()
	var width = parseNum(dyn(c, tooltip.Properties[field.Width], "500"), 500)
	var margin = parseNum(themedProp(field.TooltipMargin, root, c, tooltip), 50)
	tooltip.Width, tooltip.Height = width-margin, camH

	setupVisualsText(root, tooltip, true)

	var lines = textBox.TextLines()
	var lh = textBox.LineHeight
	var textH = float32(len(lines)*int(lh+textBox.LineGap)) + lh
	textBox.Height = textH
	textBox.X = tooltipForWidget.X + tooltipForWidget.Width/2 - textBox.Width/2
	textBox.Y = tooltipForWidget.Y - textH
	textBox.X = number.Limit(textBox.X, -camW/2, camW/2-width)
	textBox.Y = number.Limit(textBox.Y, -camH/2, camH/2-textH)
	tooltip.X, tooltip.Y, tooltip.Width, tooltip.Height = textBox.X, textBox.Y, width, textH

	textBox.X += margin / 2

	if tooltip.Y+tooltip.Height > tooltipForWidget.Y+2 { // margin of error 2 pixels
		tooltip.Y = tooltipForWidget.Y + tooltipForWidget.Height
		textBox.Y = tooltip.Y
	}

	setupVisualsTextured(root, tooltip)
	drawVisuals(cam, root, tooltip, false, nil)
}
