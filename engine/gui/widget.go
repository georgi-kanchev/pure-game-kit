package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
)

type extraProp struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

type widget struct {
	XmlProps      []xml.Attr  `xml:",any,attr"`
	XmlExtraProps []extraProp `xml:",any"`

	Properties    map[string]string
	Owner         string
	UpdateAndDraw func(cam *graphics.Camera, widget *widget, owner *container)
}

func newWidget(class, id, x, y, width, height string, properties, children [][2]string) string {
	var result = "<Widget class=\"" + class + "\"" +
		" id=\"" + id + "\"" +
		" x=\"" + x + "\"" +
		" y=\"" + y + "\"" +
		" width=\"" + width + "\"" +
		" height=\"" + height + "\""
	result += widgetExtraProps(properties)

	if len(children) == 0 {
		return result + " />"
	}

	result += ">\n"
	for _, child := range children {
		result += "\t\t\t<" + child[0] + ">" + child[1] + "</" + child[0] + ">"
	}

	return result + "\n\t\t</Widget>"
}
func widgetExtraProps(props [][2]string) string {
	var result = ""
	for _, v := range props {
		result += " " + v[0] + "=\"" + v[1] + "\""
	}
	return result
}

func (widget *widget) draw(cam *graphics.Camera, owner *container) {
	var x, y, w, h = getArea(cam, owner, widget.Properties)
	var col = getColor(widget.Properties)
	cam.DrawRectangle(x, y, w, h, 0, col)
}
