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

func setupVisualsTextured(root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = root.themedField(p.AssetId, owner, widget)

	if assetId != "" {
		var cLeft = parseNum(root.themedField(p.BoxEdgeLeft, owner, widget), 0)
		var cRight = parseNum(root.themedField(p.BoxEdgeRight, owner, widget), 0)
		var cTop = parseNum(root.themedField(p.BoxEdgeTop, owner, widget), 0)
		var cBottom = parseNum(root.themedField(p.BoxEdgeBottom, owner, widget), 0)
		var col = parseColor(root.themedField(p.Color, owner, widget), widget.isDisabled(owner))
		var _, has = internal.Boxes[assetId]

		if has {
			box.X, box.Y = widget.X, widget.Y
			box.AssetId = assetId
			box.Tint = col
			box.ScaleX, box.ScaleY = 1, 1
			box.Width, box.Height = widget.Width, widget.Height
			box.EdgeLeft, box.EdgeRight = cLeft, cRight
			box.EdgeTop, box.EdgeBottom = cTop, cBottom
			box.PivotX, box.PivotY = 0, 0
		} else {
			sprite.X, sprite.Y = widget.X, widget.Y
			sprite.PivotX, sprite.PivotY = 0, 0
			sprite.AssetId = assetId
			sprite.Tint = col
			sprite.ScaleX, sprite.ScaleY = 1, 1
			sprite.TextureRepeat = false
			sprite.Width, sprite.Height = widget.Width, widget.Height
		}

	}
}
func setupVisualsText(root *root, widget *widget, skipEmpty bool) {
	if widget.textBox == nil {
		widget.textBox = &graphics.TextBox{}
	}

	var owner = root.Containers[widget.OwnerId]
	var text = root.themedField(p.Text, owner, widget)
	if skipEmpty && text == "" {
		return
	}

	var disabled = widget.isDisabled(owner)
	widget.textBox.ScaleX, widget.textBox.ScaleY = 1, 1
	widget.textBox.X, widget.textBox.Y = widget.X, widget.Y
	widget.textBox.EmbeddedColorsTag = '`'
	widget.textBox.EmbeddedAssetsTag = '^'
	widget.textBox.EmbeddedThicknessesTag = '*'
	widget.textBox.Text = text
	widget.textBox.WordWrap = defaultValue(root.themedField(p.TextWordWrap, owner, widget), "1") == "1"
	widget.textBox.PivotX, widget.textBox.PivotY = 0, 0
	widget.textBox.FontId = root.themedField(p.TextFontId, owner, widget)
	widget.textBox.LineHeight = parseNum(root.themedField(p.TextLineHeight, owner, widget), 30)
	widget.textBox.LineGap = parseNum(root.themedField(p.TextLineGap, owner, widget), 0)
	widget.textBox.SymbolGap = parseNum(root.themedField(p.TextSymbolGap, owner, widget), 0.2)
	widget.textBox.AlignmentX = parseNum(root.themedField(p.TextAlignmentX, owner, widget), 0)
	widget.textBox.AlignmentY = parseNum(root.themedField(p.TextAlignmentY, owner, widget), 0)
	widget.textBox.Width, widget.textBox.Height = widget.Width, widget.Height

	widget.textBox.EmbeddedAssetsTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedAssetsTag, owner, widget), "^")[0])
	widget.textBox.EmbeddedAssetIds = []string{
		root.themedField(p.TextEmbeddedAssetId1, owner, widget),
		root.themedField(p.TextEmbeddedAssetId2, owner, widget),
		root.themedField(p.TextEmbeddedAssetId3, owner, widget),
		root.themedField(p.TextEmbeddedAssetId4, owner, widget),
		root.themedField(p.TextEmbeddedAssetId5, owner, widget),
	}

	widget.textBox.EmbeddedColorsTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedColorsTag, owner, widget), "`")[0])
	widget.textBox.EmbeddedColors = []uint{
		parseColor(root.themedField(p.TextEmbeddedColor1, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor2, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor3, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor4, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor5, owner, widget), disabled),
	}

	widget.textBox.EmbeddedThicknessesTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedThicknessesTag, owner, widget), "*")[0])
	widget.textBox.EmbeddedThicknesses = []float32{
		parseNum(root.themedField(p.TextEmbeddedThickness1, owner, widget), 0.5),
		parseNum(root.themedField(p.TextEmbeddedThickness2, owner, widget), 0.5),
		parseNum(root.themedField(p.TextEmbeddedThickness3, owner, widget), 0.5),
		parseNum(root.themedField(p.TextEmbeddedThickness4, owner, widget), 0.5),
		parseNum(root.themedField(p.TextEmbeddedThickness5, owner, widget), 0.5),
	}
}
func drawVisuals(cam *graphics.Camera, root *root, widget *widget, fadeText bool, betweenVisualAndText func()) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = root.themedField(p.AssetId, owner, widget)
	var col = parseColor(root.themedField(p.Color, owner, widget), widget.isDisabled(owner))
	var frameCol = parseColor(root.themedField(p.FrameColor, owner, widget), widget.isDisabled(owner))
	var frameSz = parseNum(root.themedField(p.FrameSize, owner, widget), 0)

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

	if frameSz != 0 && frameCol != 0 {
		cam.DrawQuadFrame(widget.X, widget.Y, widget.Width, widget.Height, 0, frameSz, frameCol)
	}

	if widget.textBox != nil && widget.textBox.Text != "" {
		var mx, my, mw, mh = cam.MaskX, cam.MaskY, cam.MaskWidth, cam.MaskHeight
		if maskText { // used for inputbox mask
			var x, y = cam.PointToScreen(widget.X+textMargin, widget.Y+textMargin/2)
			var realX = widget.X + widget.Width - textMargin
			var realY = widget.Y + widget.Height - textMargin/2
			var xw, yh = cam.PointToScreen(realX, realY)
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
		var outlineCol = root.themedField(p.TextColorOutline, owner, widget)
		if outlineCol != "" {
			var embeddedColors = widget.textBox.EmbeddedColors
			widget.textBox.EmbeddedColors = []uint{}
			widget.textBox.Thickness = parseNum(root.themedField(p.TextThicknessOutline, owner, widget), 0.92)
			widget.textBox.Smoothness = parseNum(root.themedField(p.TextSmoothnessOutline, owner, widget), 0.08)
			widget.textBox.Tint = parseColor(outlineCol, disabled)
			cam.DrawTextBoxes(widget.textBox)
			widget.textBox.EmbeddedColors = embeddedColors
		}

		var colVal = defaultValue(root.themedField(p.TextColor, owner, widget), "127 127 127")
		var c = parseColor(colVal, disabled || fadeText)
		widget.textBox.Tint = c
		widget.textBox.Thickness = parseNum(root.themedField(p.TextThickness, owner, widget), 0.5)
		widget.textBox.Smoothness = parseNum(root.themedField(p.TextSmoothness, owner, widget), 0.02)
		cam.DrawTextBoxes(widget.textBox)

		cam.Mask(mx, my, mw, mh)
	}

	// widget.textBox.Text = "" // skip any further texts unless they are setuped beforehand
}
