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
	"strings"
)

var sound *audio.Audio = audio.New("")
var mouseX, mouseY, prevMouseX, prevMouseY float32
var updateAndDrawFuncs = map[string]func(cam *graphics.Camera, root *root, widget *widget){
	"button": button, "slider": slider, "checkbox": checkbox, "menu": menu, "inputField": inputField,
	"draggable": draggable,
}
var camCx, camCy, camLx, camRx, camTy, camBy, camW, camH string                 // dynamic prop cache
var ownerLx, ownerRx, ownerTy, ownerBy, ownerCx, ownerCy, ownerW, ownerH string // dynamic prop cache
var tarLx, tarRx, tarTy, tarBy, tarCx, tarCy, tarW, tarH, tarHid, tarDis string // dynamic prop cache

var buttonColor uint
var btnSounds = true

var typingIn *widget
var indexCursor, indexSelect int
var cursorTime, scrollX, textMargin float32
var symbolXs = []float32{}
var maskText bool       // used for inputbox mask
var simulateRemove bool // used to delete text when typing
var frame int

var tooltip, tooltipForWidget *widget
var tooltipAt float32
var tooltipVisible, tooltipWasVisible bool

// var textBox graphics.TextBox = graphics.TextBox{}
var sprite graphics.Sprite = graphics.Sprite{}
var box graphics.Box = graphics.Box{}

var reusableWidget = &widget{Fields: map[string]string{}}

var clickedId, clickedAndHeldId = "", ""

func (gui *GUI) reset(camera *graphics.Camera) {
	if mouse.IsButtonJustPressed(b.Left) {
		gui.root.wPressedOn = nil
		tooltip = nil
		gui.root.cPressedOnScrollH = nil
		gui.root.cPressedOnScrollV = nil
	}
	if mouse.IsButtonJustReleased(b.Left) {
		gui.root.cPressedOnScrollH = nil
		gui.root.cPressedOnScrollV = nil
	}
	if mouse.IsButtonJustReleased(b.Middle) {
		gui.root.cMiddlePressed = nil
	}

	camera.Zoom = gui.Scale
	camera.Angle = 0          // force no cam rotation for UI
	camera.X, camera.Y = 0, 0 // force no position offset for UI
	mouseX, mouseY = camera.MousePosition()

	if tooltip == nil {
		mouse.SetCursor(cursor.Arrow)
	}
}
func (root *root) themedField(fld string, c *container, w *widget) string {
	// priority for widget: widget -> widget theme -> container theme

	var widgetSelf, containerSelf, widgetThemeField, containerThemeField = "", "", "", ""
	var widgetTheme, containerTheme *theme
	var hasWidget, hasContainer, hasWidgetTheme, hasContainerTheme = false, false, false, false
	var hasWidgetThemeField, hasContainerThemeField = false, false

	if w != nil {
		widgetSelf, hasWidget = w.Fields[fld]
		widgetTheme, hasWidgetTheme = root.Themes[w.ThemeId]

		if hasWidgetTheme {
			widgetThemeField, hasWidgetThemeField = widgetTheme.Fields[fld]
		}
	}
	if c != nil {
		containerSelf, hasContainer = c.Fields[fld]
		containerTheme, hasContainerTheme = root.Themes[c.Fields[field.ThemeId]]

		if hasContainerTheme {
			containerThemeField, hasContainerThemeField = containerTheme.Fields[fld]
		}
	}

	if w != nil {
		if hasWidget {
			return widgetSelf
		}
		if hasWidgetTheme && hasWidgetThemeField {
			return widgetThemeField
		}
		if hasContainerTheme && hasContainerThemeField {
			return containerThemeField
		}
		if hasContainer {
			return containerSelf
		}
	}

	if hasContainer {
		return containerSelf
	}
	if hasContainerTheme {
		return containerTheme.Fields[fld]
	}

	return ""
}
func (root *root) cacheTarget(targetId string) {
	var targetContainer = root.Containers[targetId]
	var targetWidget = root.Widgets[targetId]
	var tx, ty, tw, th, tHid, tDis string
	if targetContainer != nil {
		tx = root.themedField(field.X, targetContainer, nil)
		ty = root.themedField(field.Y, targetContainer, nil)
		tw = root.themedField(field.Width, targetContainer, nil)
		th = root.themedField(field.Height, targetContainer, nil)
		tHid = targetContainer.Fields[field.Hidden]
		tDis = targetContainer.Fields[field.Disabled]
	} else if targetWidget != nil {
		var owner = root.Containers[targetWidget.OwnerId]
		tx = root.themedField(field.X, owner, targetWidget)
		ty = root.themedField(field.Y, owner, targetWidget)
		tw = root.themedField(field.Width, owner, targetWidget)
		th = root.themedField(field.Height, owner, targetWidget)
		tHid = targetWidget.Fields[field.Hidden]
		tDis = targetWidget.Fields[field.Disabled]
	}

	tarLx, tarRx, tarTy, tarBy, tarW, tarH = tx, tx+"+"+tw, ty, ty+"+"+th, tw, th
	tarCx, tarCy = tx+"+"+tw+"/2", ty+"+"+th+"/2" // caching dynamic props
	tarHid, tarDis = tHid, tDis
}

func (root *root) restore(camera *graphics.Camera, prevAng, prevZoom, prevX, prevY float32) {
	camera.Angle, camera.Zoom = prevAng, prevZoom // reset angle, zoom & mask to how it was
	camera.X, camera.Y = prevX, prevY             // also x y
	camera.SetScreenArea(camera.ScreenX, camera.ScreenY, camera.ScreenWidth, camera.ScreenHeight)

	if mouse.IsButtonJustReleased(b.Left) {
		root.wPressedOn = nil
		tooltip = nil
	}

	root.wWasHovered = root.wHovered
	root.cWasHovered = root.cHovered
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
func defaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func dyn(owner *container, value string, defaultValue string) string {
	value = strings.ReplaceAll(value, dynamic.TargetWidth, tarW)
	value = strings.ReplaceAll(value, dynamic.TargetHeight, tarH)
	value = strings.ReplaceAll(value, dynamic.TargetLeftX, tarLx)
	value = strings.ReplaceAll(value, dynamic.TargetRightX, tarRx)
	value = strings.ReplaceAll(value, dynamic.TargetTopY, tarTy)
	value = strings.ReplaceAll(value, dynamic.TargetBottomY, tarBy)
	value = strings.ReplaceAll(value, dynamic.TargetCenterX, tarCx)
	value = strings.ReplaceAll(value, dynamic.TargetCenterY, tarCy)
	value = strings.ReplaceAll(value, dynamic.TargetHidden, tarHid)
	value = strings.ReplaceAll(value, dynamic.TargetDisabled, tarDis)

	if owner != nil {
		value = text.Replace(value, dynamic.OwnerWidth, ownerW)
		value = text.Replace(value, dynamic.OwnerHeight, ownerH)
		value = text.Replace(value, dynamic.OwnerLeftX, ownerLx)
		value = text.Replace(value, dynamic.OwnerRightX, ownerRx)
		value = text.Replace(value, dynamic.OwnerTopY, ownerTy)
		value = text.Replace(value, dynamic.OwnerBottomY, ownerBy)
		value = text.Replace(value, dynamic.OwnerCenterX, ownerCx)
		value = text.Replace(value, dynamic.OwnerCenterY, ownerCy)
	}

	value = text.Replace(value, dynamic.CameraCenterX, camCx)
	value = text.Replace(value, dynamic.CameraCenterY, camCy)
	value = text.Replace(value, dynamic.CameraLeftX, camLx)
	value = text.Replace(value, dynamic.CameraRightX, camRx)
	value = text.Replace(value, dynamic.CameraTopY, camTy)
	value = text.Replace(value, dynamic.CameraBottomY, camBy)
	value = text.Replace(value, dynamic.CameraWidth, camW)
	value = text.Replace(value, dynamic.CameraHeight, camH)

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
