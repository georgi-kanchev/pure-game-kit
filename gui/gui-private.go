package gui

import (
	"pure-game-kit/audio"
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
