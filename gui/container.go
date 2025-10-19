package gui

import (
	"encoding/xml"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	p "pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []*widget  `xml:"Widget"`
	XmlThemes  []*theme   `xml:"Theme"`

	X, Y, Width, Height, prevMouseX, prevMouseY,
	ScrollX, ScrollY float32
	Properties map[string]string
	Widgets    []string
}

func Container(id, x, y, width, height string, properties ...string) string {
	var rid = "<Container " + p.Id + "=\"" + id + "\""
	var rx = condition.If(x == "", "", " "+p.X+"=\""+x+"\"")
	var ry = condition.If(y == "", "", " "+p.Y+"=\""+y+"\"")
	var rw = condition.If(width == "", "", " "+p.Width+"=\""+width+"\"")
	var rh = condition.If(height == "", "", " "+p.Height+"=\""+height+"\"")
	return rid + rx + ry + rw + rh + extraProps(properties...) + ">"
}

//=================================================================
// private

const scrollSize, scrollSpeed = 20, 100

var cMiddlePressed *container
var cPressedOnScrollH *container
var cPressedOnScrollV *container

func (c *container) updateAndDraw(root *root, cam *graphics.Camera) {
	var hidden, _ = c.Properties[p.Hidden]
	if hidden != "" {
		return
	}

	var x, y, w, h = parseNum(ownerX, 0), parseNum(ownerY, 0), parseNum(ownerW, 0), parseNum(ownerH, 0)
	var scx, scy = cam.PointToScreen(float32(x), float32(y))
	var cGapX = parseNum(dyn(c, c.Properties[p.GapX], "0"), 0)
	var cGapY = parseNum(dyn(c, c.Properties[p.GapY], "0"), 0)
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
	if c.isFocused(cam) && mouse.IsButtonPressedOnce(b.Middle) {
		cMiddlePressed = c
	}

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var wHidden, _ = widget.Properties[p.Hidden]
		if wHidden != "" || widget.Class == "tooltip" {
			continue
		}

		var _, isBgr = widget.Properties[p.FillContainer]
		var ww = parseNum(dyn(c, themedProp(p.Width, root, c, widget), "0"), 0)
		var wh = parseNum(dyn(c, themedProp(p.Height, root, c, widget), "0"), 0)
		var gapX = parseNum(dyn(c, themedProp(p.GapX, root, c, widget), "0"), 0)
		var gapY = parseNum(dyn(c, themedProp(p.GapY, root, c, widget), "0"), 0)
		var offX = parseNum(dyn(c, widget.Properties[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(c, widget.Properties[p.OffsetY], "0"), 0)

		if isBgr {
			widget.X, widget.Y = x, y
			ww, wh = w, h
			cam.Mask(scx, scy, int(w*cam.Zoom), int(h*cam.Zoom)) // gap clipping doesn't affect bgr
		} else {
			var row, newRow = widget.Properties[p.NewRow]
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
		widget.ThemeId = themedProp(p.ThemeId, root, c, widget)

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
	var scroll = mouse.Scroll()

	if minX < c.X || maxX > c.X+c.Width {
		var w = c.Width / (maxX - minX) * c.Width

		if scroll != 0 && focused && shift {
			c.ScrollX -= float32(scroll) * scrollSpeed
		}
		if c == cMiddlePressed {
			c.ScrollX -= mx - c.prevMouseX
		}
		if focused && isHovered(c.X, c.Y+c.Height-scrollSize, c.Width, scrollSize, cam) {
			wHovered = nil
			wWasHovered = nil
			wFocused = nil
			mouse.SetCursor(cursor.Hand)

			if mouse.IsButtonPressedOnce(b.Left) {
				cPressedOnScrollV = c
			}
		}

		if c == cPressedOnScrollV {
			c.ScrollX = number.Map(mx-c.X, w/2, c.Width-w/2, 0, maxX-minX-c.Width)
		}

		c.ScrollX = number.Limit(c.ScrollX, 0, (maxX-minX)-c.Width)
		var x = number.Map(c.ScrollX, 0, (maxX-minX)-c.Width, c.X, c.X+c.Width-w)
		cam.DrawRectangle(c.X, c.Y+c.Height-scrollSize, c.Width, scrollSize, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawRectangle(x, c.Y+c.Height-scrollSize, w, scrollSize, 0, color.White)
		cam.DrawFrame(x, c.Y+c.Height-scrollSize, w, scrollSize, 0, -scrollSize*0.3, color.Black)
	}
	if minY < c.Y || maxY > c.Y+c.Height {
		var h = (c.Height / (maxY - minY)) * c.Height

		if scroll != 0 && focused && !shift {
			c.ScrollY -= float32(scroll) * scrollSpeed
		}
		if c == cMiddlePressed {
			c.ScrollY -= my - c.prevMouseY
		}
		if focused && isHovered(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, cam) {
			wHovered = nil
			wWasHovered = nil
			wFocused = nil
			mouse.SetCursor(cursor.Hand)

			if mouse.IsButtonPressedOnce(b.Left) {
				cPressedOnScrollH = c
			}
		}
		if c == cPressedOnScrollH {
			c.ScrollY = number.Map(my-c.Y, h/2, c.Height-h/2, 0, maxY-minY-c.Height)
		}

		c.ScrollY = number.Limit(c.ScrollY, 0, (maxY-minY)-c.Height)
		var y = number.Map(c.ScrollY, 0, (maxY-minY)-c.Height, c.Y, c.Y+c.Height-h)
		cam.DrawRectangle(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawRectangle(c.X+c.Width-scrollSize, y, scrollSize, h, 0, color.White)
		cam.DrawFrame(c.X+c.Width-scrollSize, y, scrollSize, h, 0, -scrollSize*0.3, color.Black)
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
		var _, isBgr = widget.Properties[p.FillContainer]
		if isBgr {
			continue
		}

		minX = condition.If(widget.X < minX, widget.X, minX)
		minY = condition.If(widget.Y < minY, widget.Y, minY)
		maxX = condition.If(widget.X+widget.Width > maxX, widget.X+widget.Width, maxX)
		maxY = condition.If(widget.Y+widget.Height > maxY, widget.Y+widget.Height, maxY)
	}
	minX -= gapX
	maxX += gapX + scrollSize
	minY -= gapY
	maxY += gapY + scrollSize
	return
}
