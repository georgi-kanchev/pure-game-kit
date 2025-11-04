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
// getters

func (gui *GUI) DragOnGrab() (draggableId string) {
	return onGrab(gui.root)
}

func (gui *GUI) DragOnDrop() (grabId, dropId string) {
	return onDrop(gui.root)
}

func (gui *GUI) DragCancel() {
	if wPressedOn != nil && wPressedOn.Class == "draggable" {
		wPressedOn = nil

		var owner = gui.root.Containers[wPressedOn.OwnerId]
		sound.AssetId = defaultValue(themedProp(field.DraggableSoundCancel, gui.root, owner, wPressedOn), "~error")
		sound.Volume = gui.root.Volume
		sound.Play()
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

func drawDraggable(widget *widget, root *root, cam *graphics.Camera) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = defaultValue(themedProp(field.DraggableSpriteId, root, owner, widget), "")
	var scale = parseNum(themedProp(field.DraggableSpriteScale, root, owner, widget), 1)

	if assetId == "" {
		return
	}

	var w, h = assets.Size(assetId)
	var assetRatio = w / h
	var spriteRatio = widget.Width / widget.Height
	var drawW, drawH float32
	var disabled = widget.isDisabled(owner)
	var col = defaultValue(themedProp(field.DraggableSpriteColor, root, owner, widget), "255 255 255")

	if assetRatio > spriteRatio {
		drawW = widget.Width
		drawH = drawW / assetRatio
	} else {
		drawH = widget.Height
		drawW = drawH * assetRatio
	}

	sprite.AssetId = assetId
	sprite.X, sprite.Y = widget.DragX, widget.DragY
	sprite.Width, sprite.Height = drawW*scale, drawH*scale
	sprite.Color = parseColor(col, disabled)
	sprite.PivotX, sprite.PivotY = 0.5, 0.5
	sprite.ScaleX, sprite.ScaleY = scale, scale
	cam.DrawSprites(&sprite)
}

func onDrop(root *root) (string, string) {
	var left = mouse.IsButtonJustReleased(b.Left)
	if wPressedOn != nil && wPressedOn.Class == "draggable" && left {
		var owner = root.Containers[wPressedOn.OwnerId]
		var assetId = defaultValue(themedProp(field.DraggableSpriteId, root, owner, wPressedOn), "")
		if assetId == "" {
			return wPressedOn.Id, ""
		}

		sound.Volume = root.Volume
		defer sound.Play()
		if wFocused == nil || wFocused.Class == "draggable" {
			sound.AssetId = defaultValue(themedProp(field.ButtonSoundRelease, root, owner, wPressedOn), "~release")
			return wPressedOn.Id, wFocused.Id
		}

		sound.AssetId = defaultValue(themedProp(field.DraggableSoundCancel, root, owner, wPressedOn), "~error")
		return wPressedOn.Id, ""
	}
	return "", ""
}
func onGrab(root *root) string {
	var result = condition.JustTurnedTrue(wPressedOn != nil && wPressedOn.Class == "draggable", ";;;;draggg-start")
	if result {
		var owner = root.Containers[wPressedOn.OwnerId]
		if themedProp(field.DraggableSpriteId, root, owner, wPressedOn) == "" {
			return ""
		}
		return wPressedOn.Id
	}
	return ""
}
