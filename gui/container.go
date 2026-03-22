package gui

import (
	"encoding/xml"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	f "pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []*widget  `xml:"Widget"`
	XmlThemes  []*theme   `xml:"Theme"`

	X, Y, Width, Height,
	ScrollX, ScrollY float32

	prevMouseX, prevMouseY,
	dragVelX, dragVelY,
	targetScrollX, targetScrollY float32

	mask      *graphics.Area
	root      *root
	Fields    map[string]string
	Widgets   []string
	WasHidden bool
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

const scrollSize, handleSpeed, dragFriction, dragMomentum = 10.0, 12.0, 0.95, 30.0

var rowWidths = map[*widget]float32{}

func (c *container) updateAndDraw() {
	var x, y, w, h = parseNum(ownerLx, 0), parseNum(ownerTy, 0), parseNum(ownerW, 0), parseNum(ownerH, 0)
	var anchorX = parseNum(c.Fields[field.AnchorX], 0)
	var anchorY = parseNum(c.Fields[field.AnchorY], 0)
	var cGapX = parseNum(c.Fields[f.GapX], 0)
	var cGapY = parseNum(c.Fields[f.GapY], 0)

	if c.mask == nil {
		c.mask = graphics.NewArea(0, 0, 0, 0)
	} // create only once cuz textBoxes have prop cache and will regenerate when assigning them this mask
	c.mask.X, c.mask.Y = x+cGapX, y+cGapY
	c.mask.Width, c.mask.Height = w-cGapX*2, h-cGapY*2
	c.X, c.Y, c.Width, c.Height = x, y, w, h

	if c.isHovered() {
		c.root.cHovered = c
	}
	if c.isFocused() && mouse.IsButtonJustPressed(b.Middle) {
		c.root.cMiddlePressed = c
	}

	c.alignWidgets(x, y, w, h, cGapX, cGapY)

	//=================================================================
	// this is done in two loops because the content size of the container needs to be known to anchor it

	var minX, minY, maxX, maxY = c.contentMinMax(cGapX, cGapY)
	var contentW, contentH = maxX - minX, maxY - minY
	var draggables []*widget = make([]*widget, 0, len(c.Widgets))
	for _, wId := range c.Widgets {
		var widget = c.root.Widgets[wId]
		if widget.IsCulled || widget.isSkipped(c) {
			continue
		}

		var _, isBgr = widget.Fields[f.FillContainer]
		if !isBgr {
			if contentW <= c.Width {
				widget.X += c.Width*anchorX - contentW*anchorX
				widget.X += (contentW - rowWidths[widget]) * anchorX
			}
			if contentH <= c.Height {
				widget.Y += c.Height*anchorY - contentH*anchorY
			}
		}

		if widget.isHovered(c) {
			c.root.wHovered = widget
		}

		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(widget)
			tryShowTooltip(widget, c)
		} else if widget.Class == "visual" {
			setupVisualsTextured(widget)
			setupVisualsText(widget, true)
			drawVisuals(widget, false, nil)
			tryShowTooltip(widget, c)
		}

		if widget.Class == "draggable" {
			draggables = append(draggables, widget)
		}
	}

	for _, draggable := range draggables {
		if draggable != c.root.wPressedOn {
			drawDraggable(draggable)
		}
	}

	c.tryShowScrolls(minX, minY, maxX, maxY)
}
func (c *container) alignWidgets(x, y, w, h, cGapX, cGapY float32) {
	var maxHeight float32
	var curX, curY = x + cGapX, y + cGapY
	var rowWidth = cGapX * 2
	var rowWidgets []*widget
	var nonBgrIndex = 0 // new row shouldn't work for first widget, used to check first nonBgr widget

	collection.MapClear(rowWidths)
	for _, wId := range c.Widgets {
		var widget = c.root.Widgets[wId]
		if widget.isSkipped(c) {
			continue
		}

		var _, isBgr = widget.Fields[f.FillContainer]
		var ww = parseNum(dyn(c, c.root.themedField(f.Width, c, widget), "0"), 0)
		var wh = parseNum(dyn(c, c.root.themedField(f.Height, c, widget), "0"), 0)
		var gapX = parseNum(c.root.themedField(f.GapX, c, widget), 0)
		var gapY = parseNum(c.root.themedField(f.GapY, c, widget), 0)
		var offX = parseNum(widget.Fields[f.OffsetX], 0)
		var offY = parseNum(widget.Fields[f.OffsetY], 0)

		if isBgr {
			widget.X, widget.Y = x, y
			ww, wh = w, h
		} else {
			var row, newRow = widget.Fields[f.NewRow]
			if newRow && nonBgrIndex > 0 {
				for _, w := range rowWidgets {
					rowWidths[w] = rowWidth
				}
				rowWidgets = nil
				rowWidth = cGapX * 2

				curX = x + cGapX
				curY += parseNum(row, maxHeight+gapY)
				maxHeight = 0
			}

			gapX = condition.If(newRow || nonBgrIndex == 0, 0, gapX)
			curX += gapX
			widget.X = curX + offX
			widget.Y = curY + offY
			curX += ww
			maxHeight = condition.If(maxHeight < wh, wh, maxHeight)
			nonBgrIndex++
		}

		widget.Width, widget.Height = ww, wh
		widget.ThemeId = c.root.themedField(f.ThemeId, c, widget)

		if !isBgr {
			widget.X -= c.ScrollX
			widget.Y -= c.ScrollY

			rowWidth += widget.Width + gapX
			rowWidgets = append(rowWidgets, widget)
		}

		var outsideX = widget.X+widget.Width < c.X || widget.X > c.X+c.Width
		var outsideY = widget.Y+widget.Height < c.Y || widget.Y > c.Y+c.Height
		widget.IsCulled = outsideX || outsideY
	}
	for _, w := range rowWidgets { // apply for last row too, no new row to trigger it
		rowWidths[w] = rowWidth
	}
}

func (c *container) tryShowScrolls(minX, minY, maxX, maxY float32) {
	var mx, my = c.root.cam.MousePosition()
	var focused = c.isFocused()
	var shift = keyboard.IsKeyPressed(key.LeftShift) || keyboard.IsKeyPressed(key.RightShift)
	var scroll = mouse.ScrollSmooth()
	var horizontal = minX+1 < c.X || maxX-1 > c.X+c.Width
	var vertical = minY+1 < c.Y || maxY-1 > c.Y+c.Height

	if horizontal && !vertical {
		shift = true // when only horizontal scroll is present, no need to press shift
	}

	if mouse.Scroll() != 0 && focused {
		c.root.cScrolledOn = c
	}
	if c.root.cWasScrolling == c && !focused {
		c.root.cScrolledOn = nil
	}

	if horizontal {
		c.handleHorizontalSlider(maxX, minX, vertical, shift)
	}
	if vertical {
		c.handleVerticalSlider(maxY, minY, shift)
	}

	if scroll != 0 && focused {
		c.root.cWasScrolling = c
	}

	c.prevMouseX, c.prevMouseY = mx, my
}
func (c *container) handleVerticalSlider(maxY, minY float32, shift bool) {
	var cam = c.root.cam
	var focused = c.isFocused()
	var scroll = mouse.ScrollSmooth()
	var _, my = cam.MousePosition()
	var handleH = (c.Height / (maxY - minY)) * c.Height
	var handleCol = color.Brighten(palette.Gray, 0.5)

	if scroll != 0 && focused && !shift && c.root.cScrolledOn == c {
		c.ScrollY -= float32(scroll)
	}

	if c == c.root.cMiddlePressed {
		var dy = my - c.prevMouseY
		c.ScrollY -= dy
		var instantVelY = -dy / internal.DeltaTime
		const weight = 0.2
		c.dragVelY = (c.dragVelY * (1.0 - weight)) + (instantVelY * weight)
	} else {
		c.ScrollY += c.dragVelY * internal.DeltaTime
		var decay = number.Exponential(-10.0 * internal.DeltaTime)
		c.dragVelY *= decay

		if number.Absolute(c.dragVelY) < 0.1 {
			c.dragVelY = 0
		}
	}

	if focused && isHovered(c.X+c.Width-scrollSize, c.Y, scrollSize, c.Height, cam) {
		c.root.wHovered = nil
		c.root.wWasHovered = nil
		c.root.wFocused = nil
		mouse.SetCursor(cursor.Hand)
		handleCol = palette.White

		if mouse.IsButtonJustPressed(b.Left) {
			c.root.cPressedOnScrollV = c

			var y = number.Map(c.ScrollY, 0, (maxY-minY)-c.Height, c.Y, c.Y+c.Height-handleH)
			if !isHovered(c.X+c.Width-scrollSize, y, scrollSize, handleH, cam) {
				c.targetScrollY = number.Map(my-c.Y, handleH/2, c.Height-handleH/2, 0, maxY-minY-c.Height)
			} // clicking on non-handle area moves the handle instantly
		}
	}
	if c == c.root.cPressedOnScrollV { // drag handle
		c.targetScrollY += (my - c.prevMouseY) / (handleH / c.Height)
		handleCol = palette.Gray

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
	cam.DrawQuadFrame(c.X+c.Width-scrollSize, y, scrollSize, handleH, 0, -scrollSize*0.3, palette.Black)
}
func (c *container) handleHorizontalSlider(maxX, minX float32, vertical, shift bool) {
	var cam = c.root.cam
	var focused = c.isFocused()
	var scroll = mouse.ScrollSmooth()
	var mx, _ = cam.MousePosition()
	var barW = condition.If(vertical, c.Width-scrollSize, c.Width) // make space for vertical scroll
	var handleW = barW / (maxX - minX) * barW
	var handleCol = color.Brighten(palette.Gray, 0.5)

	if scroll != 0 && focused && shift && c.root.cScrolledOn == c {
		c.ScrollX -= float32(scroll)
	}

	if c == c.root.cMiddlePressed {
		var dx = mx - c.prevMouseX
		c.ScrollX -= dx
		var instantVelX = -dx / internal.DeltaTime
		const weight = 0.2
		c.dragVelX = (c.dragVelX * (1.0 - weight)) + (instantVelX * weight)
	} else {
		c.ScrollX += c.dragVelX * internal.DeltaTime
		var decay = number.Exponential(-10.0 * internal.DeltaTime)
		c.dragVelX *= decay

		if number.Absolute(c.dragVelX) < 0.1 {
			c.dragVelX = 0
		}
	}

	if focused && isHovered(c.X, c.Y+c.Height-scrollSize, barW, scrollSize, cam) {
		c.root.wHovered = nil
		c.root.wWasHovered = nil
		c.root.wFocused = nil
		mouse.SetCursor(cursor.Hand)
		handleCol = palette.White

		if mouse.IsButtonJustPressed(b.Left) {
			c.root.cPressedOnScrollH = c

			var x = number.Map(c.ScrollX, 0, (maxX-minX)-barW, c.X, c.X+barW-handleW)
			if !isHovered(x, c.Y+c.Height-scrollSize, handleW, scrollSize, cam) {
				c.targetScrollX = number.Map(mx-c.X, handleW/2, barW-handleW/2, 0, maxX-minX-barW)
			} // clicking on non-handle area moves the handle instantly
		}
	}

	if c == c.root.cPressedOnScrollH {
		c.targetScrollX += (mx - c.prevMouseX) / (handleW / barW)
		handleCol = palette.Gray

		// smooth handle dragging
		var diff = c.targetScrollX - c.ScrollX
		c.ScrollX += diff * handleSpeed * internal.DeltaTime
		if number.Absolute(diff) < 0.5 {
			c.ScrollX = c.targetScrollX
		}
	} else { // the scroll may have changed by MMB dragging or scrolling
		c.targetScrollX = c.ScrollX
	}

	c.ScrollX = number.Limit(c.ScrollX, 0, (maxX-minX)-barW)
	var x = number.Map(c.ScrollX, 0, (maxX-minX)-barW, c.X, c.X+barW-handleW)
	cam.DrawQuad(c.X, c.Y+c.Height-scrollSize, barW, scrollSize, 0, color.RGBA(0, 0, 0, 150))
	cam.DrawQuad(x, c.Y+c.Height-scrollSize, handleW, scrollSize, 0, handleCol)
	cam.DrawQuadFrame(x, c.Y+c.Height-scrollSize, handleW, scrollSize, 0, -scrollSize*0.3, palette.Black)
}

func (c *container) isHovered() bool {
	return isHovered(c.X, c.Y, c.Width, c.Height, c.root.cam)
}
func (c *container) isFocused() bool {
	return c.root.cFocused == c && c.root.cWasHovered == c && c.isHovered()
}
func (c *container) contentMinMax(gapX, gapY float32) (minX, minY, maxX, maxY float32) {
	minX, minY = number.Infinity(), number.Infinity()
	maxX, maxY = number.NegativeInfinity(), number.NegativeInfinity()

	for _, w := range c.Widgets {
		var widget = c.root.Widgets[w]
		var _, isBgr = widget.Fields[f.FillContainer]
		if isBgr || widget.isSkipped(c) { // even culled items are calculated for
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
