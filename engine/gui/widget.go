package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"strconv"
	"strings"
)

type widget struct {
	Properties    []xml.Attr `xml:",any,attr"`
	Children      string     `xml:",innerxml"`
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

func (widget *widget) findProp(name string) *xml.Attr {
	for i, v := range widget.Properties {
		if v.Name.Local == name {
			return &widget.Properties[i]
		}
	}
	return nil
}
func (widget *widget) findPropValue(name, defaultValue string) string {
	for i, v := range widget.Properties {
		if v.Name.Local == name {
			return widget.Properties[i].Value
		}
	}
	return defaultValue
}

func (widget *widget) properties(camera *graphics.Camera, owner *container) (x, y, w, h float32, c uint) {
	var px = dyn(camera, owner, widget.findPropValue(property.X, "0"))
	var py = dyn(camera, owner, widget.findPropValue(property.Y, "0"))
	var pw = dyn(camera, owner, widget.findPropValue(property.Width, "0"))
	var ph = dyn(camera, owner, widget.findPropValue(property.Height, "0"))
	var rgba = strings.Split(widget.findPropValue(property.RGBA, "0 0 0 0"), " ")
	var r, g, b, a uint64

	if len(rgba) == 3 || len(rgba) == 4 {
		r, _ = strconv.ParseUint(rgba[0], 10, 8)
		g, _ = strconv.ParseUint(rgba[1], 10, 8)
		b, _ = strconv.ParseUint(rgba[2], 10, 8)
		a = 255
	}
	if len(rgba) == 4 {
		a, _ = strconv.ParseUint(rgba[3], 10, 8)
	}

	var fx, _ = strconv.ParseFloat(px, 32)
	var fy, _ = strconv.ParseFloat(py, 32)
	var fw, _ = strconv.ParseFloat(pw, 32)
	var fh, _ = strconv.ParseFloat(ph, 32)
	var col = color.RGBA(byte(r), byte(g), byte(b), byte(a))
	return float32(fx), float32(fy), float32(fw), float32(fh), col
}

func (widget *widget) draw(cam *graphics.Camera, owner *container) {
	var x, y, w, h, c = widget.properties(cam, owner)
	cam.DrawRectangle(x, y, w, h, 0, c)
}
