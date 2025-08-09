package gui

import (
	"encoding/xml"
	"fmt"
	"pure-kit/engine/graphics"
	"strings"
)

type GUI struct {
	data data
}

func New(widgets ...string) GUI {
	var gui data
	var result = "<GUI>"
	var indent = "\t"
	var inContainer = 0

	for i, v := range widgets {
		if v == "</Container>" {
			indent = indent[1:]
			inContainer--
		}

		result += "\n" + indent + v

		if i == len(widgets)-1 {
			for inContainer > 0 {
				indent = indent[1:]
				result += "\n" + indent + "</Container>"
				inContainer--
			}
		}
		if strings.HasPrefix(v, "<Container") {
			indent += "\t"
			inContainer++
		}
	}

	result += "\n</GUI>"

	fmt.Printf("%v\n", result)

	xml.Unmarshal([]byte(result), &gui)
	return GUI{data: gui}
}

func (gui *GUI) Property(widgetId, name string) string {
	var containers = gui.data.RootAndContainers()

	for _, c := range containers {
		var button = c.FindWidget(widgetId)
		if button != nil {
			var prop = button.FindProperty(name)
			if prop != nil {
				return prop.Value
			}
		}
	}

	return ""
}
func (gui *GUI) SetProperty(widgetId, name string, value string) {
	var containers = gui.data.RootAndContainers()

	for _, c := range containers {
		var button = c.FindWidget(widgetId)
		if button != nil {
			var prop = button.FindProperty(name)
			if prop != nil {
				prop.Value = value
			}
		}
	}
}

func (gui *GUI) Draw(camera *graphics.Camera) {
	// var containers = gui.data.RootAndContainers()
	// for _, c := range containers {
	// var owner = replaceDynamics(camera)
	// var x = replaceDynamics(c.FindProperty(property.X), camera)
	// }
}

// #region private

type data struct {
	XmlName    xml.Name    `xml:"GUI"`
	Containers []container `xml:"Container"`
	Buttons    []button    `xml:"Button"`
}
type widget struct {
	Properties []xml.Attr `xml:",any,attr"`
}

func (data *data) RootAndContainers() []container {
	var containers = make([]container, 0, len(data.Containers)+1)
	containers = append(containers, container{Buttons: data.Buttons})
	containers = append(containers, data.Containers...)
	return containers
}

func (widget *widget) FindProperty(name string) *xml.Attr {
	for i, v := range widget.Properties {
		if v.Name.Local == name {
			return &widget.Properties[i]
		}
	}
	return nil
}

func newWidget(class, id, x, y, width, height string) string {
	return "<" + class + " id=\"" + id + "\"" +
		" x=\"" + x + "\"" +
		" y=\"" + y + "\"" +
		" width=\"" + width + "\"" +
		" height=\"" + height + "\""
}
func extraProps(props ...string) string {
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

func replaceDynamics(cam *graphics.Camera, owner container, value string) string {
	return ""
}

// #endregion
