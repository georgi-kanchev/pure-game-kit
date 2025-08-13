package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
)

// extra props are the inner data between the xml tags such as:
// <Button att="", att2="">
// 		<extraProp>hello</extraProp>
// </Button>
//
// i removed it cuz it clutters and bloats the xml and especially the API for no added value
// other than xml readability & verticality in very few cases

// type extraProp struct {
// 	XMLName xml.Name
// 	Value   string `xml:",chardata"`
// }

type widget struct {
	XmlProps []xml.Attr `xml:",any,attr"`
	// XmlExtraProps []extraProp `xml:",any"`

	Properties    map[string]string
	Owner         string
	UpdateAndDraw func(cam *graphics.Camera, widget *widget, owner *container)
}

func newWidget(class, id, x, y, width, height string, properties ...string) string {
	var result = "<Widget class=\"" + class + "\"" +
		" id=\"" + id + "\"" +
		" x=\"" + x + "\"" +
		" y=\"" + y + "\"" +
		" width=\"" + width + "\"" +
		" height=\"" + height + "\""
	return result + widgetExtraProps(properties...) + " />"

	// if len(extraProps) == 0 {
	// 	return result + " />"
	// }
	// result += ">\n"
	// for _, prop := range extraProps {
	// 	result += "\t\t\t<" + prop[0] + ">" + prop[1] + "</" + prop[0] + ">"
	// }
	// return result + "\n\t\t</Widget>"
}
func widgetExtraProps(props ...string) string {
	var result = ""
	for i, v := range props {
		if i%2 == 0 {
			result += " " + v + "=\""
			continue
		}
		result += v + "\""
	}
	if len(props)%2 != 0 {
		result += "\""
	}
	return result
}

func (widget *widget) draw(cam *graphics.Camera, owner *container) {
	var x, y, w, h = getArea(cam, owner, widget.Properties)
	var col = getColor(widget.Properties)
	cam.DrawRectangle(x, y, w, h, 0, col)
}
