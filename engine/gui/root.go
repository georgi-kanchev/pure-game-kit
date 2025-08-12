package gui

import "encoding/xml"

type root struct {
	XmlName       xml.Name    `xml:"GUI"`
	XmlContainers []container `xml:"Container"`

	Containers map[string]container
	Widgets    map[string]widget
}
