package gui

import (
	"pure-kit/engine/graphics"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/internal"
)

func Visual(id string, properties ...string) string {
	return newWidget("visual", id, properties...)
}

//=================================================================
// private

var reusableTextBox graphics.TextBox = graphics.TextBox{}
var reusableSprite graphics.Sprite = graphics.Sprite{}
var reusableNineslice graphics.Box = graphics.Box{}

func setupVisualsTextured(root *root, widget *widget, owner *container) {
	var assetId = themedProp(p.AssetId, root, owner, widget)

	if assetId != "" {
		var cLeft = parseNum(dyn(owner, themedProp(p.BoxEdgeLeft, root, owner, widget), "100"), 0)
		var cRight = parseNum(dyn(owner, themedProp(p.BoxEdgeRight, root, owner, widget), "100"), 0)
		var cTop = parseNum(dyn(owner, themedProp(p.BoxEdgeTop, root, owner, widget), "100"), 0)
		var cBottom = parseNum(dyn(owner, themedProp(p.BoxEdgeBottom, root, owner, widget), "100"), 0)
		var col = parseColor(themedProp(p.Color, root, owner, widget), widget.isDisabled(owner))
		var _, has = internal.Boxes[assetId]
		var offX = parseNum(dyn(owner, widget.Properties[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(owner, widget.Properties[p.OffsetY], "0"), 0)

		if has {
			reusableNineslice.X, reusableNineslice.Y = widget.X+offX, widget.Y+offY
			reusableNineslice.AssetId = assetId
			reusableNineslice.Color = col
			reusableNineslice.ScaleX, reusableNineslice.ScaleY = 1, 1
			reusableNineslice.Width, reusableNineslice.Height = widget.Width, widget.Height
			reusableNineslice.EdgeLeft = cLeft
			reusableNineslice.EdgeRight = cRight
			reusableNineslice.EdgeTop = cTop
			reusableNineslice.EdgeBottom = cBottom
		} else {
			reusableSprite.X, reusableSprite.Y = widget.X+offX, widget.Y+offY
			reusableSprite.AssetId = assetId
			reusableSprite.Color = col
			reusableSprite.ScaleX, reusableSprite.ScaleY = 1, 1
			reusableSprite.RepeatX, reusableSprite.RepeatY = 1, 1
			reusableSprite.Width, reusableSprite.Height = widget.Width, widget.Height
		}

	}
}
func setupVisualsText(root *root, widget *widget, owner *container) {
	var text = themedProp(p.Text, root, owner, widget)
	if text == "" {
		return
	}

	var disabled = widget.isDisabled(owner)
	var offX = parseNum(dyn(owner, widget.Properties[p.OffsetX], "0"), 0)
	var offY = parseNum(dyn(owner, widget.Properties[p.OffsetY], "0"), 0)
	reusableTextBox.ScaleX, reusableTextBox.ScaleY = 1, 1
	reusableTextBox.X, reusableTextBox.Y = widget.X+offX, widget.Y+offY
	reusableTextBox.EmbeddedColorsTag = '`'
	reusableTextBox.EmbeddedAssetsTag = '^'
	reusableTextBox.EmbeddedThicknessesTag = '*'
	reusableTextBox.Text = text
	reusableTextBox.WordWrap = defaultValue(themedProp(p.TextWordWrap, root, owner, widget), "on") == "on"
	reusableTextBox.PivotX, reusableTextBox.PivotY = 0, 0
	reusableTextBox.FontId = themedProp(p.TextFontId, root, owner, widget)
	reusableTextBox.LineHeight = parseNum(themedProp(p.TextLineHeight, root, owner, widget), 60)
	reusableTextBox.LineGap = parseNum(themedProp(p.TextLineGap, root, owner, widget), 0)
	reusableTextBox.SymbolGap = parseNum(themedProp(p.TextSymbolGap, root, owner, widget), 0.2)
	reusableTextBox.AlignmentX = parseNum(themedProp(p.TextAlignmentX, root, owner, widget), 0)
	reusableTextBox.AlignmentY = parseNum(themedProp(p.TextAlignmentY, root, owner, widget), 0)
	reusableTextBox.Width, reusableTextBox.Height = widget.Width, widget.Height

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
		parseColor(themedProp(p.TextEmbeddedColor1, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor2, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor3, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor4, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor5, root, owner, widget), disabled),
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
		reusableTextBox.Color = parseColor(outlineCol, disabled)
	}

	reusableTextBox.Color = parseColor(themedProp(p.TextColor, root, owner, widget), disabled)
	reusableTextBox.Thickness = parseNum(themedProp(p.TextThickness, root, owner, widget), 0.5)
	reusableTextBox.Smoothness = parseNum(themedProp(p.TextSmoothness, root, owner, widget), 0.02)
}
func drawVisuals(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var assetId = themedProp(p.AssetId, root, owner, widget)
	var col = parseColor(themedProp(p.Color, root, owner, widget), widget.isDisabled(owner))
	var frameCol = parseColor(themedProp(p.FrameColor, root, owner, widget), widget.isDisabled(owner))
	var frameSz = parseNum(themedProp(p.FrameSize, root, owner, widget), 0)
	var text = themedProp(p.Text, root, owner, widget)

	if assetId != "" {
		var _, has = internal.Boxes[assetId]
		if has {
			cam.DrawBoxes(&reusableNineslice)
		} else {
			cam.DrawSprites(&reusableSprite)
		}

	} else {
		cam.DrawRectangle(widget.X, widget.Y, widget.Width, widget.Height, 0, col)
	}

	if text != "" {
		var mx, my, mw, mh = cam.MaskX, cam.MaskY, cam.MaskWidth, cam.MaskHeight
		if maskText {
			var x, y = cam.PointToScreen(widget.X+textMargin, widget.Y+textMargin/2)
			var xw, yh = cam.PointToScreen(widget.X+widget.Width-textMargin, widget.Y+widget.Height-textMargin/2)
			cam.Mask(x, y, xw-x, yh-y)
		}

		var outlineCol = themedProp(p.TextColorOutline, root, owner, widget)
		if outlineCol != "" {
			cam.DrawTextBoxes(&reusableTextBox)
		}

		cam.DrawTextBoxes(&reusableTextBox)
		cam.Mask(mx, my, mw, mh)
	}

	if frameSz != 0 && frameCol != 0 {
		cam.DrawFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, frameSz, frameCol)
	}
}
