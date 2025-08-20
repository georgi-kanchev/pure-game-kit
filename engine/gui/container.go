package gui

import (
	"encoding/xml"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/symbols"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []*widget  `xml:"Widget"`
	XmlThemes  []*theme   `xml:"Theme"`

	X, Y, Width, Height float32
	Properties          map[string]string
	Widgets             []string
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

func (c *container) UpdateAndDraw(root *root, cam *graphics.Camera) {
	var _, hidden = c.Properties[p.Hidden]
	if hidden {
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

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var _, wHidden = widget.Properties[p.Hidden]
		if wHidden || widget.Class == "tooltip" {
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
}

func (c *container) IsHovered(cam *graphics.Camera) bool {
	return isHovered(c.X, c.Y, c.Width, c.Height, cam)
}

// #endregion
