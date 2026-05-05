package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	b "pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/input/mouse/cursor"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func Mouse() {
	var cam = graphics.NewCamera(1)
	var quad = graphics.NewQuad(0, 0)
	quad.Width, quad.Height = 300, 300

	for window.KeepOpen() {
		var mx, my = cam.MousePosition()
		if quad.ContainsPoint(mx, my) && mouse.IsButtonPressed(b.Left) {
			quad.Angle -= time.FrameDelta() * 60
		}
		if mouse.IsButtonJustReleased(b.Left) {
			quad.Angle = 0
		}

		if mouse.Scroll() != 0 {
			quad.Width += float32(mouse.Scroll() * 20)
			quad.Height = quad.Width
		}

		if quad.ContainsPoint(mx, my) {
			mouse.SetCursor(cursor.Hand)
		} else {
			mouse.SetCursor(cursor.Arrow)
		}

		cam.DrawQuads(quad)
	}
}
