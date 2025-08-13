package gui

import (
	"encoding/xml"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/symbols"
)

type container struct {
	XmlProps   []xml.Attr `xml:",any,attr"`
	XmlWidgets []widget   `xml:"Widget"`

	Properties map[string]string
	Widgets    []string
}

func (c *container) UpdateAndDraw(root *root, cam *graphics.Camera) {
	var x, y, w, h = getArea(cam, nil, c.Properties)
	var scx, scy = cam.PointToScreen(float32(x), float32(y))
	var col = getColor(c.Properties)

	cam.Mask(scx, scy, int(w), int(h))
	cam.DrawRectangle(x, y, w, h, 0, col)

	var curX, curY = x, y
	var curOffX, curOffY float32
	var maxHeight float32

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		var wx, wy, ww, wh = getArea(cam, c, widget.Properties)
		var offX = parseNum(dyn(cam, c, widget.Properties[property.OffsetX], "0"))
		var offY = parseNum(dyn(cam, c, widget.Properties[property.OffsetY], "0"))

		var row, newRow = widget.Properties[property.NewRow]
		if newRow {
			curX = x
			curY += parseNum(dyn(cam, c, row, symbols.New(maxHeight+curOffY)))
		}

		curOffX += offX
		curOffY += offY
		widget.Properties[property.X] = symbols.New(curX + curOffX)
		widget.Properties[property.Y] = symbols.New(curY + curOffY)
		curX += ww + curOffX
		maxHeight = condition.If(maxHeight < wh, wh, maxHeight)

		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(cam, &widget, c)

			var text, _ = widget.Properties[property.Text]
			if text != "" {
				var textBox = graphics.NewTextBox("", wx, wy, text)
				textBox.PivotX, textBox.PivotY = 0, 0
				textBox.Width, textBox.Height = ww, wh
				textBox.LineHeight = 50
				cam.DrawTextBoxes(&textBox)
			}
		}
	}

}
