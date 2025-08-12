package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
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

	for _, wId := range c.Widgets {
		var widget = root.Widgets[wId]
		if widget.UpdateAndDraw != nil {
			widget.UpdateAndDraw(cam, &widget, c)
		}
	}
}
