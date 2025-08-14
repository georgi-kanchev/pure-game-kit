package gui

import (
	"encoding/xml"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/symbols"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []widget   `xml:"Widget"`
	XmlThemes  []theme    `xml:"Theme"`

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

func (c *container) UpdateAndDraw(root *root, cam *graphics.Camera) {
	var x = parseNum(dyn(cam, nil, c.Properties[property.X], "0"), 0)
	var y = parseNum(dyn(cam, nil, c.Properties[property.Y], "0"), 0)
	var w = parseNum(dyn(cam, nil, themedProp(property.Width, root, c, nil), "0"), 0)
	var h = parseNum(dyn(cam, nil, themedProp(property.Height, root, c, nil), "0"), 0)
	var scx, scy = cam.PointToScreen(float32(x), float32(y))
	var curX, curY = x, y
	var curOffX, curOffY float32
	var maxHeight float32
	var col = parseColor(themedProp(p.Color, root, c, nil))

	cam.Mask(scx, scy, int(w), int(h))
	cam.DrawRectangle(x, y, w, h, 0, col)

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var wx = parseNum(dyn(cam, c, widget.Properties[property.X], "0"), 0)
		var wy = parseNum(dyn(cam, c, widget.Properties[property.Y], "0"), 0)
		var ww = parseNum(dyn(cam, c, themedProp(property.Width, root, c, &widget), "0"), 0)
		var wh = parseNum(dyn(cam, c, themedProp(property.Height, root, c, &widget), "0"), 0)
		var offX = parseNum(dyn(cam, c, widget.Properties[p.OffsetX], "0"), 0)
		var offY = parseNum(dyn(cam, c, widget.Properties[p.OffsetY], "0"), 0)

		var row, newRow = widget.Properties[p.NewRow]
		if newRow {
			curX = x
			curY += parseNum(dyn(cam, c, row, symbols.New(maxHeight+curOffY)), 0)
		}

		curOffX += offX
		curOffY += offY
		widget.Properties[p.X] = symbols.New(curX + curOffX)
		widget.Properties[p.Y] = symbols.New(curY + curOffY)
		curX += ww + curOffX
		maxHeight = condition.If(maxHeight < wh, wh, maxHeight)

		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(cam, root, &widget, c)

			var text, _ = widget.Properties[p.Text]
			if text != "" {
				var textBox = graphics.NewTextBox("", wx, wy, text)
				textBox.WordWrap = false
				textBox.PivotX, textBox.PivotY = 0, 0
				textBox.Width, textBox.Height = ww, wh
				textBox.FontId = themedProp(p.TextFontId, root, c, &widget)
				textBox.LineHeight = parseNum(themedProp(p.TextLineHeight, root, c, &widget), 60)
				textBox.LineGap = parseNum(themedProp(p.TextLineGap, root, c, &widget), 0)

				cam.DrawTextBoxes(&textBox)
			}
		}
	}
}
