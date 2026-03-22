package gui

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
)

func Draggable(id string, properties ...string) string {
	return newWidget("draggable", id, properties...)
}

//=================================================================

func (g *GUI) DragJustGrabbed() (id string) {
	return g.root.onGrab()
}
func (g *GUI) DragJustDropped() (grabId, dropId string) {
	return g.root.onDrop()
}
func (g *GUI) DragCancel() {
	if g.root.wPressedOn != nil && g.root.wPressedOn.Class == "draggable" {
		g.root.wPressedOn = nil

		var owner = g.root.Containers[g.root.wPressedOn.OwnerId]
		var val = g.root.themedField(field.DraggableCancelSoundId, owner, g.root.wPressedOn)
		sound.AssetId = defaultValue(val, "~error")
		sound.Volume = g.root.Volume
		sound.Play()
	}
}

//=================================================================
// private

func draggable(w *widget) {
	if w.root.wPressedOn == w {
		mouse.SetCursor(cursor.Hand)
		w.DragX += (mouseX - prevMouseX) / w.root.cam.Zoom
		w.DragY += (mouseY - prevMouseY) / w.root.cam.Zoom
	} else {
		w.DragX, w.DragY = w.X+w.Width/2, w.Y+w.Height/2
	}

	button(w)
}

func drawDraggable(widget *widget) {
	var owner = widget.root.Containers[widget.OwnerId]
	var assetId = defaultValue(widget.root.themedField(field.DraggableAssetId, owner, widget), "")
	var scale = parseNum(widget.root.themedField(field.DraggableAssetScale, owner, widget), 1)

	if assetId == "" {
		return
	}

	var w, h = assets.Size(assetId)
	var assetRatio = float32(w) / float32(h)
	var spriteRatio = widget.Width / widget.Height
	var drawW, drawH float32
	var disabled = widget.isDisabled(owner)
	var col = defaultValue(widget.root.themedField(field.DraggableAssetColor, owner, widget), "255 255 255")
	var sprite = graphics.NewSprite(assetId, widget.DragX, widget.DragY)

	if assetRatio > spriteRatio {
		drawW = widget.Width
		drawH = drawW / assetRatio
	} else {
		drawH = widget.Height
		drawW = drawH * assetRatio
	}

	sprite.Width, sprite.Height = drawW*scale, drawH*scale
	sprite.Tint = parseColor(col, disabled)
	sprite.PivotX, sprite.PivotY = 0.5, 0.5
	sprite.ScaleX, sprite.ScaleY = scale, scale
	widget.root.sprites = append(widget.root.sprites, sprite)
}

func (r *root) onDrop() (string, string) {
	var left = mouse.IsButtonJustReleased(b.Left)
	if r.wPressedOn != nil && r.wPressedOn.Class == "draggable" && left {
		var owner = r.Containers[r.wPressedOn.OwnerId]
		var assetId = defaultValue(r.themedField(field.DraggableAssetId, owner, r.wPressedOn), "")
		if assetId == "" {
			return r.wPressedOn.Id, ""
		}

		sound.Volume = r.Volume
		defer sound.Play()
		if r.wFocused == nil || r.wFocused.Class == "draggable" {
			var val = r.themedField(field.ButtonSoundRelease, owner, r.wPressedOn)
			sound.AssetId = defaultValue(val, "~release")
			return r.wPressedOn.Id, r.wFocused.Id
		}

		sound.AssetId = defaultValue(r.themedField(field.DraggableCancelSoundId, owner, r.wPressedOn), "~error")
		return r.wPressedOn.Id, ""
	}
	return "", ""
}
func (r *root) onGrab() string {
	var cond = r.wPressedOn != nil && r.wPressedOn.Class == "draggable"
	var result = condition.JustTurnedTrue(cond, ";;;;draggg-start")
	if result {
		var owner = r.Containers[r.wPressedOn.OwnerId]
		if r.themedField(field.DraggableAssetId, owner, r.wPressedOn) == "" {
			return ""
		}
		return r.wPressedOn.Id
	}
	return ""
}
