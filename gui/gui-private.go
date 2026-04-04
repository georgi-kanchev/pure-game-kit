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

var sound audio.Audio = audio.New("")
var mouseX, mouseY, prevMouseX, prevMouseY float32
var updates = map[string]func(widget *widget){
	"button": button, "slider": slider, "checkbox": checkbox, "menu": menu, "inputField": inputField,
	"draggable": draggable,
}

var camCx, camCy, camLx, camRx, camTy, camBy, camW, camH float32                 // dynamic prop cache
var ownerLx, ownerRx, ownerTy, ownerBy, ownerCx, ownerCy, ownerW, ownerH float32 // dynamic prop cache
var tarLx, tarRx, tarTy, tarBy, tarCx, tarCy, tarW, tarH float32                 // dynamic prop cache
var tarHid, tarDis string                                                        // dynamic prop cache

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
var dynContainer *container // avoids closure allocation in dynNum

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
	var tx, ty, tw, th float32
	if targetContainer != nil {
		tx = dynNum(nil, r.themedField(field.X, targetContainer, nil), 0)
		ty = dynNum(nil, r.themedField(field.Y, targetContainer, nil), 0)
		tw = dynNum(nil, r.themedField(field.Width, targetContainer, nil), 0)
		th = dynNum(nil, r.themedField(field.Height, targetContainer, nil), 0)
		tarHid = targetContainer.Fields[field.Hidden]
		tarDis = targetContainer.Fields[field.Disabled]
	} else if targetWidget != nil {
		var owner = r.Containers[targetWidget.OwnerId]
		tx = dynNum(owner, r.themedField(field.X, owner, targetWidget), 0)
		ty = dynNum(owner, r.themedField(field.Y, owner, targetWidget), 0)
		tw = dynNum(owner, r.themedField(field.Width, owner, targetWidget), 0)
		th = dynNum(owner, r.themedField(field.Height, owner, targetWidget), 0)
		tarHid = targetWidget.Fields[field.Hidden]
		tarDis = targetWidget.Fields[field.Disabled]
	}

	tarLx, tarRx, tarTy, tarBy, tarW, tarH = tx, tx+tw, ty, ty+th, tw, th
	tarCx, tarCy = tx+tw/2, ty+th/2
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
	r.boxes = make([]*graphics.NinePatch, 0, 64)
	r.textBoxes = make([]*graphics.TextBox, 0, 64)
}
func (r *root) drawEnd() {
	r.cam.DrawNinePatches(r.boxes...)
	r.cam.DrawSprites(r.sprites...)
	r.cam.DrawTextBoxes(r.textBoxes...)
	r.cam.DrawSprites(r.spritesAbove...)
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
	camCx, camCy, camLx, camRx = cx, cy, tlx, brx
	camTy, camBy, camW, camH = tly, bry, w, h
}
func dynNum(c *container, value string, defaultValue float32) float32 {
	if value == "" {
		return defaultValue
	}
	dynContainer = c
	var calc = text.Calculate(value, func(name string) float32 {
		switch name {
		case dynamic.CameraCenterX:
			return camCx
		case dynamic.CameraCenterY:
			return camCy
		case dynamic.CameraLeftX:
			return camLx
		case dynamic.CameraRightX:
			return camRx
		case dynamic.CameraTopY:
			return camTy
		case dynamic.CameraBottomY:
			return camBy
		case dynamic.CameraWidth:
			return camW
		case dynamic.CameraHeight:
			return camH
		case dynamic.OwnerCenterX:
			if dynContainer != nil {
				return ownerCx
			}
		case dynamic.OwnerCenterY:
			if dynContainer != nil {
				return ownerCy
			}
		case dynamic.OwnerLeftX:
			if dynContainer != nil {
				return ownerLx
			}
		case dynamic.OwnerRightX:
			if dynContainer != nil {
				return ownerRx
			}
		case dynamic.OwnerTopY:
			if dynContainer != nil {
				return ownerTy
			}
		case dynamic.OwnerBottomY:
			if dynContainer != nil {
				return ownerBy
			}
		case dynamic.OwnerWidth:
			if dynContainer != nil {
				return ownerW
			}
		case dynamic.OwnerHeight:
			if dynContainer != nil {
				return ownerH
			}
		case dynamic.TargetCenterX:
			return tarCx
		case dynamic.TargetCenterY:
			return tarCy
		case dynamic.TargetLeftX:
			return tarLx
		case dynamic.TargetRightX:
			return tarRx
		case dynamic.TargetTopY:
			return tarTy
		case dynamic.TargetBottomY:
			return tarBy
		case dynamic.TargetWidth:
			return tarW
		case dynamic.TargetHeight:
			return tarH
		}
		return number.NaN()
	})
	if number.IsNaN(calc) {
		return defaultValue
	}
	return calc
}
func dyn(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	if text.Contains(value, "Target") {
		value = text.Replace(value, dynamic.TargetHidden, tarHid)
		value = text.Replace(value, dynamic.TargetDisabled, tarDis)
	}
	return value
}
func parseColor(value string, disabled ...bool) uint {
	if value == "" || value == "0 0 0 0" {
		return color.RGBA(0, 0, 0, 0)
	}

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
		a = 0
	}

	return color.RGBA(byte(r), byte(g), byte(b), byte(a))
}
func parseNum(value string, defaultValue float32) float32 {
	if value == "" {
		return defaultValue
	}
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
