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

func (g *GUI) DragOnGrab() (draggableId string) {
	return onGrab(g.root)
}

func (g *GUI) DragOnDrop() (grabId, dropId string) {
	return onDrop(g.root)
}

func (g *GUI) DragCancel() {
	if g.root.wPressedOn != nil && g.root.wPressedOn.Class == "draggable" {
		g.root.wPressedOn = nil

		var owner = g.root.Containers[g.root.wPressedOn.OwnerId]
		var val = g.root.themedField(field.DraggableSoundCancel, owner, g.root.wPressedOn)
		sound.AssetId = defaultValue(val, "~error")
		sound.Volume = g.root.Volume
		sound.Play()
	}
}

//=================================================================
// private

func draggable(cam *graphics.Camera, root *root, widget *widget) {
	if root.wPressedOn == widget {
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
	var assetId = defaultValue(root.themedField(field.DraggableSpriteId, owner, widget), "")
	var scale = parseNum(root.themedField(field.DraggableSpriteScale, owner, widget), 1)

	if assetId == "" {
		return
	}

	var w, h = assets.Size(assetId)
	var assetRatio = float32(w) / float32(h)
	var spriteRatio = widget.Width / widget.Height
	var drawW, drawH float32
	var disabled = widget.isDisabled(owner)
	var col = defaultValue(root.themedField(field.DraggableSpriteColor, owner, widget), "255 255 255")

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
	sprite.Tint = parseColor(col, disabled)
	sprite.PivotX, sprite.PivotY = 0.5, 0.5
	sprite.ScaleX, sprite.ScaleY = scale, scale
	cam.DrawSprites(&sprite)
}

func onDrop(root *root) (string, string) {
	var left = mouse.IsButtonJustReleased(b.Left)
	if root.wPressedOn != nil && root.wPressedOn.Class == "draggable" && left {
		var owner = root.Containers[root.wPressedOn.OwnerId]
		var assetId = defaultValue(root.themedField(field.DraggableSpriteId, owner, root.wPressedOn), "")
		if assetId == "" {
			return root.wPressedOn.Id, ""
		}

		sound.Volume = root.Volume
		defer sound.Play()
		if root.wFocused == nil || root.wFocused.Class == "draggable" {
			var val = root.themedField(field.ButtonSoundRelease, owner, root.wPressedOn)
			sound.AssetId = defaultValue(val, "~release")
			return root.wPressedOn.Id, root.wFocused.Id
		}

		sound.AssetId = defaultValue(root.themedField(field.DraggableSoundCancel, owner, root.wPressedOn), "~error")
		return root.wPressedOn.Id, ""
	}
	return "", ""
}
func onGrab(root *root) string {
	var cond = root.wPressedOn != nil && root.wPressedOn.Class == "draggable"
	var result = condition.JustTurnedTrue(cond, ";;;;draggg-start")
	if result {
		var owner = root.Containers[root.wPressedOn.OwnerId]
		if root.themedField(field.DraggableSpriteId, owner, root.wPressedOn) == "" {
			return ""
		}
		return root.wPressedOn.Id
	}
	return ""
}
