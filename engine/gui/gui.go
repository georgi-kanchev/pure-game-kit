package gui

import (
	"encoding/xml"
	"fmt"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/dynamic"
	"pure-kit/engine/gui/property"
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

	// container is missing, add root container
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

	for i := range gui.Containers {
		var c = &gui.Containers[i]

		for j := range c.Widgets {
			var w = &c.Widgets[j]
			var class = w.findPropValue("class", "")
			var fn, has = updateAndDrawFuncs[class]

			if has {
				w.UpdateAndDraw = fn
			}
		}
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
	var _, widget = gui.root.findWidget(widgetId)
	if widget == nil {
		return ""
	}

	return widget.findPropValue(property, "")
}
func (gui *GUI) SetProperty(widgetId, property string, value string) {
	var _, widget = gui.root.findWidget(widgetId)
	if widget == nil {
		return
	}

	var prop = widget.findProp(property)
	if prop != nil {
		prop.Value = value
	}
}

func (gui *GUI) Draw(camera *graphics.Camera) {
	var prevAng = camera.Angle
	var containers = gui.root.Containers
	var prevScX, prevScY = camera.ScreenX, camera.ScreenY // each container masks an area, store to turn back later
	var prevScW, prevScH = camera.ScreenWidth, camera.ScreenHeight

	camera.Angle = 0 // force no cam rotation for UI

	for i := range containers {
		var c = &containers[i]
		var x, _ = strconv.ParseFloat(dyn(camera, c, c.findPropValue(property.X, "0")), 32)
		var y, _ = strconv.ParseFloat(dyn(camera, c, c.findPropValue(property.Y, "0")), 32)
		var w, _ = strconv.ParseFloat(dyn(camera, c, c.findPropValue(property.Width, "0")), 32)
		var h, _ = strconv.ParseFloat(dyn(camera, c, c.findPropValue(property.Height, "0")), 32)
		var scx, scy = camera.PointToScreen(float32(x), float32(y))

		camera.Mask(scx, scy, int(w), int(h))
		c.draw(camera, nil)
		for j := range c.Widgets {
			var widget = &c.Widgets[j]
			widget.UpdateAndDraw(camera, widget, c)
		}
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
		var ox = symbols.New(dyn(cam, nil, owner.findPropValue(property.X, "0")))
		var oy = symbols.New(dyn(cam, nil, owner.findPropValue(property.Y, "0")))
		var ow = symbols.New(dyn(cam, nil, owner.findPropValue(property.Width, "0")))
		var oh = symbols.New(dyn(cam, nil, owner.findPropValue(property.Height, "0")))
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

// #endregion
