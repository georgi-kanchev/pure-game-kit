package gui

import (
	"pure-game-kit/audio"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/dynamic"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

var sound *audio.Audio = audio.New("")
var mouseX, mouseY, prevMouseX, prevMouseY float32
var updateAndDrawFuncs = map[string]func(widget *widget){
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

var clickedId, clickedAndHeldId, sliderSlidId = "", "", ""

func (g *GUI) reset(inputState bool) (prAng, prZoom, prX, prY float32) {
	var cam = g.root.cam
	prAng, prZoom, prX, prY = cam.Angle, cam.Zoom, cam.X, cam.Y

	if inputState {
		mouseX, mouseY = cam.MousePosition()
		if mouse.IsButtonJustPressed(b.Left) {
			g.root.wPressedOn = nil
			tooltip = nil
			g.root.cPressedOnScrollH = nil
			g.root.cPressedOnScrollV = nil
		}
		if mouse.IsButtonJustReleased(b.Left) {
			g.root.cPressedOnScrollH = nil
			g.root.cPressedOnScrollV = nil
		}
		if mouse.IsButtonJustReleased(b.Middle) {
			g.root.cMiddlePressed = nil
		}
		if tooltip == nil {
			mouse.SetCursor(cursor.Arrow)
		}
	}

	cam.Zoom = g.Scale
	cam.Angle = 0       // force no cam rotation for UI
	cam.X, cam.Y = 0, 0 // force no position offset for UI
	return
}
func (r *root) themedField(fld string, c *container, w *widget) string {
	// priority for widget: widget -> widget theme -> container theme

	var widgetSelf, containerSelf, widgetThemeField, containerThemeField = "", "", "", ""
	var widgetTheme, containerTheme *theme
	var hasWidget, hasContainer, hasWidgetTheme, hasContainerTheme = false, false, false, false
	var hasWidgetThemeField, hasContainerThemeField = false, false

	if w != nil {
		widgetSelf, hasWidget = w.Fields[fld]
		widgetTheme, hasWidgetTheme = r.Themes[w.ThemeId]

		if hasWidgetTheme {
			widgetThemeField, hasWidgetThemeField = widgetTheme.Fields[fld]
		}
	}
	if c != nil {
		containerSelf, hasContainer = c.Fields[fld]
		containerTheme, hasContainerTheme = r.Themes[c.Fields[field.ThemeId]]

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
func (r *root) cacheDynTargetProps(targetId string) {
	var targetContainer = r.Containers[targetId]
	var targetWidget = r.Widgets[targetId]
	var tx, ty, tw, th, tHid, tDis string
	if targetContainer != nil {
		tx = r.themedField(field.X, targetContainer, nil)
		ty = r.themedField(field.Y, targetContainer, nil)
		tw = r.themedField(field.Width, targetContainer, nil)
		th = r.themedField(field.Height, targetContainer, nil)
		tHid = targetContainer.Fields[field.Hidden]
		tDis = targetContainer.Fields[field.Disabled]
	} else if targetWidget != nil {
		var owner = r.Containers[targetWidget.OwnerId]
		tx = r.themedField(field.X, owner, targetWidget)
		ty = r.themedField(field.Y, owner, targetWidget)
		tw = r.themedField(field.Width, owner, targetWidget)
		th = r.themedField(field.Height, owner, targetWidget)
		tHid = targetWidget.Fields[field.Hidden]
		tDis = targetWidget.Fields[field.Disabled]
	}

	tarLx, tarRx, tarTy, tarBy, tarW, tarH = tx, tx+"+"+tw, ty, ty+"+"+th, tw, th
	tarCx, tarCy = tx+"+"+tw+"/2", ty+"+"+th+"/2" // caching dynamic props
	tarHid, tarDis = tHid, tDis
}

func (r *root) restore(prAng, prZoom, prX, prY float32) {
	var cam = r.cam
	cam.Angle, cam.Zoom = prAng, prZoom
	cam.X, cam.Y = prX, prY

	r.wWasHovered = r.wHovered
	r.cWasHovered = r.cHovered
	prevMouseX, prevMouseY = mouseX, mouseY
}

func (r *root) drawStart() {
	r.sprites = make([]*graphics.Sprite, 0, 64)
	r.spritesAbove = make([]*graphics.Sprite, 0, 8)
	r.boxes = make([]*graphics.Box, 0, 64)
	r.textBoxes = make([]*graphics.TextBox, 0, 64)
}
func (r *root) drawEnd() {
	r.cam.DrawBoxes(r.boxes...)
	r.cam.DrawSprites(r.sprites...)
	r.cam.DrawTextBoxes(r.textBoxes...)
	r.cam.DrawSprites(r.spritesAbove...)
	if keyboard.IsKeyJustPressed(key.A) {
		print()
	}
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

func cacheDynCamProps(camera *graphics.Camera) {
	var tlx, tly = camera.PointFromEdge(0, 0)
	var brx, bry = camera.PointFromEdge(1, 1)
	var cx, cy = camera.PointFromEdge(0.5, 0.5)
	var w, h = camera.Size()
	camCx, camCy, camLx, camRx = text.New(cx), text.New(cy), text.New(tlx), text.New(brx)
	camTy, camBy, camW, camH = text.New(tly), text.New(bry), text.New(w), text.New(h)
}

func dyn(c *container, value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	if text.Contains(value, "Target") {
		value = text.Replace(value, dynamic.TargetWidth, tarW)
		value = text.Replace(value, dynamic.TargetHeight, tarH)
		value = text.Replace(value, dynamic.TargetLeftX, tarLx)
		value = text.Replace(value, dynamic.TargetRightX, tarRx)
		value = text.Replace(value, dynamic.TargetTopY, tarTy)
		value = text.Replace(value, dynamic.TargetBottomY, tarBy)
		value = text.Replace(value, dynamic.TargetCenterX, tarCx)
		value = text.Replace(value, dynamic.TargetCenterY, tarCy)
		value = text.Replace(value, dynamic.TargetHidden, tarHid)
		value = text.Replace(value, dynamic.TargetDisabled, tarDis)
	}
	if c != nil && text.Contains(value, "Owner") {
		value = text.Replace(value, dynamic.OwnerWidth, ownerW)
		value = text.Replace(value, dynamic.OwnerHeight, ownerH)
		value = text.Replace(value, dynamic.OwnerLeftX, ownerLx)
		value = text.Replace(value, dynamic.OwnerRightX, ownerRx)
		value = text.Replace(value, dynamic.OwnerTopY, ownerTy)
		value = text.Replace(value, dynamic.OwnerBottomY, ownerBy)
		value = text.Replace(value, dynamic.OwnerCenterX, ownerCx)
		value = text.Replace(value, dynamic.OwnerCenterY, ownerCy)
	}
	if text.Contains(value, "Camera") {
		value = text.Replace(value, dynamic.CameraCenterX, camCx)
		value = text.Replace(value, dynamic.CameraCenterY, camCy)
		value = text.Replace(value, dynamic.CameraLeftX, camLx)
		value = text.Replace(value, dynamic.CameraRightX, camRx)
		value = text.Replace(value, dynamic.CameraTopY, camTy)
		value = text.Replace(value, dynamic.CameraBottomY, camBy)
		value = text.Replace(value, dynamic.CameraWidth, camW)
		value = text.Replace(value, dynamic.CameraHeight, camH)
	}

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
	var mx, my = cam.MousePosition()
	var result = mx > x && mx < x+w && my > y && my < y+h
	cam.Angle = prevAng
	return result
}
