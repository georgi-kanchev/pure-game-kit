package gui

import (
	"encoding/xml"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	f "pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []*widget  `xml:"Widget"`
	XmlThemes  []*theme   `xml:"Theme"`

	X, Y, Width, Height,
	ScrollX, ScrollY float32

	prevMouseX, prevMouseY,
	velocityX, velocityY,
	targetScrollX, targetScrollY float32

	Properties map[string]string
	Widgets    []string
	WasHidden  bool
}

func Container(id, x, y, width, height string, properties ...string) string {
	var rid = "<Container " + f.Id + "=\"" + id + "\""
	var rx = condition.If(x == "", "", " "+f.X+"=\""+x+"\"")
	var ry = condition.If(y == "", "", " "+f.Y+"=\""+y+"\"")
	var rw = condition.If(width == "", "", " "+f.Width+"=\""+width+"\"")
	var rh = condition.If(height == "", "", " "+f.Height+"=\""+height+"\"")
	return rid + rx + ry + rw + rh + extraProps(properties...) + ">"
}

//=================================================================
// private

const scrollSize, handleSpeed, dragFriction, dragMomentum = 25.0, 12.0, 0.95, 30.0

var cMiddlePressed, cPressedOnScrollH, cPressedOnScrollV *container

func (c *container) updateAndDraw(root *root, cam *graphics.Camera) {
	var hidden, _ = c.Properties[f.Hidden]
	if hidden != "" {
		return
	}

	var x, y, w, h = parseNum(ownerX, 0), parseNum(ownerY, 0), parseNum(ownerW, 0), parseNum(ownerH, 0)
	var scx, scy = cam.PointToScreen(float32(x), float32(y))
	var cGapX = parseNum(dyn(c, c.Properties[f.GapX], "0"), 0)
	var cGapY = parseNum(dyn(c, c.Properties[f.GapY], "0"), 0)
	var curX, curY = x + cGapX, y + cGapY
	var maxHeight float32
	var maskW, maskH = (w - cGapX*2) * cam.Zoom, (h - cGapY*2) * cam.Zoom
	var nonBgrIndex = 0
	var draggables []*widget = make([]*widget, 0)

	cam.Mask(scx+int(cGapX*cam.Zoom), scy+int(cGapY*cam.Zoom), int(maskW), int(maskH))
	c.X, c.Y, c.Width, c.Height = x, y, w, h

	if c.isHovered(cam) {
		cHovered = c
	}
	if c.isFocused(cam) && mouse.IsButtonJustPressed(b.Middle) {
		cMiddlePressed = c
	}

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var wHidden, _ = widget.Properties[f.Hidden]
		if wHidden != "" || widget.Class == "tooltip" {
			continue
		}

		var _, isBgr = widget.Properties[f.FillContainer]
		var ww = parseNum(dyn(c, root.themedField(f.Width, c, widget), "0"), 0)
		var wh = parseNum(dyn(c, root.themedField(f.Height, c, widget), "0"), 0)
		var gapX = parseNum(dyn(c, root.themedField(f.GapX, c, widget), "0"), 0)
		var gapY = parseNum(dyn(c, root.themedField(f.GapY, c, widget), "0"), 0)
		var offX = parseNum(dyn(c, widget.Properties[f.OffsetX], "0"), 0)
		var offY = parseNum(dyn(c, widget.Properties[f.OffsetY], "0"), 0)

		if isBgr {
			widget.X, widget.Y = x, y
			ww, wh = w, h
			cam.Mask(cam.ScreenX, cam.ScreenY, cam.ScreenWidth, cam.ScreenHeight) // mask doesn't affect bgr
		} else {
			var row, newRow = widget.Properties[f.NewRow]
			if newRow {
				curX = x + cGapX
				curY += parseNum(dyn(c, row, text.New(maxHeight+gapY)), 0)
				maxHeight = 0
			}

			curX += condition.If(newRow || nonBgrIndex == 0, 0, gapX)
			widget.X = curX + offX
			widget.Y = curY + offY
			curX += ww
			maxHeight = condition.If(maxHeight < wh, wh, maxHeight)
			nonBgrIndex++
		}

		widget.Width, widget.Height = ww, wh
		widget.ThemeId = root.themedField(f.ThemeId, c, widget)

		if !isBgr {
			widget.X -= c.ScrollX
			widget.Y -= c.ScrollY
		}

		if widget.isHovered(c, cam) {
			wHovered = widget
		}

		var outsideX = widget.X+widget.Width < c.X || widget.X > c.X+c.Width
		var outsideY = widget.Y+widget.Height < c.Y || widget.Y > c.Y+c.Height
		widget.IsCulled = outsideX || outsideY

		if !widget.IsCulled { // culling widgets outside of the container (masked, invisible)
			if widget.UpdateAndDraw != nil {
				widget.UpdateAndDraw(cam, root, widget)
				tryShowTooltip(widget, root, c, cam)
			} else if widget.Class == "visual" {
				setupVisualsTextured(root, widget)
				setupVisualsText(root, widget, true)
				drawVisuals(cam, root, widget, false, nil)
				tryShowTooltip(widget, root, c, cam)
			}

			if widget.Class == "draggable" {
				draggables = append(draggables, widget)
			}
		}

		if isBgr { // back to gap clipping
			cam.Mask(scx+int(cGapX*cam.Zoom), scy+int(cGapY*cam.Zoom), int(maskW), int(maskH))
		}
	}

	for _, draggable := range draggables {
		if draggable != wPressedOn {
			drawDraggable(draggable, root, cam)
		}
	}

	cam.Mask(scx, scy, int(w*cam.Zoom), int(h*cam.Zoom))
	c.tryShowScroll(cGapX, cGapY, root, cam)
}

func (c *container) tryShowScroll(gapX, gapY float32, root *root, cam *graphics.Camera) {
	var minX, minY, maxX, maxY = c.contentMinMax(gapX, gapY, root)
	var mx, my = cam.MousePosition()
	var focused = c.isFocused(cam)
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)
	var scroll = mouse.ScrollSmooth()

	if minX < c.X || maxX > c.X+c.Width {
		var handleWidth = c.Width / (maxX - minX) * c.Width
		var handleColor = color.Brighten(color.Gray, 0.5)

		if scroll != 0 && focused && shift {
			c.ScrollX -= float32(scroll)
		}

		if c == cMiddlePressed {
			var delta = mx - c.prevMouseX
			c.ScrollX -= delta
			c.velocityX = -delta * dragMomentum
		} else { // apply momentum
			c.ScrollX += c.velocityX * internal.DeltaTime
			c.velocityX *= dragFriction // friction
			if number.Absolute(c.velocityX) < 0.01 {
				c.velocityX = 0
			}
		}

		if focused && isHovered(c.X, c.Y+c.Height-scrollSize, c.Width, scrollSize, cam) {
			wHovered = nil
			wWasHovered = nil
			wFocused = nil
			mouse.SetCursor(cursor.Hand)
			handleColor = color.White

			if mouse.IsButtonJustPressed(b.Left) {
				cPressedOnScrollV = c

				var x = number.Map(c.ScrollX, 0, (maxX-minX)-c.Width, c.X, c.X+c.Width-handleWidth)
				if !isHovered(x, c.Y+c.Height-scrollSize, handleWidth, scrollSize, cam) {
					c.targetScrollX = number.Map(mx-c.X, handleWidth/2, c.Width-handleWidth/2, 0, maxX-minX-c.Width)
				} // clicking on non-handle area moves the handle instantly
			}
		}

		if c == cPressedOnScrollV {
			c.targetScrollX = number.Map(mx-c.X, handleWidth/2, c.Width-handleWidth/2, 0, maxX-minX-c.Width)
			handleColor = color.Gray

			// smooth handle dragging
			var diff = c.targetScrollX - c.ScrollX
			c.ScrollX += diff * handleSpeed * internal.DeltaTime
			if number.Absolute(diff) < 0.5 {
				c.ScrollX = c.targetScrollX
			}
		} else { // the scroll may have changed by MMB dragging or scrolling
			c.targetScrollX = c.ScrollX
		}

		c.ScrollX = number.Limit(c.ScrollX, 0, (maxX-minX)-c.Width)
		var x = number.Map(c.ScrollX, 0, (maxX-minX)-c.Width, c.X, c.X+c.Width-handleWidth)
		cam.DrawQuad(c.X, c.Y+c.Height-scrollSize, c.Width, scrollSize, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawQuad(x, c.Y+c.Height-scrollSize, handleWidth, scrollSize, 0, handleColor)
		cam.DrawQuadFrame(x, c.Y+c.Height-scrollSize, handleWidth, scrollSize, 0, -scrollSize*0.3, color.Black)
	}

	//=================================================================

	if minY < c.Y || maxY > c.Y+c.Height {
		var handleH = (c.Height / (maxY - minY)) * c.Height
		var handleCol = color.Brighten(color.Gray, 0.5)

		if scroll != 0 && focused && !shift {
			c.ScrollY -= float32(scroll)
		}

		if c == cMiddlePressed {
			var delta = my - c.prevMouseY
			c.ScrollY -= delta
			c.velocityY = -delta * dragMomentum
		} else { // apply momentum
			c.ScrollY += c.velocityY * internal.DeltaTime
			c.velocityY *= dragFriction // friction
			if number.Absolute(c.velocityY) < 0.01 {
				c.velocityY = 0
			}
		}

		if focused && isHovered(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, cam) {
			wHovered = nil
			wWasHovered = nil
			wFocused = nil
			mouse.SetCursor(cursor.Hand)
			handleCol = color.White

			if mouse.IsButtonJustPressed(b.Left) {
				cPressedOnScrollH = c

				var y = number.Map(c.ScrollY, 0, (maxY-minY)-c.Height, c.Y, c.Y+c.Height-handleH)
				if !isHovered(c.X+c.Width-scrollSize, y, scrollSize, handleH, cam) {
					c.targetScrollY = number.Map(my-c.Y, handleH/2, c.Height-handleH/2, 0, maxY-minY-c.Height)
				} // clicking on non-handle area moves the handle instantly
			}
		}
		if c == cPressedOnScrollH { // drag handle
			c.targetScrollY += (my - c.prevMouseY) / (handleH / c.Height)
			handleCol = color.Gray

			// smooth handle dragging
			var diff = c.targetScrollY - c.ScrollY
			c.ScrollY += diff * handleSpeed * internal.DeltaTime
			if number.Absolute(diff) < 0.5 {
				c.ScrollY = c.targetScrollY
			}
		} else { // the scroll may have changed by MMB dragging or scrolling
			c.targetScrollY = c.ScrollY
		}

		c.ScrollY = number.Limit(c.ScrollY, 0, (maxY-minY)-c.Height)
		var y = number.Map(c.ScrollY, 0, (maxY-minY)-c.Height, c.Y, c.Y+c.Height-handleH)
		cam.DrawQuad(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawQuad(c.X+c.Width-scrollSize, y, scrollSize, handleH, 0, handleCol)
		cam.DrawQuadFrame(c.X+c.Width-scrollSize, y, scrollSize, handleH, 0, -scrollSize*0.3, color.Black)
	}
	c.prevMouseX, c.prevMouseY = mx, my
}

func (c *container) isHovered(cam *graphics.Camera) bool {
	return isHovered(c.X, c.Y, c.Width, c.Height, cam)
}
func (c *container) isFocused(cam *graphics.Camera) bool {
	return cFocused == c && cWasHovered == c && c.isHovered(cam)
}

func (c *container) contentMinMax(gapX, gapY float32, root *root) (minX, minY, maxX, maxY float32) {
	minX, minY = number.Infinity(), number.Infinity()
	maxX, maxY = number.NegativeInfinity(), number.NegativeInfinity()

	for _, w := range c.Widgets {
		var widget = root.Widgets[w]
		var _, isBgr = widget.Properties[f.FillContainer]
		if isBgr {
			continue
		}

		minX = condition.If(widget.X < minX, widget.X, minX)
		minY = condition.If(widget.Y < minY, widget.Y, minY)
		maxX = condition.If(widget.X+widget.Width > maxX, widget.X+widget.Width, maxX)
		maxY = condition.If(widget.Y+widget.Height > maxY, widget.Y+widget.Height, maxY)
	}
	minX -= gapX
	maxX += gapX
	minY -= gapY
	maxY += gapY
	return
}
