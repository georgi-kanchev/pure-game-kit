package gui

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	b "pure-kit/engine/input/mouse/button"
	"pure-kit/engine/input/mouse/cursor"
)

func Draggable(id string, properties ...string) string {
	return newWidget("draggable", id, properties...)
}

//=================================================================
// getters

func (gui *GUI) DragOnGrab() (draggableId string) {
	var result = condition.TrueOnce(wPressedOn != nil && wPressedOn.Class == "draggable", ";;;;draggg-start")
	if result {
		var owner = gui.root.Containers[wPressedOn.OwnerId]
		if themedProp(property.DraggableSpriteId, gui.root, owner, wPressedOn) == "" {
			return ""
		}
		return wPressedOn.Id
	}
	return ""

}
func (gui *GUI) DragOnDrop() (grabId, dropId string) {
	var left = mouse.IsButtonReleasedOnce(b.Left)
	if wPressedOn != nil && wPressedOn.Class == "draggable" && left {
		if wFocused == nil || wFocused.Class == "draggable" {
			return wPressedOn.Id, wFocused.Id
		}

		return wPressedOn.Id, ""
	}
	return "", ""
}
func (gui *GUI) DragCancel() {
	if wPressedOn != nil && wPressedOn.Class == "draggable" {
		wPressedOn = nil
	}
}

//=================================================================
// private

func draggable(cam *graphics.Camera, root *root, widget *widget) {
	if wPressedOn == widget {
		mouse.SetCursor(cursor.Hand)
		widget.DragX += mouseX - prevMouseX
		widget.DragY += mouseY - prevMouseY
	} else {
		widget.DragX, widget.DragY = widget.X+widget.Width/2, widget.Y+widget.Height/2
	}

	button(cam, root, widget)
}

func drawDraggable(draggable *widget, root *root, cam *graphics.Camera) {
	var owner = root.Containers[draggable.OwnerId]
	var assetId = defaultValue(themedProp(property.DraggableSpriteId, root, owner, draggable), ";;;;;")
	var scale = parseNum(themedProp(property.DraggableSpriteScale, root, owner, draggable), 1)

	if assetId == "" {
		return
	}

	var w, h = assets.Size(assetId)
	var assetRatio = w / h
	var spriteRatio = draggable.Width / draggable.Height
	var drawW, drawH float32
	var disabled = draggable.isDisabled(owner)
	var col = defaultValue(themedProp(property.DraggableSpriteColor, root, owner, draggable), "255 255 255")

	if assetRatio > spriteRatio {
		drawW = draggable.Width
		drawH = drawW / assetRatio
	} else {
		drawH = draggable.Height
		drawW = drawH * assetRatio
	}

	reusableSprite.AssetId = assetId
	reusableSprite.X, reusableSprite.Y = draggable.DragX, draggable.DragY
	reusableSprite.Width, reusableSprite.Height = drawW*scale, drawH*scale
	reusableSprite.Color = parseColor(col, disabled)
	reusableSprite.PivotX, reusableSprite.PivotY = 0.5, 0.5
	reusableSprite.ScaleX, reusableSprite.ScaleY = scale, scale
	cam.DrawSprites(&reusableSprite)
}
