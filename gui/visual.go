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

func setupVisualsTextured(w *widget) {
	var owner = w.root.Containers[w.OwnerId]
	var assetId = w.root.themedField(f.AssetId, owner, w)

	if w.sprite == nil {
		w.sprite = graphics.NewSprite("", 0, 0)
		w.top = graphics.NewSprite("", 0, 0)
		w.left = graphics.NewSprite("", 0, 0)
		w.right = graphics.NewSprite("", 0, 0)
		w.bottom = graphics.NewSprite("", 0, 0)
		w.sprite.PivotX, w.sprite.PivotY = 0, 0
	}
	if w.box == nil {
		w.box = graphics.NewBox("", 0, 0)
		w.box.PivotX, w.box.PivotY = 0, 0
	}

	var col = parseColor(w.root.themedField(f.Color, owner, w), w.isDisabled(owner))
	var _, has = internal.Boxes[assetId]
	var sprite, box = w.sprite, w.box

	if has {
		var cLeft = parseNum(w.root.themedField(f.BoxEdgeLeft, owner, w), 0)
		var cRight = parseNum(w.root.themedField(f.BoxEdgeRight, owner, w), 0)
		var cTop = parseNum(w.root.themedField(f.BoxEdgeTop, owner, w), 0)
		var cBottom = parseNum(w.root.themedField(f.BoxEdgeBottom, owner, w), 0)

		box.X, box.Y = w.X, w.Y
		box.AssetId = assetId
		box.Tint = col
		box.Width, box.Height = w.Width, w.Height
		box.EdgeLeft, box.EdgeRight = cLeft, cRight
		box.EdgeTop, box.EdgeBottom = cTop, cBottom
	} else {
		sprite.X, sprite.Y = w.X, w.Y
		sprite.AssetId = assetId
		sprite.Tint = col
		sprite.TextureRepeat = false
		sprite.Width, sprite.Height = w.Width, w.Height
	}
}
func setupVisualsText(w *widget, skipEmpty bool) {
	if w.textBox == nil {
		w.textBox = &graphics.TextBox{}
	}

	var owner = w.root.Containers[w.OwnerId]
	var text = w.root.themedField(f.Text, owner, w)
	if skipEmpty && text == "" {
		return
	}

	w.textBox.ScaleX, w.textBox.ScaleY = 1, 1
	w.textBox.X, w.textBox.Y = w.X, w.Y
	w.textBox.Text = text
	w.textBox.WordWrap = defaultValue(w.root.themedField(f.TextWordWrap, owner, w), "1") == "1"
	w.textBox.PivotX, w.textBox.PivotY = 0, 0
	w.textBox.FontId = w.root.themedField(f.TextFontId, owner, w)
	w.textBox.LineHeight = parseNum(w.root.themedField(f.TextLineHeight, owner, w), 30)
	w.textBox.LineGap = parseNum(w.root.themedField(f.TextLineGap, owner, w), 0)
	w.textBox.SymbolGap = parseNum(w.root.themedField(f.TextSymbolGap, owner, w), 0.2)
	w.textBox.AlignmentX = parseNum(w.root.themedField(f.TextAlignmentX, owner, w), 0)
	w.textBox.AlignmentY = parseNum(w.root.themedField(f.TextAlignmentY, owner, w), 0)
	w.textBox.Width, w.textBox.Height = w.Width, w.Height
	w.textBox.Fast = w.root.themedField(f.TextFast, owner, w) != ""
}
func drawVisuals(w *widget, fadeText bool, betweenVisualAndText func()) {
	var cam = w.root.cam
	var owner = w.root.Containers[w.OwnerId]
	var assetId = w.root.themedField(f.AssetId, owner, w)
	var frameCol = parseColor(w.root.themedField(f.FrameColor, owner, w), w.isDisabled(owner))
	var frameSz = parseNum(w.root.themedField(f.FrameSize, owner, w), 0)

	var _, has = internal.Boxes[assetId]
	if has {
		cam.DrawBoxes(w.box)
	} else {
		cam.DrawSprites(w.sprite)
	}

	if frameSz != 0 && frameCol != 0 {
		w.top.PivotX, w.top.PivotY = 0, 0
		w.left.PivotX, w.left.PivotY = 0, 0
		w.right.PivotX, w.right.PivotY = 0, 0
		w.bottom.PivotX, w.bottom.PivotY = 0, 0
		w.top.Tint = frameCol
		w.left.Tint = frameCol
		w.right.Tint = frameCol
		w.bottom.Tint = frameCol

		if frameSz < 0 {
			var t = -frameSz
			w.top.X, w.top.Y = w.X, w.Y
			w.top.Width, w.top.Height = w.Width, t
			w.right.X, w.right.Y = w.X+w.Width-t, w.Y
			w.right.Width, w.right.Height = t, w.Height
			w.bottom.X, w.bottom.Y = w.X, w.Y+w.Height-t
			w.bottom.Width, w.bottom.Height = w.Width, t
			w.left.X, w.left.Y = w.X, w.Y
			w.left.Width, w.left.Height = t, w.Height

		} else {
			var t = frameSz
			w.top.X, w.top.Y = w.X-t, w.Y-t
			w.top.Width, w.top.Height = w.Width+(t*2), t
			w.right.X, w.right.Y = w.X+w.Width, w.Y-t
			w.right.Width, w.right.Height = t, w.Height+(t*2)
			w.bottom.X, w.bottom.Y = w.X-t, w.Y+w.Height
			w.bottom.Width, w.bottom.Height = w.Width+(t*2), t
			w.left.X, w.left.Y = w.X-t, w.Y-t
			w.left.Width, w.left.Height = t, w.Height+(t*2)
		}
		cam.DrawSprites(w.top, w.left, w.right, w.bottom)
	}

	if w.textBox == nil || w.textBox.Text == "" {
		return
	}

	var mx, my, mw, mh = cam.MaskX, cam.MaskY, cam.MaskWidth, cam.MaskHeight
	if maskText { // used for inputbox mask
		var x, y = cam.PointToScreen(w.X+textMargin, w.Y+textMargin/2)
		var realX = w.X + w.Width - textMargin
		var realY = w.Y + w.Height - textMargin/2
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

	var disabled = w.isDisabled(owner)
	var colVal = defaultValue(w.root.themedField(f.TextColor, owner, w), "127 127 127")
	var c = parseColor(colVal, disabled || fadeText)
	w.textBox.Tint = c
	cam.DrawTextBoxes(w.textBox)
	cam.Mask(mx, my, mw, mh)
}
