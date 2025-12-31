package example

import (
	"pure-game-kit/graphics"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Mouse() {
	var cam = graphics.NewCamera(1)
	var node = graphics.NewNode(0, 0)
	node.Width, node.Height = 300, 300

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		if node.IsHovered(cam) && mouse.IsButtonPressed(b.Left) {
			node.Angle -= time.FrameDelta() * 60
		}
		if mouse.IsButtonJustReleased(b.Left) {
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

		cam.DrawNodes(node)
	}
}
