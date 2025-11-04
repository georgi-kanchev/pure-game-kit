package gui

import (
	"pure-game-kit/audio"
	"pure-game-kit/data/storage"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
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
// setters

func (gui *GUI) UpdateAndDraw(camera *graphics.Camera) {
	var prevAng, prevZoom, prevX, prevY = camera.Angle, camera.Zoom, camera.X, camera.Y
	var containers = gui.root.ContainerIds

	gui.root.Volume = gui.Volume

	reset(camera, gui) // keep order of variables & reset

	var tlx, tly = camera.PointFromPivot(0, 0)
	var brx, bry = camera.PointFromPivot(1, 1)
	var cx, cy = camera.PointFromPivot(0.5, 0.5)
	var w, h = camera.Size() // caching dynamic cam props
	camCx, camCy, camLx, camRx = text.New(cx), text.New(cy), text.New(tlx), text.New(brx)
	camTy, camBy, camW, camH = text.New(tly), text.New(bry), text.New(w), text.New(h)

	for _, id := range containers {
		var c = gui.root.Containers[id]
		var ox = text.New(dyn(nil, c.Properties[field.X], "0"))
		var oy = text.New(dyn(nil, c.Properties[field.Y], "0"))
		var ow = text.New(dyn(nil, c.Properties[field.Width], "0"))
		var oh = text.New(dyn(nil, c.Properties[field.Height], "0"))
		ownerX, ownerY = ox, oy // caching dynamic owner/container props
		ownerLx, ownerRx, ownerTy, ownerBy, ownerW, ownerH = ox, ox+"+"+ow, oy, oy+"+"+oh, ow, oh

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
// getters

// works for widgets & containers
func (gui *GUI) Field(id, field string) string {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		var owner = gui.root.Containers[w.OwnerId]
		return themedProp(field, gui.root, owner, w)
	}
	if hasC {
		return c.Properties[field]
	}
	if hasT {
		return t.Properties[field]
	}

	return ""
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

//=================================================================
// private

var sound *audio.Audio = audio.New("")
var mouseX, mouseY, prevMouseX, prevMouseY float32
var wFocused, wHovered, wWasHovered *widget
var cFocused, cHovered, cWasHovered *container
var updateAndDrawFuncs = map[string]func(cam *graphics.Camera, root *root, widget *widget){
	"button": button, "slider": slider, "checkbox": checkbox, "menu": menu, "inputField": inputField,
	"draggable": draggable,
}
var camCx, camCy, camLx, camRx, camTy, camBy, camW, camH string               // dynamic prop cache
var ownerX, ownerY, ownerLx, ownerRx, ownerTy, ownerBy, ownerW, ownerH string // dynamic prop cache

func reset(camera *graphics.Camera, gui *GUI) {
	if mouse.IsButtonJustPressed(b.Left) {
		wPressedOn = nil
		tooltip = nil
		cPressedOnScrollH = nil
		cPressedOnScrollV = nil
	}
	if mouse.IsButtonJustReleased(b.Left) {
		cPressedOnScrollH = nil
		cPressedOnScrollV = nil
	}
	if mouse.IsButtonJustReleased(b.Middle) {
		cMiddlePressed = nil
	}

	camera.Zoom = gui.Scale
	camera.Angle = 0          // force no cam rotation for UI
	camera.X, camera.Y = 0, 0 // force no position offset for UI
	mouseX, mouseY = camera.MousePosition()

	if tooltip == nil {
		mouse.SetCursor(cursor.Arrow)
	}
}
func restore(camera *graphics.Camera, prevAng, prevZoom, prevX, prevY float32) {
	camera.Angle, camera.Zoom = prevAng, prevZoom // reset angle, zoom & mask to how it was
	camera.X, camera.Y = prevX, prevY             // also x y
	camera.SetScreenArea(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)

	if mouse.IsButtonJustReleased(b.Left) {
		wPressedOn = nil
		tooltip = nil
	}

	wWasHovered = wHovered
	cWasHovered = cHovered
	prevMouseX, prevMouseY = mouseX, mouseY
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

func themedProp(prop string, root *root, c *container, w *widget) string {
	// priority for widget: widget -> widget theme -> container theme

	var widgetSelf, containerSelf = "", ""
	var widgetTheme, containerTheme *theme
	var hasWidget, hasContainer, hasWidgetTheme, hasContainerTheme = false, false, false, false

	if w != nil {
		widgetSelf, hasWidget = w.Properties[prop]
		widgetTheme, hasWidgetTheme = root.Themes[w.ThemeId]
	}
	if c != nil {
		containerSelf, hasContainer = c.Properties[prop]
		containerTheme, hasContainerTheme = root.Themes[c.Properties[field.ThemeId]]
	}

	if w != nil {
		if hasWidget {
			return widgetSelf
		}
		if hasWidgetTheme {
			return widgetTheme.Properties[prop]
		}
		if hasContainerTheme {
			return containerTheme.Properties[prop]
		}
		if hasContainer {
			return containerSelf
		}
	}

	if hasContainer {
		return containerSelf
	}
	if hasContainerTheme {
		return containerTheme.Properties[prop]
	}

	return ""
}
func defaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func dyn(owner *container, value string, defaultValue string) string {
	value = text.Replace(value, dynamic.CameraCenterX, camCx)
	value = text.Replace(value, dynamic.CameraCenterY, camCy)
	value = text.Replace(value, dynamic.CameraLeftX, camLx)
	value = text.Replace(value, dynamic.CameraRightX, camRx)
	value = text.Replace(value, dynamic.CameraTopY, camTy)
	value = text.Replace(value, dynamic.CameraBottomY, camBy)
	value = text.Replace(value, dynamic.CameraWidth, camW)
	value = text.Replace(value, dynamic.CameraHeight, camH)

	if owner != nil {
		value = text.Replace(value, dynamic.OwnerX, ownerX)
		value = text.Replace(value, dynamic.OwnerY, ownerY)
		value = text.Replace(value, dynamic.OwnerWidth, ownerW)
		value = text.Replace(value, dynamic.OwnerHeight, ownerH)
		value = text.Replace(value, dynamic.OwnerLeftX, ownerLx)
		value = text.Replace(value, dynamic.OwnerRightX, ownerRx)
		value = text.Replace(value, dynamic.OwnerTopY, ownerTy)
		value = text.Replace(value, dynamic.OwnerBottomY, ownerBy)
	}

	// value = strings.ReplaceAll(value, dynamic.MyX, "")
	// value = strings.ReplaceAll(value, dynamic.MyY, "")
	// value = strings.ReplaceAll(value, dynamic.MyWidth, "")
	// value = strings.ReplaceAll(value, dynamic.MyHeight, "")
	// value = strings.ReplaceAll(value, dynamic.MyLeftX, "")
	// value = strings.ReplaceAll(value, dynamic.MyRightX, "")
	// value = strings.ReplaceAll(value, dynamic.MyTopY, "")
	// value = strings.ReplaceAll(value, dynamic.MyBottomY, "")

	var calc = text.Calculate(value)
	if number.IsNaN(calc) {
		return defaultValue
	}
	return text.New(calc)
}

func parseColor(value string, disabled ...bool) uint {
	var rgba = text.Split(value, " ")
	var r, g, b, a byte

	if len(rgba) == 3 || len(rgba) == 4 {
		r = text.ToNumber[byte](rgba[0])
		g = text.ToNumber[byte](rgba[1])
		b = text.ToNumber[byte](rgba[2])
		a = 255
	}
	if len(rgba) == 4 {
		a = text.ToNumber[byte](rgba[3])
	}

	if len(disabled) == 1 && disabled[0] {
		a /= 3
	}

	return color.RGBA(byte(r), byte(g), byte(b), byte(a))
}
func parseNum(value string, defaultValue float32) float32 {
	var v = text.ToNumber[float32](value)
	if number.IsNaN(v) {
		return defaultValue
	}
	return v
}

func isHovered(x, y, w, h float32, cam *graphics.Camera) bool {
	var prevAng = cam.Angle
	cam.Angle = 0
	var sx, sy = cam.PointToScreen(x, y)
	var mx, my = cam.PointToScreen(cam.MousePosition())
	var result = mx > sx && mx < sx+int(w*cam.Zoom) && my > sy && my < sy+int(h*cam.Zoom)
	cam.Angle = prevAng
	return result
}
