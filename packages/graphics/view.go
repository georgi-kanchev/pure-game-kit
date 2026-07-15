// Strictly tied to the window, drawing on it through a view and converting between the two coordinate systems.
// The view's drawing consists of two categories: primitives and objects. While using the assets for drawing,
// the graphical objects are still very lightweight and exist independently of them.
package graphics

import (
	geometry "pure-game-kit/packages/geometry"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/number"
)

type View struct {
	X, Y, Zoom, Angle float32
	WindowArea        geometry.Area // The drawing area in window space. Zero value = entire window.

	// =================================================================

	velocityX, velocityY,
	dragVelX, dragVelY float32
	debugBuffer []byte
}

func NewView(zoom float32) View { return View{Zoom: zoom} }

// =================================================================

func (v *View) MouseDragAndZoom() {
	var oldZoom, scroll = v.Zoom, mouse.ScrollY()

	if scroll != 0 {
		v.Zoom *= 1 + 0.05*scroll

		var mx, my = v.MousePosition()
		v.X += (mx - v.X) * (v.Zoom/oldZoom - 1)
		v.Y += (my - v.Y) * (v.Zoom/oldZoom - 1)
	}

	if mouse.IsButtonPressed(button.Middle) {
		var dx, dy = mouse.CursorDelta()
		var sin, cos = internal.SinCos(-v.Angle)
		v.X -= (dx*cos - dy*sin) / v.Zoom
		v.Y -= (dx*sin + dy*cos) / v.Zoom
	}
}
func (v *View) MouseDragAndZoomSmoothly() {
	var oldZoom, scroll = v.Zoom, mouse.ScrollSmoothY()

	if scroll != 0 {
		v.Zoom *= 1 + 0.001*scroll

		var mx, my = v.MousePosition()
		v.X += (mx - v.X) * (v.Zoom/oldZoom - 1)
		v.Y += (my - v.Y) * (v.Zoom/oldZoom - 1)
	}

	const dragFriction, dragStrength = 8.0, 8.0
	var dt = internal.FrameDelta
	var decay = number.Exponential(-dragFriction * dt)
	v.velocityX, v.velocityY = v.velocityX*decay, v.velocityY*decay

	if mouse.IsButtonPressed(button.Middle) {
		var sin, cos = internal.SinCos(-v.Angle)
		var dx, dy = mouse.CursorDelta()

		dx, dy = dx/dt, dy/dt
		var wdx, wdy = (dx*cos - dy*sin) / v.Zoom, (dx*sin + dy*cos) / v.Zoom
		v.velocityX, v.velocityY = v.velocityX-(wdx*dragStrength*dt), v.velocityY-(wdy*dragStrength*dt)
	}

	v.X, v.Y = v.X+(v.velocityX*dt), v.Y+v.velocityY*dt

	if number.Absolute(v.velocityX) < 0.0001 {
		v.velocityX = 0
	}
	if number.Absolute(v.velocityY) < 0.0001 {
		v.velocityY = 0
	}
}

func (v *View) FitSize(width, height float32) {
	var windowArea = v.windowArea()
	if windowArea.Width <= 0 || windowArea.Height <= 0 {
		return
	}

	var w, h = width, height
	var minX, maxX, minY, maxY float32 = 0, 0, 0, 0
	var corners = [4][2]float32{{-w / 2, -h / 2}, {w / 2, -h / 2}, {-w / 2, h / 2}, {w / 2, h / 2}}
	var sin, cos = internal.SinCos(-v.Angle)

	for i, p := range corners {
		var dx, dy = p[0], p[1]
		var localX, localY = dx*cos - dy*sin, dx*sin + dy*cos

		if i == 0 {
			minX, maxX, minY, maxY = localX, localX, localY, localY
		} else {
			minX, maxX, minY, maxY = min(minX, localX), max(maxX, localX), min(minY, localY), max(maxY, localY)
		}
	}

	var localWidth, localHeight = max(maxX-minX, 1), max(maxY-minY, 1) // no division by zero for tiny/empty bounds
	v.Zoom = min(windowArea.Width/localWidth, windowArea.Height/localHeight)
}

// =================================================================

func (v *View) IsAreaVisible(x, y, width, height float32) bool {
	var sx1, sy1 = v.PointToScreen(x, y)
	var sx2, sy2 = v.PointToScreen(x+width, y)
	var sx3, sy3 = v.PointToScreen(x, y+height)
	var sx4, sy4 = v.PointToScreen(x+width, y+height)
	var minX, maxX = min(min(sx1, sx2), min(sx3, sx4)), max(max(sx1, sx2), max(sx3, sx4))
	var minY, maxY = min(min(sy1, sy2), min(sy3, sy4)), max(max(sy1, sy2), max(sy3, sy4))
	return maxX >= 0 && minX <= internal.WindowWidth && maxY >= 0 && minY <= internal.WindowHeight
}
func (v *View) IsHovered() bool {
	return internal.MouseX > 0 && internal.MouseX < internal.WindowWidth &&
		internal.MouseY > 0 && internal.MouseY < internal.WindowHeight
}
func (v *View) MousePosition() (x, y float32) {
	return v.PointFromScreen(internal.MouseX, internal.MouseY)
}
func (v *View) MouseDelta() (deltaX, deltaY float32) {
	var mx, my = mouse.CursorDelta()
	return mx * v.Zoom, my * v.Zoom
}
func (v *View) Size() (width, height float32) {
	var windowArea = v.windowArea()
	return windowArea.Width / v.Zoom, windowArea.Height / v.Zoom
}
func (v *View) Bounds() (x, y, width, height float32) {
	var x1, y1 = v.PointFromEdge(0, 0)
	var x2, y2 = v.PointFromEdge(1, 0)
	var x3, y3 = v.PointFromEdge(1, 1)
	var x4, y4 = v.PointFromEdge(0, 1)
	var minX, minY = number.Minimum(x1, x2, x3, x4), number.Minimum(y1, y2, y3, y4)
	var maxX, maxY = number.Maximum(x1, x2, x3, x4), number.Maximum(y1, y2, y3, y4)
	return minX, minY, maxX - minX, maxY - minY
}

func (v *View) PointFromScreen(screenX, screenY float32) (x, y float32) {
	var wa = v.windowArea()
	x, y = screenX-wa.X, screenY-wa.Y
	if v.Zoom != 0 {
		x, y = x/v.Zoom, y/v.Zoom
	}
	if v.Angle != 0 {
		var sin, cos = internal.SinCos(v.Angle)
		x, y = x*cos-y*sin, x*sin+y*cos
	}
	return x + v.X, y + v.Y
}
func (v *View) PointToScreen(x, y float32) (screenX, screenY float32) {
	x, y = x-v.X, y-v.Y
	if v.Angle != 0 {
		var sin, cos = internal.SinCos(v.Angle)
		x, y = x*cos+y*sin, -x*sin+y*cos
	}
	if v.Zoom != 0 {
		x, y = x*v.Zoom, y*v.Zoom
	}
	var wa = v.windowArea()
	return x + wa.X, y + wa.Y
}
func (v *View) PointFromView(otherView *View, otherX, otherY float32) (myX, myY float32) {
	return v.PointFromScreen(otherView.PointToScreen(otherX, otherY))
}
func (v *View) PointToView(otherView *View, myX, myY float32) (otherX, otherY float32) {
	return otherView.PointFromView(v, myX, myY)
}
func (v *View) PointFromEdge(horizontal, vertical float32) (x, y float32) {
	var wa = v.windowArea()
	return v.PointFromScreen(wa.X-wa.Width/2+wa.Width*horizontal, wa.Y-wa.Height/2+wa.Height*vertical)
}
