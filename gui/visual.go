package gui

import (
	"pure-game-kit/graphics"
	p "pure-game-kit/gui/field"
	"pure-game-kit/internal"
	"pure-game-kit/utility/number"
)

func Visual(id string, fields ...string) string {
	return newWidget("visual", id, fields...)
}

//=================================================================
// private

var textBox graphics.TextBox = graphics.TextBox{}
var sprite graphics.Sprite = graphics.Sprite{}
var box graphics.Box = graphics.Box{}

func setupVisualsTextured(root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
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
			box.X, box.Y = widget.X+offX, widget.Y+offY
			box.AssetId = assetId
			box.Color = col
			box.ScaleX, box.ScaleY = 1, 1
			box.Width, box.Height = widget.Width, widget.Height
			box.EdgeLeft, box.EdgeRight = cLeft, cRight
			box.EdgeTop, box.EdgeBottom = cTop, cBottom
			box.PivotX, box.PivotY = 0, 0
		} else {
			sprite.X, sprite.Y = widget.X+offX, widget.Y+offY
			sprite.PivotX, sprite.PivotY = 0, 0
			sprite.AssetId = assetId
			sprite.Color = col
			sprite.ScaleX, sprite.ScaleY = 1, 1
			sprite.TextureRepeat = false
			sprite.Width, sprite.Height = widget.Width, widget.Height
		}

	}
}
func setupVisualsText(root *root, widget *widget, skipEmpty bool) {
	var owner = root.Containers[widget.OwnerId]
	var text = themedProp(p.Text, root, owner, widget)
	if skipEmpty && text == "" {
		return
	}

	var disabled = widget.isDisabled(owner)
	textBox.ScaleX, textBox.ScaleY = 1, 1
	textBox.X, textBox.Y = widget.X, widget.Y
	textBox.EmbeddedColorsTag = '`'
	textBox.EmbeddedAssetsTag = '^'
	textBox.EmbeddedThicknessesTag = '*'
	textBox.Text = text
	textBox.WordWrap = defaultValue(themedProp(p.TextWordWrap, root, owner, widget), "on") == "on"
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.FontId = themedProp(p.TextFontId, root, owner, widget)
	textBox.LineHeight = parseNum(themedProp(p.TextLineHeight, root, owner, widget), 60)
	textBox.LineGap = parseNum(themedProp(p.TextLineGap, root, owner, widget), 0)
	textBox.SymbolGap = parseNum(themedProp(p.TextSymbolGap, root, owner, widget), 0.2)
	textBox.AlignmentX = parseNum(themedProp(p.TextAlignmentX, root, owner, widget), 0)
	textBox.AlignmentY = parseNum(themedProp(p.TextAlignmentY, root, owner, widget), 0)
	textBox.Width, textBox.Height = widget.Width, widget.Height

	textBox.EmbeddedAssetsTag =
		rune(defaultValue(themedProp(p.TextEmbeddedAssetsTag, root, owner, widget), "^")[0])
	textBox.EmbeddedAssetIds = []string{
		themedProp(p.TextEmbeddedAssetId1, root, owner, widget),
		themedProp(p.TextEmbeddedAssetId2, root, owner, widget),
		themedProp(p.TextEmbeddedAssetId3, root, owner, widget),
		themedProp(p.TextEmbeddedAssetId4, root, owner, widget),
		themedProp(p.TextEmbeddedAssetId5, root, owner, widget),
	}

	textBox.EmbeddedColorsTag =
		rune(defaultValue(themedProp(p.TextEmbeddedColorsTag, root, owner, widget), "`")[0])
	textBox.EmbeddedColors = []uint{
		parseColor(themedProp(p.TextEmbeddedColor1, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor2, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor3, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor4, root, owner, widget), disabled),
		parseColor(themedProp(p.TextEmbeddedColor5, root, owner, widget), disabled),
	}

	textBox.EmbeddedThicknessesTag =
		rune(defaultValue(themedProp(p.TextEmbeddedThicknessesTag, root, owner, widget), "*")[0])
	textBox.EmbeddedThicknesses = []float32{
		parseNum(themedProp(p.TextEmbeddedThickness1, root, owner, widget), 0.5),
		parseNum(themedProp(p.TextEmbeddedThickness2, root, owner, widget), 0.5),
		parseNum(themedProp(p.TextEmbeddedThickness3, root, owner, widget), 0.5),
		parseNum(themedProp(p.TextEmbeddedThickness4, root, owner, widget), 0.5),
		parseNum(themedProp(p.TextEmbeddedThickness5, root, owner, widget), 0.5),
	}
}
func drawVisuals(cam *graphics.Camera, root *root, widget *widget, fadeText bool, betweenVisualAndText func()) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = themedProp(p.AssetId, root, owner, widget)
	var col = parseColor(themedProp(p.Color, root, owner, widget), widget.isDisabled(owner))
	var frameCol = parseColor(themedProp(p.FrameColor, root, owner, widget), widget.isDisabled(owner))
	var frameSz = parseNum(themedProp(p.FrameSize, root, owner, widget), 0)

	if assetId != "" {
		var _, has = internal.Boxes[assetId]
		if has {
			cam.DrawBoxes(&box)
		} else {
			cam.DrawSprites(&sprite)
		}

	} else {
		cam.DrawQuad(widget.X, widget.Y, widget.Width, widget.Height, 0, col)
	}

	if textBox.Text != "" {
		var mx, my, mw, mh = cam.MaskX, cam.MaskY, cam.MaskWidth, cam.MaskHeight
		if maskText { // used for inputbox mask
			var x, y = cam.PointToScreen(widget.X+textMargin, widget.Y+textMargin/2)
			var xw, yh = cam.PointToScreen(widget.X+widget.Width-textMargin, widget.Y+widget.Height-textMargin/2)
			xw = number.Limit(xw, cam.MaskX, cam.MaskX+cam.MaskWidth)
			yh = number.Limit(yh, cam.MaskY, cam.MaskY+cam.MaskHeight)
			x = number.Limit(x, cam.MaskX, cam.MaskX+cam.MaskWidth)
			y = number.Limit(y, cam.MaskY, cam.MaskY+cam.MaskHeight)
			cam.Mask(x, y, xw-x, yh-y)
		}

		if betweenVisualAndText != nil {
			betweenVisualAndText()
		}

		var disabled = widget.isDisabled(owner)
		var outlineCol = themedProp(p.TextColorOutline, root, owner, widget)
		if outlineCol != "" {
			var embeddedColors = textBox.EmbeddedColors
			textBox.EmbeddedColors = []uint{}
			textBox.Thickness = parseNum(themedProp(p.TextThicknessOutline, root, owner, widget), 0.92)
			textBox.Smoothness = parseNum(themedProp(p.TextSmoothnessOutline, root, owner, widget), 0.08)
			textBox.Color = parseColor(outlineCol, disabled)
			cam.DrawTextBoxes(&textBox)
			textBox.EmbeddedColors = embeddedColors
		}

		var c = parseColor(defaultValue(themedProp(p.TextColor, root, owner, widget), "0 0 0"), disabled || fadeText)
		textBox.Color = c
		textBox.Thickness = parseNum(themedProp(p.TextThickness, root, owner, widget), 0.5)
		textBox.Smoothness = parseNum(themedProp(p.TextSmoothness, root, owner, widget), 0.02)
		cam.DrawTextBoxes(&textBox)

		cam.Mask(mx, my, mw, mh)
	}

	if frameSz != 0 && frameCol != 0 {
		cam.DrawQuadFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, frameSz, frameCol)
	}

	textBox.Text = "" // skip any further texts unless they are setuped beforehand
}
