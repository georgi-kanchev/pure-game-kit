package gui

import (
	"encoding/xml"
	"fmt"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/dynamic"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/symbols"
	"strconv"
	"strings"
)

type GUI struct {
	root root
}

func New(widgets ...string) GUI {
	var gui root
	var result = "<GUI>"

	// container is missing on top, add root container
	if len(widgets) > 0 && !strings.HasPrefix(widgets[0], "<Container") {
		result += "\n\t<Container " + property.Id + "=\"root\" " +
			property.X + "=\"" + dynamic.CameraLeftX + "\" " +
			property.Y + "=\"" + dynamic.CameraTopY + "\" " +
			property.Width + "=\"" + dynamic.CameraWidth + "\" " +
			property.Height + "=\"" + dynamic.CameraHeight + "\">"
	}

	for i, v := range widgets {
		if strings.HasPrefix(v, "<Container") {
			if i > 0 {
				result += "\n\t</Container>"
			}
		} else {
			v = "\t" + v
		}

		result += "\n\t" + v

		if i == len(widgets)-1 {
			result += "\n\t</Container>"
		}
	}

	result += "\n</GUI>"

	fmt.Printf("%v\n", result)

	xml.Unmarshal([]byte(result), &gui)

	gui.Containers = map[string]container{}
	gui.Widgets = map[string]widget{}

	for i := range gui.XmlContainers {
		var c = &gui.XmlContainers[i]
		var cId = c.XmlProps[0].Value
		c.Widgets = make([]string, len(c.XmlWidgets))
		c.Properties = make(map[string]string, len(c.XmlProps))

		for _, xmlProp := range c.XmlProps {
			c.Properties[xmlProp.Name.Local] = xmlProp.Value
		}

		for j := range c.XmlWidgets {
			var w = &c.XmlWidgets[j]
			var wClass = w.XmlProps[0].Value
			var wId = w.XmlProps[1].Value
			var fn, has = updateAndDrawFuncs[wClass]
			c.Widgets[j] = wId
			w.Owner = cId
			w.Properties = make(map[string]string, len(w.XmlProps)+len(w.XmlExtraProps))

			if has {
				w.UpdateAndDraw = fn
			}

			for _, xmlProp := range w.XmlProps {
				w.Properties[xmlProp.Name.Local] = xmlProp.Value
			}
			for _, xmlProp := range w.XmlExtraProps {
				w.Properties[xmlProp.XMLName.Local] = xmlProp.Value
			}

			gui.Widgets[wId] = *w
		}
		gui.Containers[cId] = *c
	}

	return GUI{root: gui}
}
func Container(id, x, y, width, height string, properties [][2]string) string {
	var props = "<Container id=\"" + id + "\"" +
		" x=\"" + x + "\"" +
		" y=\"" + y + "\"" +
		" width=\"" + width + "\"" +
		" height=\"" + height + "\""
	return props + widgetExtraProps(properties) + ">"
}
func NewButton(id, x, y, width, height string, properties [][2]string, children [][2]string) string {
	return newWidget("button", id, x, y, width, height, properties, children)
}

func (gui *GUI) Property(widgetId, property string) string {
	var widget, has = gui.root.Widgets[widgetId]
	if !has {
		return ""
	}

	return widget.Properties[property]
}
func (gui *GUI) SetProperty(widgetId, property string, value string) {
	var widget, has = gui.root.Widgets[widgetId]
	if !has {
		return
	}

	widget.Properties[property] = value
}

func (gui *GUI) Draw(camera *graphics.Camera) {
	var prevAng = camera.Angle
	var containers = gui.root.XmlContainers
	var prevScX, prevScY = camera.ScreenX, camera.ScreenY // each container masks an area, store to turn back later
	var prevScW, prevScH = camera.ScreenWidth, camera.ScreenHeight

	camera.Angle = 0 // force no cam rotation for UI

	for _, v := range containers {
		v.UpdateAndDraw(&gui.root, camera)
	}

	camera.Angle = prevAng // return camera stuff back to how it was
	camera.SetScreenArea(prevScX, prevScY, prevScW, prevScH)
}

// #region private

var updateAndDrawFuncs = map[string]func(cam *graphics.Camera, widget *widget, owner *container){
	"button": buttonUpdateAndDraw,
}

func dyn(cam *graphics.Camera, owner *container, value string) string {
	var tlx, tly = cam.PointFromPivot(0, 0)
	var brx, bry = cam.PointFromPivot(1, 1)
	var cx, cy = cam.PointFromPivot(0.5, 0.5)
	var w, h = cam.Size()

	value = strings.ReplaceAll(value, dynamic.CameraCenterX, symbols.New(cx))
	value = strings.ReplaceAll(value, dynamic.CameraCenterY, symbols.New(cy))
	value = strings.ReplaceAll(value, dynamic.CameraLeftX, symbols.New(tlx))
	value = strings.ReplaceAll(value, dynamic.CameraRightX, symbols.New(brx))
	value = strings.ReplaceAll(value, dynamic.CameraTopY, symbols.New(tly))
	value = strings.ReplaceAll(value, dynamic.CameraBottomY, symbols.New(bry))
	value = strings.ReplaceAll(value, dynamic.CameraWidth, symbols.New(w))
	value = strings.ReplaceAll(value, dynamic.CameraHeight, symbols.New(h))

	if owner != nil {
		var ox = symbols.New(dyn(cam, nil, owner.Properties[property.X]))
		var oy = symbols.New(dyn(cam, nil, owner.Properties[property.Y]))
		var ow = symbols.New(dyn(cam, nil, owner.Properties[property.Width]))
		var oh = symbols.New(dyn(cam, nil, owner.Properties[property.Height]))
		var olx = ox
		var orx = olx + "+" + ow
		var oty = oy
		var oby = oty + "+" + oh

		value = strings.ReplaceAll(value, dynamic.OwnerX, ox)
		value = strings.ReplaceAll(value, dynamic.OwnerY, oy)
		value = strings.ReplaceAll(value, dynamic.OwnerWidth, ow)
		value = strings.ReplaceAll(value, dynamic.OwnerHeight, oh)
		value = strings.ReplaceAll(value, dynamic.OwnerLeftX, olx)
		value = strings.ReplaceAll(value, dynamic.OwnerRightX, orx)
		value = strings.ReplaceAll(value, dynamic.OwnerTopY, oty)
		value = strings.ReplaceAll(value, dynamic.OwnerBottomY, oby)
	}

	return symbols.New(symbols.Calculate(value))
}

func getArea(cam *graphics.Camera, owner *container, props map[string]string) (x, y, w, h float32) {
	var vx, _ = strconv.ParseFloat(dyn(cam, owner, props[property.X]), 32)
	var vy, _ = strconv.ParseFloat(dyn(cam, owner, props[property.Y]), 32)
	var vw, _ = strconv.ParseFloat(dyn(cam, owner, props[property.Width]), 32)
	var vh, _ = strconv.ParseFloat(dyn(cam, owner, props[property.Height]), 32)
	return float32(vx), float32(vy), float32(vw), float32(vh)
}

func getColor(props map[string]string) uint {
	var rgba = strings.Split(props[property.RGBA], " ")
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
	return color.RGBA(byte(r), byte(g), byte(b), byte(a))
}

// #endregion
