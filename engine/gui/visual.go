package gui

import (
	"pure-kit/engine/graphics"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/internal"
)

func Visual(id string, properties ...string) string {
	return newWidget("visual", id, properties...)
}

// #region private

var reusableTextBox graphics.TextBox = graphics.TextBox{}
var reusableSprite graphics.Sprite = graphics.Sprite{}
var reusableNineslice graphics.Box = graphics.Box{}

func visual(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var assetId = themedProp(p.AssetId, root, owner, widget)
	var cLeft = parseNum(dyn(owner, themedProp(p.BoxEdgeLeft, root, owner, widget), "100"), 0)
	var cRight = parseNum(dyn(owner, themedProp(p.BoxEdgeRight, root, owner, widget), "100"), 0)
	var cTop = parseNum(dyn(owner, themedProp(p.BoxEdgeTop, root, owner, widget), "100"), 0)
	var cBottom = parseNum(dyn(owner, themedProp(p.BoxEdgeBottom, root, owner, widget), "100"), 0)
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
		reusableTextBox.ScaleX, reusableTextBox.ScaleY = 1, 1
		reusableTextBox.X, reusableTextBox.Y = widget.X, widget.Y
		reusableTextBox.EmbeddedColorsTag = '`'
		reusableTextBox.EmbeddedAssetsTag = '^'
		reusableTextBox.EmbeddedThicknessesTag = '*'
		reusableTextBox.Text = text
		reusableTextBox.WordWrap = defaultValue(themedProp(p.TextWordWrap, root, owner, widget), "on") == "on"
		reusableTextBox.PivotX, reusableTextBox.PivotY = 0, 0
		reusableTextBox.Width, reusableTextBox.Height = w, h
		reusableTextBox.FontId = themedProp(p.TextFontId, root, owner, widget)
		reusableTextBox.LineHeight = parseNum(themedProp(p.TextLineHeight, root, owner, widget), 60)
		reusableTextBox.LineGap = parseNum(themedProp(p.TextLineGap, root, owner, widget), 0)
		reusableTextBox.SymbolGap = parseNum(themedProp(p.TextSymbolGap, root, owner, widget), 0.2)
		reusableTextBox.AlignmentX = parseNum(themedProp(p.TextAlignmentX, root, owner, widget), 0)
		reusableTextBox.AlignmentY = parseNum(themedProp(p.TextAlignmentY, root, owner, widget), 0)

		reusableTextBox.EmbeddedAssetsTag =
			rune(defaultValue(themedProp(p.TextEmbeddedAssetsTag, root, owner, widget), "^")[0])
		reusableTextBox.EmbeddedAssetIds = []string{
			themedProp(p.TextEmbeddedAssetId1, root, owner, widget),
			themedProp(p.TextEmbeddedAssetId2, root, owner, widget),
			themedProp(p.TextEmbeddedAssetId3, root, owner, widget),
			themedProp(p.TextEmbeddedAssetId4, root, owner, widget),
			themedProp(p.TextEmbeddedAssetId5, root, owner, widget),
		}

		reusableTextBox.EmbeddedColorsTag =
			rune(defaultValue(themedProp(p.TextEmbeddedColorsTag, root, owner, widget), "`")[0])
		reusableTextBox.EmbeddedColors = []uint{
			parseColor(themedProp(p.TextEmbeddedColor1, root, owner, widget)),
			parseColor(themedProp(p.TextEmbeddedColor2, root, owner, widget)),
			parseColor(themedProp(p.TextEmbeddedColor3, root, owner, widget)),
			parseColor(themedProp(p.TextEmbeddedColor4, root, owner, widget)),
			parseColor(themedProp(p.TextEmbeddedColor5, root, owner, widget)),
		}

		reusableTextBox.EmbeddedThicknessesTag =
			rune(defaultValue(themedProp(p.TextEmbeddedThicknessesTag, root, owner, widget), "*")[0])
		reusableTextBox.EmbeddedThicknesses = []float32{
			parseNum(themedProp(p.TextEmbeddedThickness1, root, owner, widget), 0.5),
			parseNum(themedProp(p.TextEmbeddedThickness2, root, owner, widget), 0.5),
			parseNum(themedProp(p.TextEmbeddedThickness3, root, owner, widget), 0.5),
			parseNum(themedProp(p.TextEmbeddedThickness4, root, owner, widget), 0.5),
			parseNum(themedProp(p.TextEmbeddedThickness5, root, owner, widget), 0.5),
		}

		var outlineCol = themedProp(p.TextColorOutline, root, owner, widget)
		if outlineCol != "" {
			reusableTextBox.Thickness = parseNum(themedProp(p.TextThicknessOutline, root, owner, widget), 0.92)
			reusableTextBox.Smoothness = parseNum(themedProp(p.TextSmoothnessOutline, root, owner, widget), 0.08)
			reusableTextBox.Color = parseColor(outlineCol)
			cam.DrawTextBoxes(&reusableTextBox)
		}

		reusableTextBox.Color = parseColor(themedProp(p.TextColor, root, owner, widget))
		reusableTextBox.Thickness = parseNum(themedProp(p.TextThickness, root, owner, widget), 0.5)
		reusableTextBox.Smoothness = parseNum(themedProp(p.TextSmoothness, root, owner, widget), 0.02)
		cam.DrawTextBoxes(&reusableTextBox)
	}
}

// #endregion
