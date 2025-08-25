package gui

import (
	"encoding/xml"
	"fmt"
	"math"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/dynamic"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/symbols"
	"strconv"
	"strings"
)

type GUI struct {
	root *root
}

func New(elements ...string) *GUI {
	var gui = GUI{root: &root{}}
	var result = "<GUI>"

	// container is missing on top, add root container
	if len(elements) > 0 && !strings.HasPrefix(elements[0], "<Container") {
		result += "\n\t<Container " + property.Id + "=\"root\" " +
			property.X + "=\"" + dynamic.CameraLeftX + "\" " +
			property.Y + "=\"" + dynamic.CameraTopY + "\" " +
			property.Width + "=\"" + dynamic.CameraWidth + "\" " +
			property.Height + "=\"" + dynamic.CameraHeight + "\">"
	}

	for i, v := range elements {
		if strings.HasPrefix(v, "<Container") {
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

	fmt.Printf("%v\n", result)

	var err = xml.Unmarshal([]byte(result), &gui.root)
	fmt.Printf("err: %v\n", err)

	gui.root.Containers = map[string]*container{}
	gui.root.Widgets = map[string]*widget{}
	gui.root.Themes = map[string]*theme{}

	for _, c := range gui.root.XmlContainers {
		var cId = c.XmlProps[0].Value
		c.Widgets = make([]string, len(c.XmlWidgets))
		c.Properties = make(map[string]string, len(c.XmlProps))

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
	}

	return &gui
}

func (gui *GUI) Property(id, property string) string {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		var owner = gui.root.Containers[w.OwnerId]
		return themedProp(property, gui.root, owner, w)
	}
	if hasC {
		return c.Properties[property]
	}
	if hasT {
		return t.Properties[property]
	}

	return ""
}
func (gui *GUI) SetProperty(id, property string, value string) {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]
	var t, hasT = gui.root.Themes[id]

	if hasW {
		w.Properties[property] = value
	}
	if hasC {
		c.Properties[property] = value
	}
	if hasT {
		t.Properties[property] = value
	}
}

func (gui *GUI) Draw(camera *graphics.Camera) {
	var prevAng = camera.Angle
	var containers = gui.root.Containers

	if mouse.IsButtonPressedOnce(mouse.ButtonLeft) {
		pressedOn = nil
		tooltip = nil
	}

	camera.Angle = 0 // force no cam rotation for UI

	if tooltip == nil {
		mouse.SetCursor(mouse.CursorArrow)
	}

	var tlx, tly = camera.PointFromPivot(0, 0)
	var brx, bry = camera.PointFromPivot(1, 1)
	var cx, cy = camera.PointFromPivot(0.5, 0.5)
	var w, h = camera.Size() // caching dynamic cam props
	camCx, camCy = symbols.New(cx), symbols.New(cy)
	camLx, camRx = symbols.New(tlx), symbols.New(brx)
	camTy, camBy = symbols.New(tly), symbols.New(bry)
	camW, camH = symbols.New(w), symbols.New(h)

	for _, c := range containers {
		var ox = symbols.New(dyn(nil, c.Properties[property.X], "0"))
		var oy = symbols.New(dyn(nil, c.Properties[property.Y], "0"))
		var ow = symbols.New(dyn(nil, c.Properties[property.Width], "0"))
		var oh = symbols.New(dyn(nil, c.Properties[property.Height], "0"))
		ownerX, ownerY = ox, oy // caching dynamic owner/container props
		ownerLx, ownerRx = ox, ox+"+"+ow
		ownerTy, ownerBy = oy, oy+"+"+oh
		ownerW, ownerH = ow, oh

		c.UpdateAndDraw(gui.root, camera)
	}

	if tooltip != nil {
		camera.Mask(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
		var tooltipOwner = gui.root.Containers[tooltip.OwnerId]
		drawTooltip(gui.root, tooltipOwner, camera)
	}

	camera.Angle = prevAng // reset angle & mask to how it was
	camera.SetScreenArea(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)

	if mouse.IsButtonReleasedOnce(mouse.ButtonLeft) {
		pressedOn = nil
		tooltip = nil
	}
}

func (gui *GUI) IsHovered(id string, camera *graphics.Camera) bool {
	var w, hasW = gui.root.Widgets[id]
	var c, hasC = gui.root.Containers[id]

	if hasW {
		var owner = gui.root.Containers[w.OwnerId]
		return w.IsHovered(owner, camera)
	}
	if hasC {
		return c.IsHovered(camera)
	}
	return false
}

// #region private

var updateAndDrawFuncs = map[string]func(cam *graphics.Camera, root *root, widget *widget, owner *container){
	"button": button, "slider": slider, "checkbox": checkbox,
}
var camCx, camCy, camLx, camRx, camTy, camBy, camW, camH string               // dynamic prop cache
var ownerX, ownerY, ownerLx, ownerRx, ownerTy, ownerBy, ownerW, ownerH string // dynamic prop cache

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

	var wSelf = ""
	var wTheme, cTheme *theme
	var hasW, hasWt, hasCt = false, false, false

	if w != nil {
		wSelf, hasW = w.Properties[prop]
		wTheme, hasWt = root.Themes[w.ThemeId]
	}
	if c != nil {
		cTheme, hasCt = root.Themes[c.Properties[property.ThemeId]]
	}

	// widget checks
	if hasW {
		return wSelf
	}
	if hasWt {
		return wTheme.Properties[prop]
	}

	if hasCt {
		return cTheme.Properties[prop]
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
	value = strings.ReplaceAll(value, dynamic.CameraCenterX, camCx)
	value = strings.ReplaceAll(value, dynamic.CameraCenterY, camCy)
	value = strings.ReplaceAll(value, dynamic.CameraLeftX, camLx)
	value = strings.ReplaceAll(value, dynamic.CameraRightX, camRx)
	value = strings.ReplaceAll(value, dynamic.CameraTopY, camTy)
	value = strings.ReplaceAll(value, dynamic.CameraBottomY, camBy)
	value = strings.ReplaceAll(value, dynamic.CameraWidth, camW)
	value = strings.ReplaceAll(value, dynamic.CameraHeight, camH)

	if owner != nil {
		value = strings.ReplaceAll(value, dynamic.ContainerX, ownerX)
		value = strings.ReplaceAll(value, dynamic.ContainerY, ownerY)
		value = strings.ReplaceAll(value, dynamic.ContainerWidth, ownerW)
		value = strings.ReplaceAll(value, dynamic.ContainerHeight, ownerH)
		value = strings.ReplaceAll(value, dynamic.ContainerLeftX, ownerLx)
		value = strings.ReplaceAll(value, dynamic.ContainerRightX, ownerRx)
		value = strings.ReplaceAll(value, dynamic.ContainerTopY, ownerTy)
		value = strings.ReplaceAll(value, dynamic.ContainerBottomY, ownerBy)
	}

	var calc = symbols.Calculate(value)
	if math.IsNaN(calc) {
		return defaultValue
	}
	return symbols.New(calc)
}

func parseColor(value string, disabled ...bool) uint {
	var rgba = strings.Split(value, " ")
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

	if len(disabled) == 1 && disabled[0] {
		a /= 3
	}

	return color.RGBA(byte(r), byte(g), byte(b), byte(a))
}
func parseNum(value string, defaultValue float32) float32 {
	var v, err = strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}
	return float32(v)
}

func isHovered(x, y, w, h float32, cam *graphics.Camera) bool {
	var prevAng = cam.Angle
	cam.Angle = 0
	var sx, sy = cam.PointToScreen(x, y)
	var mx, my = cam.PointToScreen(cam.MousePosition())
	var result = mx > sx && mx < sx+int(w) && my > sy && my < sy+int(h)
	cam.Angle = prevAng
	return result
}

// #endregion
