package gui

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/text"
)

// Joins multiple XMLs into a single GUI - useful for splitting single large files into multiple.
// Keep in mind that the GUI will have the Scale, Volume & Language of only the first XML root, the rest are ignored.
//
// Pseudo-XML format example:
//
//	GUI // root start
//		Container // cannot contain other containers
//			Theme // optional for buttons
//			Theme // optional for labels
//			...   // other themes
//		Container // may contain only widgets & themes
//			Widget // visual label
//			Widget // button
//			Widget // visual image
//			Widget // slider
//			...	   // other widgets
//		Container
//			Widget // input box
//			Widget // check box
//		... // other containers
//	GUI // root end
func NewFromXMLs(camera *graphics.Camera, xmlsData ...string) *GUI {
	var gui = &GUI{root: &root{cam: camera}}
	var mergedXML = text.NewBuilder()

	for i, xmlData := range xmlsData {
		xmlData = text.Trim(xmlData)
		var lines = text.Split(xmlData, "\n")
		if xmlData == "" || len(lines) < 3 {
			continue
		}

		var startIndex = condition.If(i == 0, 0, 1)
		var portion = lines[startIndex : len(lines)-1]
		xmlData = collection.ToText(portion, "\n")
		mergedXML.WriteText(xmlData + "\n")
	}
	mergedXML.WriteText("</GUI>")

	storage.FromXML(mergedXML.ToText(), &gui.root)
	gui.root.Containers = map[string]*container{}
	gui.root.Widgets = map[string]*widget{}
	gui.root.Themes = map[string]*theme{}
	gui.root.ContainerIds = []string{}

	for _, c := range gui.root.XmlContainers {
		var cId = c.XmlProps[0].Value
		c.Widgets = make([]string, len(c.XmlWidgets))
		c.Fields = make(map[string]string, len(c.XmlProps))
		c.WasHidden = true
		c.root = gui.root

		for _, xmlProp := range c.XmlProps {
			c.Fields[xmlProp.Name.Local] = xmlProp.Value
		}

		for j, w := range c.XmlWidgets {
			var wClass = w.XmlProps[0].Value
			var wId = w.XmlProps[1].Value
			var fn, has = updates[wClass]
			c.Widgets[j] = wId
			w.OwnerId = cId
			w.Class = wClass
			w.Fields = make(map[string]string, len(w.XmlProps))
			w.Id = wId
			w.holdId = ";;hold-" + wId
			w.root = gui.root

			if has {
				w.Update = fn
			}

			for _, xmlProp := range w.XmlProps {
				w.Fields[xmlProp.Name.Local] = xmlProp.Value
			}

			gui.root.Widgets[wId] = w
		}
		for _, t := range c.XmlThemes {
			var tId = t.XmlProps[0].Value
			t.Fields = make(map[string]string, len(t.XmlProps))
			t.Root = gui.root

			for _, xmlProp := range t.XmlProps {
				t.Fields[xmlProp.Name.Local] = xmlProp.Value
			}
			gui.root.Themes[tId] = t
		}

		gui.root.Containers[cId] = c
		gui.root.ContainerIds = append(gui.root.ContainerIds, cId)
	}

	gui.Scale = gui.root.XmlScale
	gui.Volume = gui.root.XmlVolume
	return gui
}

// Constructs an XML from a chain of elements (Widgets, Containers and Themes) with Scale & Volume of 1 & Language "".
// Useful for creating the GUI in an autocompleted code environment instead of in a raw XML file, like so:
//
//	var xml = gui.NewElementsXML(
//		gui.Container("menu", "0", "0", "800", "400"),
//		gui.Visual("menu-bgr", field.FillContainer, "", field.Color, "200 200 200 255"),
//		gui.Button("menu-1", field.ThemeId, "button", field.Text, "Monday"),
//		gui.Button("menu-2", field.ThemeId, "button", field.Text, "Tuesday"),
//		gui.Button("menu-3", field.ThemeId, "button", field.Text, "Wednesday"),
//		gui.Button("menu-4", field.ThemeId, "button", field.Text, "Thursday"),
//		gui.Button("menu-5", field.ThemeId, "button", field.Text, "Friday"),
//		gui.Visual("weekend-label", field.ThemeId, "label", field.Text, "Weekend"),
//		gui.Button("menu-6", field.ThemeId, "button", field.Text, "Saturday"),
//		gui.Button("menu-7", field.ThemeId, "button", field.Text, "Sunday")
//	)
//	var menu = gui.NewFromXMLs(xml) // <-- result
func NewElementsXML(elements ...string) string {
	var result = text.NewBuilder()
	result.WriteText("<GUI scale=\"1\" volume=\"1\">")

	// container is missing on top, add root container
	if len(elements) > 0 && !text.StartsWith(elements[0], "<Container") {
		result.WriteText("\n\t<Container " + f.Id + "=\"root\" " +
			f.X + "=\"" + dynamic.CameraLeftX + "\" " +
			f.Y + "=\"" + dynamic.CameraTopY + "\" " +
			f.Width + "=\"" + dynamic.CameraWidth + "\" " +
			f.Height + "=\"" + dynamic.CameraHeight + "\">")
	}

	for i, v := range elements {
		if text.StartsWith(v, "<Container") {
			if i > 0 {
				result.WriteText("\n\t</Container>")
			}
		} else {
			v = "\t" + v
		}

		result.WriteText("\n\t" + v)

		if i == len(elements)-1 {
			result.WriteText("\n\t</Container>")
		}
	}

	result.WriteText("\n</GUI>")
	return result.ToText()
}
