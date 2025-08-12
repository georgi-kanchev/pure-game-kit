package gui

import "pure-kit/engine/graphics"

func buttonUpdateAndDraw(cam *graphics.Camera, widget *widget, owner *container) {
	widget.draw(cam, owner)
}
