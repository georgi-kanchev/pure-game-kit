package gui

import "encoding/xml"

type root struct {
	XmlName    xml.Name    `xml:"GUI"`
	Containers []container `xml:"Container"`
}

func (root *root) findWidget(widgetId string) (owner, widget *widget) {
	var containers = root.Containers
	for _, c := range containers {
		var widget = c.findWidget(widgetId)
		if widget != nil {
			return &c.widget, widget
		}
	}
	return nil, nil
}
