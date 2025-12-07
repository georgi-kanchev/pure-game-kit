package gui

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/utility/text"
)

// https://showcase.primefaces.org - basic default browser widgets showcase (scroll down to forms on the left)

type GUI struct {
	Scale  float32
	Volume float32
	root   *root
}

// Joins multiple XMLs into a single GUI - useful for splitting single large files into multiple.
// Keep in mind that the GUI will have the Scale & Volume of only the first XML, the rest are ignored.
func NewFromXMLs(xmlsData ...string) *GUI {
	var gui = GUI{root: &root{}}
	var roots []*root

	for i, xmlData := range xmlsData {
		if xmlData == "" {
			continue
		}

		var root = &root{}
		storage.FromXML(xmlData, &root)

		if i == 0 { // only take scale & volume from the first xml
			gui.root.XmlScale = root.XmlScale
			gui.root.XmlVolume = root.XmlVolume
		}

		roots = append(roots, root)
	}

	gui.root.Containers = map[string]*container{}
	gui.root.Widgets = map[string]*widget{}
	gui.root.Themes = map[string]*theme{}
	gui.root.ContainerIds = []string{}

	for _, r := range roots { // merge contents from all xml roots
		gui.root.XmlContainers = append(gui.root.XmlContainers, r.XmlContainers...)
	}

	for _, c := range gui.root.XmlContainers {
		var cId = c.XmlProps[0].Value
		c.Widgets = make([]string, len(c.XmlWidgets))
		c.Fields = make(map[string]string, len(c.XmlProps))
		c.WasHidden = true

		for _, xmlProp := range c.XmlProps {
			c.Fields[xmlProp.Name.Local] = xmlProp.Value
		}

		for j, w := range c.XmlWidgets {
			var wClass = w.XmlProps[0].Value
			var wId = w.XmlProps[1].Value
			var fn, has = updateAndDrawFuncs[wClass]
			c.Widgets[j] = wId
			w.OwnerId = cId
			w.Class = wClass
			w.Fields = make(map[string]string, len(w.XmlProps))
			w.Id = wId

			if has {
				w.UpdateAndDraw = fn
			}

			for _, xmlProp := range w.XmlProps {
				w.Fields[xmlProp.Name.Local] = xmlProp.Value
			}

			gui.root.Widgets[wId] = w
		}
		for _, t := range c.XmlThemes {
			var tId = t.XmlProps[0].Value
			t.Properties = make(map[string]string, len(t.XmlProps))

			for _, xmlProp := range t.XmlProps {
				t.Properties[xmlProp.Name.Local] = xmlProp.Value
			}
			gui.root.Themes[tId] = t
		}

		gui.root.Containers[cId] = c
		gui.root.ContainerIds = append(gui.root.ContainerIds, cId)
	}

	gui.Scale = gui.root.XmlScale
	gui.Volume = gui.root.XmlVolume
	return &gui
}

// Constructs an XML from a chain of elements (Widgets, Containers and Themes) with Scale & Volume of 1.
// Useful for creating the GUI in an autocompleted code environment instead of in a raw XML file.
//
//	gui.NewFromXMLs(...) // <- put the resulting XML in here to create the GUI
func NewElementsXML(elements ...string) string {
	var result = "<GUI scale=\"1\" volume=\"1\">"

	// container is missing on top, add root container
	if len(elements) > 0 && !text.StartsWith(elements[0], "<Container") {
		result += "\n\t<Container " + f.Id + "=\"root\" " +
			f.X + "=\"" + dynamic.CameraLeftX + "\" " +
			f.Y + "=\"" + dynamic.CameraTopY + "\" " +
			f.Width + "=\"" + dynamic.CameraWidth + "\" " +
			f.Height + "=\"" + dynamic.CameraHeight + "\">"
	}

	for i, v := range elements {
		if text.StartsWith(v, "<Container") {
			if i > 0 {
				result += "\n\t</Container>"
			}
		} else {
			v = "\t" + v
		}

		result += "\n\t" + v

		if i == len(elements)-1 {
			result += "\n\t</Container>"
		}
	}

	result += "\n</GUI>"
	return result
}

// =================================================================

func (gui *GUI) UpdateAndDraw(camera *graphics.Camera) {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	var containers = gui.root.ContainerIds

	gui.root.Volume = gui.Volume

	gui.reset(camera) // keep order of variables & reset

	var tlx, tly = camera.PointFromEdge(0, 0)
	var brx, bry = camera.PointFromEdge(1, 1)
	var cx, cy = camera.PointFromEdge(0.5, 0.5)
	var w, h = camera.Size() // caching dynamic cam props
	camCx, camCy, camLx, camRx = text.New(cx), text.New(cy), text.New(tlx), text.New(brx)
	camTy, camBy, camW, camH = text.New(tly), text.New(bry), text.New(w), text.New(h)

	for _, id := range containers {
		var c = gui.root.Containers[id]
		var _, hasTarget = c.Fields[f.TargetId]
		if hasTarget {
			gui.root.cacheTarget(gui.root.themedField(f.TargetId, c, nil))
		}

		var hidden = dyn(c, c.Fields[f.Hidden], "")
		if hidden != "" { // dyn uses target so needs to be after
			continue
		}

		var ox = text.New(dyn(nil, c.Fields[f.X], "0"))
		var oy = text.New(dyn(nil, c.Fields[f.Y], "0"))
		var ow = text.New(dyn(nil, c.Fields[f.Width], "0"))
		var oh = text.New(dyn(nil, c.Fields[f.Height], "0"))
		ownerLx, ownerRx, ownerTy, ownerBy, ownerW, ownerH = ox, ox+"+"+ow, oy, oy+"+"+oh, ow, oh
		ownerCx, ownerCy = ox+"+"+ow+"/2", oy+"+"+oh+"/2" // caching dynamic props

		c.updateAndDraw(gui.root, camera)
	}

	if cWasHovered == cHovered {
		cFocused = cHovered // only containers that are hovered 2 frames in a row accept input (top-down prio)
	}
	if wWasHovered == wHovered {
		wFocused = wHovered // only widgets that are hovered 2 frames in a row accept input (top-down prio)
	}

	if wPressedOn != nil && wPressedOn.Class == "draggable" {
		camera.Mask(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
		drawDraggable(wPressedOn, gui.root, camera)
	}
	if tooltip != nil {
		camera.Mask(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
		drawTooltip(gui.root, gui.root.Containers[tooltip.OwnerId], camera)
	}

	restore(camera, prevAng, prevZoom, prevX, prevY) // undo what reset does, everything as it was for cam
}

// Works for Widgets & Containers.
func (gui *GUI) SetField(id, field string, value string) {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		w.Fields[field] = value
	}
	if hasC {
		c.Fields[field] = value
	}
	if hasT {
		t.Properties[field] = value
	}
}

//=================================================================

// Works for Widgets & Containers.
func (gui *GUI) Field(id, field string) string {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		var owner = gui.root.Containers[w.OwnerId]
		return gui.root.themedField(field, owner, w)
	}
	if hasC {
		return c.Fields[field]
	}
	if hasT {
		return t.Properties[field]
	}

	return ""
}

func (gui *GUI) IsAnyHovered(camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { restore(camera, prevAng, prevZoom, prevX, prevY) }()
	gui.reset(camera)

	for _, c := range gui.root.Containers {
		var hidden = c.Fields[f.Hidden]
		if hidden == "" && c.isHovered(camera) {
			return true
		}
	}

	return false
}

// Works for Widgets & Containers.
func (gui *GUI) IsHovered(id string, camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { restore(camera, prevAng, prevZoom, prevX, prevY) }()
	gui.reset(camera)

	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]

	if hasW {
		return w.isFocused(gui.root, camera)
	}
	if hasC {
		return c.isFocused(camera)
	}
	return false
}

// Works for Widgets & Containers.
func (gui *GUI) IsFocused(widgetId string, camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { restore(camera, prevAng, prevZoom, prevX, prevY) }()
	gui.reset(camera)

	var w, has = gui.root.Widgets[widgetId]
	if has {
		return w.isFocused(gui.root, camera)
	}
	return false
}
