/*
The most complex package of all - handling graphical user interfaces by depending heavily on multiple
other packages used for file loading, drawing graphics, accepting input, playing sounds, etc.

The GUI topic is long & thorough and there are many designs, but the few main ones for
games seem to be:
  - Object-oriented (OOP) - Offers the most freedom but way too verbose.
  - Immediate mode (Im) - The simplest one but lacks customization & serialization for re-usability.
  - Data-oriented (css) - Reusable and customizable but hard to parse, create with code, and handle custom logic.

This package takes benefits from each one of them and tries to solve their problems.
It relies on a simple design idea but its usage remains a bit complex due to its sheer depth.
It is constructed of 3 types of elements: Containers, Widgets & Themes.

The GUI creation relies on the Data-oriented approach by parsing XML. Its problems and how they are solved:
  - Hard to parse - by not allowing nesting, having a max depth of 2.
  - Hard to create with code + no autocomplete - by optionally chaining function calls to construct the XML.
  - Hard to handle custom logic - by mixing in the Immediate mode approach.

The Immediate mode approach brings its problems as well, here is how they are solved:
  - Lacking customization - by separating creation details and functionality
  - Lacking serialization for re-usability - by the Data-oriented XML approach
  - Relying on code structures as existing elements - by a single GUI structure & accessing everything through ids.

Another huge problem is GUI elements respecting any window aspect ratios. This is solved by replacing
certain dynamic variable keywords during the XML parsing while handling math expresions.
Due to the nature of those dynamic values, scaling the GUI comes for free by zooming its provided camera.

While loading an XML is handled, saving is not and this is a deliberate choice. Saving a GUI state
has a risk of doing damage to the initial state and has to deal with versioning or multiple GUI states.
Another reason not to do it is that fundamentally it does not make sense to save a GUI.
It's rather better to save its data instead, then load the GUI in its initial state each time and
have it react to the separately loaded data.

Alongside solving all of those problems, here are some of the very useful features in no particular order:

  - Widgets inheriting/reusing properties from their Themes or Container owners and optionally overwriting them.
  - Elements supporting custom properties that only custom logic may rely on.
  - Dividing long & complex GUI systems into multiple XMLs by optionally merging them during parsing.
  - Containers handling scrolling, masking and ordering widgets out-of-the-box.
  - Easy to reference themes, widgets, containers and assets due to the nature of ids.
  - Rendering fallbacks to basic colored shapes in case no assets are provided.
  - Out-of-the-box Z ordering for input & drawing.
  - Having tooltips for all widgets, including text labels & images.
*/
package gui

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

// https://showcase.primefaces.org - basic default browser widgets showcase (scroll down to forms on the left)

type GUI struct {
	Scale  float32
	Volume float32
	root   *root
}

/*
Joins multiple XMLs into a single GUI - useful for splitting single large files into multiple.
Keep in mind that the GUI will have the Scale, Volume & Language of only the first XML root, the rest are ignored.

Pseudo-XML format example:

	GUI // root start
		Container // cannot contain other containers
			Theme // optional for buttons
			Theme // optional for labels
			...   // other themes
		Container // may contain only widgets & themes
			Widget // visual label
			Widget // button
			Widget // visual image
			Widget // slider
			...	   // other widgets
		Container
			Widget // input box
			Widget // check box
		... // other containers
	GUI // root end
*/
func NewFromXMLs(xmlsData ...string) *GUI {
	var gui = GUI{root: &root{}}
	var mergedXML = ""

	for i, xmlData := range xmlsData {
		xmlData = text.Trim(xmlData)
		var lines = text.Split(xmlData, "\n")
		if xmlData == "" || len(lines) < 3 {
			continue
		}

		var startIndex = condition.If(i == 0, 0, 1)
		var portion = lines[startIndex : len(lines)-1]
		xmlData = collection.ToText(portion, "\n")
		mergedXML += xmlData + "\n"
	}
	mergedXML += "</GUI>"

	var root = &root{}
	storage.FromXML(mergedXML, &root)
	gui.root = root
	gui.root.Containers = map[string]*container{}
	gui.root.Widgets = map[string]*widget{}
	gui.root.Themes = map[string]*theme{}
	gui.root.ContainerIds = []string{}

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
			t.Fields = make(map[string]string, len(t.XmlProps))

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
	return &gui
}

/*
Constructs an XML from a chain of elements (Widgets, Containers and Themes) with Scale & Volume of 1 & Language "".
Useful for creating the GUI in an autocompleted code environment instead of in a raw XML file, like so:

	var xml = gui.NewElementsXML(
		gui.Container("menu", "0", "0", "800", "400"),
		gui.Visual("menu-bgr", field.FillContainer, "", field.Color, "200 200 200 255"),
		gui.Button("menu-1", field.ThemeId, "button", field.Text, "Monday"),
		gui.Button("menu-2", field.ThemeId, "button", field.Text, "Tuesday"),
		gui.Button("menu-3", field.ThemeId, "button", field.Text, "Wednesday"),
		gui.Button("menu-4", field.ThemeId, "button", field.Text, "Thursday"),
		gui.Button("menu-5", field.ThemeId, "button", field.Text, "Friday"),
		gui.Visual("weekend-label", field.ThemeId, "label", field.Text, "Weekend"),
		gui.Button("menu-6", field.ThemeId, "button", field.Text, "Saturday"),
		gui.Button("menu-7", field.ThemeId, "button", field.Text, "Sunday")
	)
	var menu = gui.NewFromXMLs(xml) // <-- result
*/
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

func (g *GUI) UpdateAndDraw(camera *graphics.Camera) {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	var containers = g.root.ContainerIds

	g.root.Volume = g.Volume

	g.reset(camera) // keep order of variables & reset

	var tlx, tly = camera.PointFromEdge(0, 0)
	var brx, bry = camera.PointFromEdge(1, 1)
	var cx, cy = camera.PointFromEdge(0.5, 0.5)
	var w, h = camera.Size() // caching dynamic cam props
	camCx, camCy, camLx, camRx = text.New(cx), text.New(cy), text.New(tlx), text.New(brx)
	camTy, camBy, camW, camH = text.New(tly), text.New(bry), text.New(w), text.New(h)

	for _, id := range containers {
		var c = g.root.Containers[id]
		var _, hasTarget = c.Fields[f.TargetId]
		if hasTarget {
			g.root.cacheTarget(g.root.themedField(f.TargetId, c, nil))
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

		c.updateAndDraw(g.root, camera)
	}

	if g.root.cWasHovered == g.root.cHovered {
		g.root.cFocused = g.root.cHovered // only containers hovered 2 frames in a row get input (top-down prio)
	}
	if g.root.wWasHovered == g.root.wHovered {
		g.root.wFocused = g.root.wHovered // only widgets hovered 2 frames in a row get input (top-down prio)
	}

	if g.root.wPressedOn != nil && g.root.wPressedOn.Class == "draggable" {
		camera.Mask(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
		drawDraggable(g.root.wPressedOn, g.root, camera)
	}
	if tooltip != nil {
		camera.Mask(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
		drawTooltip(g.root, g.root.Containers[tooltip.OwnerId], camera)
	}

	clickedId = condition.If(clickedId != "", "", clickedId)
	clickedAndHeldId = condition.If(clickedAndHeldId != "", "", clickedAndHeldId)

	if g.root.wPressedOn != nil {
		if g.root.IsButtonJustClicked(g.root.wPressedOn.Id, camera) {
			clickedId = g.root.wPressedOn.Id
		}
		if g.root.IsButtonClickedAndHeld(g.root.wPressedOn.Id, camera) {
			clickedAndHeldId = g.root.wPressedOn.Id
		}
	}

	g.root.restore(camera, prevAng, prevZoom, prevX, prevY) // undo what reset does, everything as it was for cam
}

// Works for Widgets & Containers.
func (g *GUI) SetField(id, field string, value string) {
	var w, hasW = g.root.Widgets[id]
	var c, hasC = g.root.Containers[id]
	var t, hasT = g.root.Themes[id]

	if hasW {
		w.Fields[field] = value
	}
	if hasC {
		c.Fields[field] = value
	}
	if hasT {
		t.Fields[field] = value
	}
}

//=================================================================

// Works for Widgets & Containers.
func (g *GUI) Field(id, field string) string {
	var w, hasW = g.root.Widgets[id]
	var c, hasC = g.root.Containers[id]
	var t, hasT = g.root.Themes[id]

	if hasW {
		var owner = g.root.Containers[w.OwnerId]
		return g.root.themedField(field, owner, w)
	}
	if hasC {
		return c.Fields[field]
	}
	if hasT {
		return t.Fields[field]
	}

	return ""
}
func (g *GUI) FieldNumber(id, field string) float32 {
	var w, hasW = g.root.Widgets[id]
	var owner *container
	if hasW {
		owner = g.root.Containers[w.OwnerId]
	}
	var value = dyn(owner, g.Field(id, field), "NaN")
	return parseNum(value, number.NaN())
}

func (g *GUI) IsAnyHovered(camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { g.root.restore(camera, prevAng, prevZoom, prevX, prevY) }()
	g.reset(camera)

	for _, c := range g.root.Containers {
		var hidden = c.Fields[f.Hidden]
		if hidden == "" && c.isHovered(camera) {
			return true
		}
	}

	return false
}

// Works for Widgets & Containers.
func (g *GUI) IsHovered(id string, camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { g.root.restore(camera, prevAng, prevZoom, prevX, prevY) }()
	g.reset(camera)

	var w, hasW = g.root.Widgets[id]
	var c, hasC = g.root.Containers[id]

	if hasW {
		return w.isFocused(g.root, camera)
	}
	if hasC {
		return c.isFocused(g.root, camera)
	}
	return false
}

// Works for Widgets & Containers.
func (g *GUI) IsFocused(widgetId string, camera *graphics.Camera) bool {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	defer func() { g.root.restore(camera, prevAng, prevZoom, prevX, prevY) }()
	g.reset(camera)

	var w, has = g.root.Widgets[widgetId]
	if has {
		return w.isFocused(g.root, camera)
	}
	return false
}

func (g *GUI) IdsWidgets() []string {
	return collection.MapKeys(g.root.Widgets)
}
func (g *GUI) IdsContainers() []string {
	return collection.MapKeys(g.root.Containers)
}
func (g *GUI) IdsThemes() []string {
	return collection.MapKeys(g.root.Themes)
}
