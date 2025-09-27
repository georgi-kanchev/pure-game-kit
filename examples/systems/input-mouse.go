package example

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/mouse"
	b "pure-kit/engine/input/mouse/button"
	"pure-kit/engine/input/mouse/cursor"
	"pure-kit/engine/utility/time"
	"pure-kit/engine/window"
)

func Mouse() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	node.Width, node.Height = 300, 300

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		if mouse.IsButtonPressed(b.Left) {
			node.Angle -= time.FrameDelta() * 60
		}
		if mouse.IsButtonReleasedOnce(b.Left) {
			node.Angle = 0
		}

		if mouse.Scroll() != 0 {
			node.Width += float32(mouse.Scroll() * 20)
			node.Height = node.Width
		}

		if node.IsHovered(cam) {
			mouse.SetCursor(cursor.Hand)
		} else {
			mouse.SetCursor(cursor.Arrow)
		}

		cam.DrawNodes(&node)
	}
}
