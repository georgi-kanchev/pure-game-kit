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
	XmlWidgets []widget   `xml:"Widget"`
	XmlThemes  []theme    `xml:"Theme"`

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

func (c *container) UpdateAndDraw(root *root, cam *graphics.Camera) {
	var x = parseNum(dyn(cam, nil, c.Properties[p.X], "0"), 0)
	var y = parseNum(dyn(cam, nil, c.Properties[p.Y], "0"), 0)
	var w = parseNum(dyn(cam, nil, c.Properties[p.Width], "0"), 0)
	var h = parseNum(dyn(cam, nil, c.Properties[p.Height], "0"), 0)
	var scx, scy = cam.PointToScreen(float32(x), float32(y))
	var cGapX = parseNum(dyn(cam, c, c.Properties[p.GapX], "0"), 0)
	var cGapY = parseNum(dyn(cam, c, c.Properties[p.GapY], "0"), 0)
	var curX, curY = x + cGapX, y + cGapY
	var maxHeight float32

	cam.Mask(scx, scy, int(w), int(h))
	c.X, c.Y, c.Width, c.Height = x, y, w, h

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var ww = parseNum(dyn(cam, c, themedProp(p.Width, root, c, &widget), "0"), 0)
		var wh = parseNum(dyn(cam, c, themedProp(p.Height, root, c, &widget), "0"), 0)
		var gapX = parseNum(dyn(cam, c, themedProp(p.GapX, root, c, &widget), "0"), 0)
		var gapY = parseNum(dyn(cam, c, themedProp(p.GapY, root, c, &widget), "0"), 0)
		var offX = parseNum(dyn(cam, c, widget.Properties[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(cam, c, widget.Properties[p.OffsetY], "0"), 0)
		var _, isBgr = widget.Properties[p.FillContainer]

		if isBgr {
			widget.X, widget.Y = x, y
			ww, wh = w, h
		} else {
			var row, newRow = widget.Properties[p.NewRow]
			if newRow {
				curX = x + cGapX
				curY += parseNum(dyn(cam, c, row, symbols.New(maxHeight+gapY)), 0)
			}

			widget.X = curX + offX
			widget.Y = curY + offY
			curX += ww + gapX
			maxHeight = condition.If(maxHeight < wh, wh, maxHeight)
		}

		widget.Width, widget.Height = ww, wh
		widget.AssetId = themedProp(p.AssetId, root, c, &widget)

		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(ww, wh, cam, root, &widget, c)
		}
	}
}

func (c *container) IsHovered(root *root, cam *graphics.Camera) bool {
	var x, y = cam.PointToScreen(c.X, c.Y)
	var mx, my = cam.PointToScreen(cam.MousePosition())
	return mx > x && mx < x+int(c.Width) && my > y && my < y+int(c.Height)
}
