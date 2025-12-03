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
	var assetId = root.themedField(p.AssetId, owner, widget)

	if assetId != "" {
		var cLeft = parseNum(dyn(owner, root.themedField(p.BoxEdgeLeft, owner, widget), "100"), 0)
		var cRight = parseNum(dyn(owner, root.themedField(p.BoxEdgeRight, owner, widget), "100"), 0)
		var cTop = parseNum(dyn(owner, root.themedField(p.BoxEdgeTop, owner, widget), "100"), 0)
		var cBottom = parseNum(dyn(owner, root.themedField(p.BoxEdgeBottom, owner, widget), "100"), 0)
		var col = parseColor(root.themedField(p.Color, owner, widget), widget.isDisabled(owner))
		var _, has = internal.Boxes[assetId]
		var offX = parseNum(dyn(owner, widget.Fields[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(owner, widget.Fields[p.OffsetY], "0"), 0)

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
	var text = root.themedField(p.Text, owner, widget)
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
	textBox.WordWrap = defaultValue(root.themedField(p.TextWordWrap, owner, widget), "on") == "on"
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.FontId = root.themedField(p.TextFontId, owner, widget)
	textBox.LineHeight = parseNum(root.themedField(p.TextLineHeight, owner, widget), 30)
	textBox.LineGap = parseNum(root.themedField(p.TextLineGap, owner, widget), 0)
	textBox.SymbolGap = parseNum(root.themedField(p.TextSymbolGap, owner, widget), 0.2)
	textBox.AlignmentX = parseNum(root.themedField(p.TextAlignmentX, owner, widget), 0)
	textBox.AlignmentY = parseNum(root.themedField(p.TextAlignmentY, owner, widget), 0)
	textBox.Width, textBox.Height = widget.Width, widget.Height

	textBox.EmbeddedAssetsTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedAssetsTag, owner, widget), "^")[0])
	textBox.EmbeddedAssetIds = []string{
		root.themedField(p.TextEmbeddedAssetId1, owner, widget),
		root.themedField(p.TextEmbeddedAssetId2, owner, widget),
		root.themedField(p.TextEmbeddedAssetId3, owner, widget),
		root.themedField(p.TextEmbeddedAssetId4, owner, widget),
		root.themedField(p.TextEmbeddedAssetId5, owner, widget),
	}

	textBox.EmbeddedColorsTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedColorsTag, owner, widget), "`")[0])
	textBox.EmbeddedColors = []uint{
		parseColor(root.themedField(p.TextEmbeddedColor1, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor2, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor3, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor4, owner, widget), disabled),
		parseColor(root.themedField(p.TextEmbeddedColor5, owner, widget), disabled),
	}

	textBox.EmbeddedThicknessesTag =
		rune(defaultValue(root.themedField(p.TextEmbeddedThicknessesTag, owner, widget), "*")[0])
	textBox.EmbeddedThicknesses = []float32{
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
		var outlineCol = root.themedField(p.TextColorOutline, owner, widget)
		if outlineCol != "" {
			var embeddedColors = textBox.EmbeddedColors
			textBox.EmbeddedColors = []uint{}
			textBox.Thickness = parseNum(root.themedField(p.TextThicknessOutline, owner, widget), 0.92)
			textBox.Smoothness = parseNum(root.themedField(p.TextSmoothnessOutline, owner, widget), 0.08)
			textBox.Color = parseColor(outlineCol, disabled)
			cam.DrawTextBoxes(&textBox)
			textBox.EmbeddedColors = embeddedColors
		}

		var colVal = defaultValue(root.themedField(p.TextColor, owner, widget), "127 127 127")
		var c = parseColor(colVal, disabled || fadeText)
		textBox.Color = c
		textBox.Thickness = parseNum(root.themedField(p.TextThickness, owner, widget), 0.5)
		textBox.Smoothness = parseNum(root.themedField(p.TextSmoothness, owner, widget), 0.02)
		cam.DrawTextBoxes(&textBox)

		cam.Mask(mx, my, mw, mh)
	}

	textBox.Text = "" // skip any further texts unless they are setuped beforehand
}
