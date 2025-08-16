package gui

import (
	"pure-kit/engine/graphics"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/internal"
)

func WidgetVisual(id string, properties ...string) string {
	return newWidget("visual", id, properties...)
}

// #region private

var reusableSprite graphics.Sprite = graphics.Sprite{}
var reusableNineslice graphics.Box = graphics.Box{}

func visual(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var assetId = widget.AssetId
	var cLeft = parseNum(dyn(cam, owner, themedProp(p.BoxEdgeLeft, root, owner, widget), "100"), 0)
	var cRight = parseNum(dyn(cam, owner, themedProp(p.BoxEdgeRight, root, owner, widget), "100"), 0)
	var cTop = parseNum(dyn(cam, owner, themedProp(p.BoxEdgeTop, root, owner, widget), "100"), 0)
	var cBottom = parseNum(dyn(cam, owner, themedProp(p.BoxEdgeBottom, root, owner, widget), "100"), 0)
	var col = parseColor(themedProp(p.Color, root, owner, widget))

	if assetId != "" {
		var _, has = internal.Boxes[assetId]
		if has {
			reusableNineslice.X, reusableNineslice.Y = widget.X, widget.Y
			reusableNineslice.AssetId = assetId
			reusableNineslice.Color = col
			reusableNineslice.ScaleX, reusableNineslice.ScaleY = 1, 1
			reusableNineslice.Width, reusableNineslice.Height = w, h
			reusableNineslice.EdgeLeft = cLeft
			reusableNineslice.EdgeRight = cRight
			reusableNineslice.EdgeTop = cTop
			reusableNineslice.EdgeBottom = cBottom
			cam.DrawBoxes(&reusableNineslice)
		} else {
			reusableSprite.X, reusableSprite.Y = widget.X, widget.Y
			reusableSprite.AssetId = assetId
			reusableSprite.Color = col
			reusableSprite.ScaleX, reusableSprite.ScaleY = 1, 1
			reusableSprite.RepeatX, reusableSprite.RepeatY = 1, 1
			reusableSprite.Width, reusableSprite.Height = w, h
			cam.DrawSprites(&reusableSprite)
		}

	} else {
		cam.DrawRectangle(widget.X, widget.Y, w, h, 0, col)
	}

	var text, _ = widget.Properties[p.Text]
	if text != "" {
		var textBox = graphics.NewTextBox("", widget.X, widget.Y, text)
		textBox.WordWrap = defaultValue(themedProp(p.TextWordWrap, root, owner, widget), "on") == "on"
		textBox.PivotX, textBox.PivotY = 0, 0
		textBox.Width, textBox.Height = w, h
		textBox.FontId = themedProp(p.TextFontId, root, owner, widget)
		textBox.LineHeight = parseNum(themedProp(p.TextLineHeight, root, owner, widget), 60)
		textBox.LineGap = parseNum(themedProp(p.TextLineGap, root, owner, widget), 0)
		textBox.AlignmentX = parseNum(themedProp(p.TextAlignmentX, root, owner, widget), 0)
		textBox.AlignmentY = parseNum(themedProp(p.TextAlignmentY, root, owner, widget), 0)

		var outlineCol = themedProp(p.TextColorOutline, root, owner, widget)
		if outlineCol != "" {
			textBox.Thickness = parseNum(themedProp(p.TextThicknessOutline, root, owner, widget), 0.92)
			textBox.Smoothness = parseNum(themedProp(p.TextSmoothnessOutline, root, owner, widget), 0.08)
			textBox.Color = parseColor(outlineCol)
			cam.DrawTextBoxes(&textBox)
		}

		textBox.Color = parseColor(themedProp(p.TextColor, root, owner, widget))
		textBox.Thickness = parseNum(themedProp(p.TextThickness, root, owner, widget), 0.5)
		textBox.Smoothness = parseNum(themedProp(p.TextSmoothness, root, owner, widget), 0.02)
		cam.DrawTextBoxes(&textBox)
	}
}

// #endregion
