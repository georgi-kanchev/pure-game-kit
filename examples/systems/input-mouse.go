package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
	"pure-kit/engine/window"
)

func Mouse() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	node.Width, node.Height = 300, 300

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		if mouse.IsButtonPressed(mouse.ButtonLeft) {
			node.Angle -= seconds.FrameDelta() * 60
		}
		if mouse.IsButtonReleasedOnce(mouse.ButtonLeft) {
			node.Angle = 0
		}

		if mouse.Scroll() != 0 {
			node.Width += float32(mouse.Scroll() * 20)
			node.Height = node.Width
		}

		if node.IsHovered(cam) {
			mouse.SetCursor(mouse.CursorHand)
		} else {
			mouse.SetCursor(mouse.CursorArrow)
		}

		cam.DrawNodes(&node)
	}
}
