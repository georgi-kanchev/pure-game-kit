package gui

import "encoding/xml"

type root struct {
	XmlName       xml.Name    `xml:"GUI"`
	XmlContainers []container `xml:"Container"`

	Themes     map[string]theme
	Containers map[string]container
	Widgets    map[string]widget
}
