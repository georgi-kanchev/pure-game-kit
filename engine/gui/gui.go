package gui

import (
	"encoding/xml"
	"fmt"
	"math"
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
	var root root
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

	var err = xml.Unmarshal([]byte(result), &root)
	fmt.Printf("err: %v\n", err)

	root.Containers = map[string]container{}
	root.Widgets = map[string]widget{}
	root.Themes = map[string]theme{}

	for i := range root.XmlContainers {
		var c = &root.XmlContainers[i]
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
			w.Properties = make(map[string]string, len(w.XmlProps))

			if has {
				w.UpdateAndDraw = fn
			}

			for _, xmlProp := range w.XmlProps {
				w.Properties[xmlProp.Name.Local] = xmlProp.Value
			}

			root.Widgets[wId] = *w
		}
		for j := range c.XmlThemes {
			var t = &c.XmlThemes[j]
			var tId = t.XmlProps[0].Value
			t.Properties = make(map[string]string, len(t.XmlProps))

			for _, xmlProp := range t.XmlProps {
				t.Properties[xmlProp.Name.Local] = xmlProp.Value
			}
			root.Themes[tId] = *t
		}

		root.Containers[cId] = *c
	}

	return GUI{root: root}
}

func (gui *GUI) ThemeProperty(themeId, property string) string {
	return gui.root.Themes[themeId].Properties[property]
}
func (gui *GUI) SetThemeProperty(themeId, property, value string) {
	gui.root.Themes[themeId].Properties[property] = value
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

	camera.Angle = 0 // force no cam rotation for UI

	for _, c := range containers {
		c.UpdateAndDraw(&gui.root, camera)
	}

	camera.Angle = prevAng // reset angle & mask to how it was
	camera.SetScreenArea(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)
}

// #region private

var updateAndDrawFuncs = map[string]func(cam *graphics.Camera, root *root, widget *widget, owner *container){
	"button": buttonUpdateAndDraw,
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
	// priority for widget: widget -> widget theme -> container theme -> container
	// priority for container: container -> container theme

	var wSelf, cSelf = "", ""
	var wTheme, cTheme theme
	var hasW, hasC, hasWt, hasCt = false, false, false, false

	if w != nil {
		wSelf, hasW = w.Properties[prop]
		wTheme, hasWt = root.Themes[w.Properties[property.ThemeId]]
	}
	if c != nil {
		cSelf, hasC = c.Properties[prop]
		cTheme, hasCt = root.Themes[c.Properties[property.ThemeId]]
	}

	// container checks
	if w == nil {
		if hasC {
			return cSelf
		}
		if hasCt {
			return cTheme.Properties[prop]
		}

		return ""
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
	if hasC {
		return cSelf
	}

	return ""
}

func dyn(cam *graphics.Camera, owner *container, value string, defaultValue string) string {
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
		var ox = symbols.New(dyn(cam, nil, owner.Properties[property.X], "0"))
		var oy = symbols.New(dyn(cam, nil, owner.Properties[property.Y], "0"))
		var ow = symbols.New(dyn(cam, nil, owner.Properties[property.Width], "0"))
		var oh = symbols.New(dyn(cam, nil, owner.Properties[property.Height], "0"))
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

	var calc = symbols.Calculate(value)
	if math.IsNaN(calc) {
		return defaultValue
	}
	return symbols.New(calc)
}

func parseColor(value string) uint {
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
	return color.RGBA(byte(r), byte(g), byte(b), byte(a))
}
func parseNum(value string, defaultValue float32) float32 {
	var v, err = strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}
	return float32(v)
}

// #endregion
