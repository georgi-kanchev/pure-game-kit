package gui

import (
	"pure-game-kit/graphics"
	f "pure-game-kit/gui/field"
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
	var assetId = root.themedField(f.AssetId, owner, widget)

	if widget.sprite == nil {
		widget.sprite = graphics.NewSprite("", 0, 0)
		widget.top = graphics.NewSprite("", 0, 0)
		widget.left = graphics.NewSprite("", 0, 0)
		widget.right = graphics.NewSprite("", 0, 0)
		widget.bottom = graphics.NewSprite("", 0, 0)
		widget.sprite.PivotX, widget.sprite.PivotY = 0, 0
	}
	if widget.box == nil {
		widget.box = graphics.NewBox("", 0, 0)
		widget.box.PivotX, widget.box.PivotY = 0, 0
	}

	var col = parseColor(root.themedField(f.Color, owner, widget), widget.isDisabled(owner))
	var _, has = internal.Boxes[assetId]
	var sprite, box = widget.sprite, widget.box

	if has {
		var cLeft = parseNum(root.themedField(f.BoxEdgeLeft, owner, widget), 0)
		var cRight = parseNum(root.themedField(f.BoxEdgeRight, owner, widget), 0)
		var cTop = parseNum(root.themedField(f.BoxEdgeTop, owner, widget), 0)
		var cBottom = parseNum(root.themedField(f.BoxEdgeBottom, owner, widget), 0)

		box.X, box.Y = widget.X, widget.Y
		box.AssetId = assetId
		box.Tint = col
		box.Width, box.Height = widget.Width, widget.Height
		box.EdgeLeft, box.EdgeRight = cLeft, cRight
		box.EdgeTop, box.EdgeBottom = cTop, cBottom
	} else {
		sprite.X, sprite.Y = widget.X, widget.Y
		sprite.AssetId = assetId
		sprite.Tint = col
		sprite.TextureRepeat = false
		sprite.Width, sprite.Height = widget.Width, widget.Height
	}
}
func setupVisualsText(root *root, widget *widget, skipEmpty bool) {
	if widget.textBox == nil {
		widget.textBox = &graphics.TextBox{}
	}

	var owner = root.Containers[widget.OwnerId]
	var text = root.themedField(f.Text, owner, widget)
	if skipEmpty && text == "" {
		return
	}

	widget.textBox.ScaleX, widget.textBox.ScaleY = 1, 1
	widget.textBox.X, widget.textBox.Y = widget.X, widget.Y
	widget.textBox.Text = text
	widget.textBox.WordWrap = defaultValue(root.themedField(f.TextWordWrap, owner, widget), "1") == "1"
	widget.textBox.PivotX, widget.textBox.PivotY = 0, 0
	widget.textBox.FontId = root.themedField(f.TextFontId, owner, widget)
	widget.textBox.LineHeight = parseNum(root.themedField(f.TextLineHeight, owner, widget), 30)
	widget.textBox.LineGap = parseNum(root.themedField(f.TextLineGap, owner, widget), 0)
	widget.textBox.SymbolGap = parseNum(root.themedField(f.TextSymbolGap, owner, widget), 0.2)
	widget.textBox.AlignmentX = parseNum(root.themedField(f.TextAlignmentX, owner, widget), 0)
	widget.textBox.AlignmentY = parseNum(root.themedField(f.TextAlignmentY, owner, widget), 0)
	widget.textBox.Width, widget.textBox.Height = widget.Width, widget.Height
	widget.textBox.Fast = root.themedField(f.TextFast, owner, widget) != ""
}
func drawVisuals(cam *graphics.Camera, root *root, widget *widget, fadeText bool, betweenVisualAndText func()) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = root.themedField(f.AssetId, owner, widget)
	var frameCol = parseColor(root.themedField(f.FrameColor, owner, widget), widget.isDisabled(owner))
	var frameSz = parseNum(root.themedField(f.FrameSize, owner, widget), 0)

	var _, has = internal.Boxes[assetId]
	if has {
		cam.DrawBoxes(widget.box)
	} else {
		cam.DrawSprites(widget.sprite)
	}

	if frameSz != 0 && frameCol != 0 {
		widget.top.PivotX, widget.top.PivotY = 0, 0
		widget.left.PivotX, widget.left.PivotY = 0, 0
		widget.right.PivotX, widget.right.PivotY = 0, 0
		widget.bottom.PivotX, widget.bottom.PivotY = 0, 0
		widget.top.Tint = frameCol
		widget.left.Tint = frameCol
		widget.right.Tint = frameCol
		widget.bottom.Tint = frameCol

		if frameSz < 0 {
			var t = -frameSz
			widget.top.X, widget.top.Y = widget.X, widget.Y
			widget.top.Width, widget.top.Height = widget.Width, t
			widget.right.X, widget.right.Y = widget.X+widget.Width-t, widget.Y
			widget.right.Width, widget.right.Height = t, widget.Height
			widget.bottom.X, widget.bottom.Y = widget.X, widget.Y+widget.Height-t
			widget.bottom.Width, widget.bottom.Height = widget.Width, t
			widget.left.X, widget.left.Y = widget.X, widget.Y
			widget.left.Width, widget.left.Height = t, widget.Height

		} else {
			var t = frameSz
			widget.top.X, widget.top.Y = widget.X-t, widget.Y-t
			widget.top.Width, widget.top.Height = widget.Width+(t*2), t
			widget.right.X, widget.right.Y = widget.X+widget.Width, widget.Y-t
			widget.right.Width, widget.right.Height = t, widget.Height+(t*2)
			widget.bottom.X, widget.bottom.Y = widget.X-t, widget.Y+widget.Height
			widget.bottom.Width, widget.bottom.Height = widget.Width+(t*2), t
			widget.left.X, widget.left.Y = widget.X-t, widget.Y-t
			widget.left.Width, widget.left.Height = t, widget.Height+(t*2)
		}
		cam.DrawSprites(widget.top, widget.left, widget.right, widget.bottom)
	}

	if widget.textBox == nil || widget.textBox.Text == "" {
		return
	}

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
		cam.Mask(x, y, xw-x+1, yh-y)
	}

	if betweenVisualAndText != nil {
		betweenVisualAndText()
	}

	var disabled = widget.isDisabled(owner)
	var colVal = defaultValue(root.themedField(f.TextColor, owner, widget), "127 127 127")
	var c = parseColor(colVal, disabled || fadeText)
	widget.textBox.Tint = c
	cam.DrawTextBoxes(widget.textBox)
	cam.Mask(mx, my, mw, mh)
}

func drawText(cam *graphics.Camera, root *root, widget *widget, fadeText bool, betweenVisualAndText func()) {

}
