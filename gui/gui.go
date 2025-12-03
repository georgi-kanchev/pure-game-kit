package gui

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/text"
)

// https://showcase.primefaces.org - basic default browser widgets showcase (scroll down to forms on the left)

type GUI struct {
	Scale  float32
	Volume float32
	root   *root
}

func NewXML(xmlData string) *GUI {
	var gui = GUI{root: &root{}}
	storage.FromXML(xmlData, &gui.root)

	gui.root.Containers = map[string]*container{}
	gui.root.Widgets = map[string]*widget{}
	gui.root.Themes = map[string]*theme{}

	for _, c := range gui.root.XmlContainers {
		var cId = c.XmlProps[0].Value
		c.Widgets = make([]string, len(c.XmlWidgets))
		c.Properties = make(map[string]string, len(c.XmlProps))
		c.WasHidden = true

		for _, xmlProp := range c.XmlProps {
			c.Properties[xmlProp.Name.Local] = xmlProp.Value
		}

		for j, w := range c.XmlWidgets {
			var wClass = w.XmlProps[0].Value
			var wId = w.XmlProps[1].Value
			var fn, has = updateAndDrawFuncs[wClass]
			c.Widgets[j] = wId
			w.OwnerId = cId
			w.Class = wClass
			w.Properties = make(map[string]string, len(w.XmlProps))
			w.Id = wId

			if has {
				w.UpdateAndDraw = fn
			}

			for _, xmlProp := range w.XmlProps {
				w.Properties[xmlProp.Name.Local] = xmlProp.Value
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
func NewElements(elements ...string) *GUI {
	var result = "<GUI scale=\"1\" volume=\"1\">"

	// container is missing on top, add root container
	if len(elements) > 0 && !text.StartsWith(elements[0], "<Container") {
		result += "\n\t<Container " + field.Id + "=\"root\" " +
			field.X + "=\"" + dynamic.CameraLeftX + "\" " +
			field.Y + "=\"" + dynamic.CameraTopY + "\" " +
			field.Width + "=\"" + dynamic.CameraWidth + "\" " +
			field.Height + "=\"" + dynamic.CameraHeight + "\">"
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
	return NewXML(result)
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
		var _, hasTarget = c.Properties[field.TargetId]
		if hasTarget {
			gui.root.cacheTarget(gui.root.themedField(field.TargetId, c, nil))
		}

		var ox = text.New(dyn(nil, c.Properties[field.X], "0"))
		var oy = text.New(dyn(nil, c.Properties[field.Y], "0"))
		var ow = text.New(dyn(nil, c.Properties[field.Width], "0"))
		var oh = text.New(dyn(nil, c.Properties[field.Height], "0"))
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

// works for widgets & containers
func (gui *GUI) SetField(id, field string, value string) {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		w.Properties[field] = value
	}
	if hasC {
		c.Properties[field] = value
	}
	if hasT {
		t.Properties[field] = value
	}
}

//=================================================================

// works for widgets & containers
func (gui *GUI) Field(id, field string) string {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		var owner = gui.root.Containers[w.OwnerId]
		return gui.root.themedField(field, owner, w)
	}
	if hasC {
		return c.Properties[field]
	}
	if hasT {
		return t.Properties[field]
	}

	return ""
}

func (gui *GUI) IsAnyHovered(camera *graphics.Camera) bool {
	for _, c := range gui.root.Containers {
		var hidden = c.Properties[field.Hidden]
		if hidden != "" && c.isHovered(camera) {
			return true
		}
	}

	return false
}

// works for widgets & containers
func (gui *GUI) IsHovered(id string, camera *graphics.Camera) bool {
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

// works for widgets & containers
func (gui *GUI) IsFocused(widgetId string, camera *graphics.Camera) bool {
	var w, has = gui.root.Widgets[widgetId]
	if has {
		return w.isFocused(gui.root, camera)
	}
	return false
}
