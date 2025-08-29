package gui

import (
	"encoding/xml"
	"math"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/symbols"
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

// #region private

const scrollSize, scrollSpeed = 20, 100

func (c *container) UpdateAndDraw(root *root, cam *graphics.Camera) {
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

	cam.Mask(scx, scy, int(w), int(h))
	c.X, c.Y, c.Width, c.Height = x, y, w, h

	if c.IsHovered(cam) {
		cHovered = c
	}

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var wHidden, _ = widget.Properties[p.Hidden]
		if wHidden != "" || widget.Class == "tooltip" {
			continue
		}

		var ww = parseNum(dyn(c, themedProp(p.Width, root, c, widget), "0"), 0)
		var wh = parseNum(dyn(c, themedProp(p.Height, root, c, widget), "0"), 0)
		var gapX = parseNum(dyn(c, themedProp(p.GapX, root, c, widget), "0"), 0)
		var gapY = parseNum(dyn(c, themedProp(p.GapY, root, c, widget), "0"), 0)
		var offX = parseNum(dyn(c, widget.Properties[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(c, widget.Properties[p.OffsetY], "0"), 0)
		var _, isBgr = widget.Properties[p.FillContainer]

		if isBgr {
			widget.X, widget.Y = x, y
			ww, wh = w, h
		} else {
			var row, newRow = widget.Properties[p.NewRow]
			if newRow {
				curX = x + cGapX
				curY += parseNum(dyn(c, row, symbols.New(maxHeight+gapY)), 0)
			}

			curX += gapX
			widget.X = curX + offX
			widget.Y = curY + offY
			curX += ww
			maxHeight = condition.If(maxHeight < wh, wh, maxHeight)
		}

		widget.Width, widget.Height = ww, wh
		widget.ThemeId = themedProp(p.ThemeId, root, c, widget)

		if !isBgr {
			widget.X -= c.ScrollX
			widget.Y -= c.ScrollY
		}

		if widget.IsHovered(c, cam) {
			wHovered = widget
		}

		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(cam, root, widget, c)
			tryShowTooltip(widget, root, c, cam)
		} else if widget.Class == "visual" {
			setupVisualsTextured(root, widget, c)
			setupVisualsText(root, widget, c)
			drawVisuals(cam, root, widget, c)
			tryShowTooltip(widget, root, c, cam)
		}
	}

	c.TryShowScroll(cGapX, cGapY, root, cam)
}

func (c *container) IsHovered(cam *graphics.Camera) bool {
	return isHovered(c.X, c.Y, c.Width, c.Height, cam)
}

func (c *container) TryShowScroll(gapX, gapY float32, root *root, cam *graphics.Camera) {
	var minX, minY = float32(math.Inf(1)), float32(math.Inf(1))
	var maxX, maxY = float32(math.Inf(-1)), float32(math.Inf(-1))

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

	var mx, my = cam.MousePosition()
	var focused = c.IsFocused(root, cam)
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)

	if minX < c.X+1 || maxX > c.X+c.Width-1 {
		var scroll = mouse.Scroll()
		var ratio = c.Width / (maxX - minX)
		var w = ratio * c.Width

		if scroll != 0 && focused && shift {
			c.ScrollX -= float32(scroll) * scrollSpeed
		}
		if mouse.IsButtonPressed(mouse.ButtonMiddle) && focused {
			c.ScrollX -= mx - c.prevMouseX
		}

		c.ScrollX = number.Limit(c.ScrollX, 0, (maxX-minX)-c.Width)
		var x = number.Map(c.ScrollX, 0, (maxX-minX)-c.Width, c.X, c.X+c.Width-w)
		cam.DrawRectangle(c.X, c.Y+c.Height-scrollSize, c.Width, scrollSize, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawRectangle(x, c.Y+c.Height-scrollSize, w, scrollSize, 0, color.White)
		cam.DrawFrame(x, c.Y+c.Height-scrollSize, w, scrollSize, 0, -scrollSize*0.3, color.Black)
	}
	if minY < c.Y || maxY > c.Y+c.Height {
		var scroll = mouse.Scroll()
		var ratio = c.Height / (maxY - minY)
		var h = ratio * c.Height

		if scroll != 0 && focused && !shift {
			c.ScrollY -= float32(scroll) * scrollSpeed
		}
		if mouse.IsButtonPressed(mouse.ButtonMiddle) && focused {
			c.ScrollY -= my - c.prevMouseY
		}

		c.ScrollY = number.Limit(c.ScrollY, 0, (maxY-minY)-c.Height)
		var y = number.Map(c.ScrollY, 0, (maxY-minY)-c.Height, c.Y, c.Y+c.Height-h)
		cam.DrawRectangle(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, 0, color.RGBA(0, 0, 0, 150))
		cam.DrawRectangle(c.X+c.Width-scrollSize, y, scrollSize, h, 0, color.White)
		cam.DrawFrame(c.X+c.Width-scrollSize, y, scrollSize, h, 0, -scrollSize*0.3, color.Black)
	}
	c.prevMouseX, c.prevMouseY = mx, my
}

func (c *container) IsFocused(root *root, cam *graphics.Camera) bool {
	return cFocused == c && cWasHovered == c && c.IsHovered(cam)
}

// #endregion
