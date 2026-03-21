package gui

import (
	"pure-game-kit/execution/condition"
	f "pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
)

// https://showcase.primefaces.org - basic default browser widgets showcase (scroll down to forms on the left)

type GUI struct {
	Scale  float32
	Volume float32
	root   *root
}

// =================================================================

func (g *GUI) UpdateAndDraw() {
	var cam = g.root.cam
	var containers = g.root.ContainerIds
	var prAng, prZoom, prX, prY = g.reset(true) // keep order of variables & reset
	cacheDynCamProps(cam)
	g.root.Volume = g.Volume

	sliderSlidId = condition.If(sliderSlidId != "", "", sliderSlidId)

	var prevMask = g.root.cam.Mask
	g.root.cam.Mask = nil

	for _, id := range containers {
		var c = g.root.Containers[id]
		var _, hasTarget = c.Fields[f.TargetId]
		if hasTarget {
			g.root.cacheDynTargetProps(g.root.themedField(f.TargetId, c, nil))
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

		c.updateAndDraw()
	}

	if g.root.cWasHovered == g.root.cHovered {
		g.root.cFocused = g.root.cHovered // only containers hovered 2 frames in a row get input (top-down prio)
	}
	if g.root.wWasHovered == g.root.wHovered {
		g.root.wFocused = g.root.wHovered // only widgets hovered 2 frames in a row get input (top-down prio)
	}

	if g.root.wPressedOn != nil && g.root.wPressedOn.Class == "draggable" {
		cam.Mask = nil
		drawDraggable(g.root.wPressedOn)
	}
	if tooltip != nil {
		cam.Mask = nil
		drawTooltip(g.root.Containers[tooltip.OwnerId])
	}

	clickedId = condition.If(clickedId != "", "", clickedId)
	clickedAndHeldId = condition.If(clickedAndHeldId != "", "", clickedAndHeldId)

	if g.root.wPressedOn != nil {
		if g.root.IsButtonJustClicked(g.root.wPressedOn.Id) {
			clickedId = g.root.wPressedOn.Id
		}
		if g.root.IsButtonClickedAndHeld(g.root.wPressedOn.Id) {
			clickedAndHeldId = g.root.wPressedOn.Id
		}
	}

	if mouse.IsButtonJustReleased(b.Left) {
		g.root.wPressedOn = nil
		tooltip = nil
	}

	g.root.restore(prAng, prZoom, prX, prY) // undo what reset does, everything as it was for cam
	g.root.cam.Mask = prevMask
}

// Works for Widgets & Containers.
func (g *GUI) SetField(anyId, field string, value string) {
	var w, hasW = g.root.Widgets[anyId]
	var c, hasC = g.root.Containers[anyId]
	var t, hasT = g.root.Themes[anyId]

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

// Works for Widgets & Containers. Use
//
//	FieldNumber(...)
//
// for dynamic values.
func (g *GUI) Field(anyId, field string) string {
	var w, hasW = g.root.Widgets[anyId]
	var c, hasC = g.root.Containers[anyId]
	var t, hasT = g.root.Themes[anyId]

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

// Works for Widgets & Containers. Converts the appropriate fields to numbers while replacing their dynamic parts.
func (g *GUI) FieldNumber(anyId, field string) float32 {
	var w, hasW = g.root.Widgets[anyId]
	var owner *container
	if hasW {
		owner = g.root.Containers[w.OwnerId]
	}
	var value = dyn(owner, g.Field(anyId, field), "NaN")
	return parseNum(value, number.NaN())
}

func (g *GUI) AreaText(widgetId string) (width, height float32) {
	var w, hasW = g.root.Widgets[widgetId]

	if hasW && w.textBox != nil {
		var t = w.textBox
		var text = condition.If(t.WordWrap, t.TextWrap(t.Text), t.Text)
		width, height = t.TextMeasure(text)
		return width, height
	}
	return number.NaN(), number.NaN()
}

// Works for Widgets & Containers.
func (g *GUI) Area(anyId string) (x, y, width, height, angle float32) {
	var cam = g.root.cam
	var zoom = cam.Zoom / g.Scale
	var w, hasW = g.root.Widgets[anyId]
	var c, hasC = g.root.Containers[anyId]
	if hasC {
		x, y = c.X, c.Y
		width, height = c.Width, c.Height
	}

	if hasW {
		x, y = w.X, w.Y
		width, height = w.Width, w.Height
	}
	angle = -cam.Angle
	x, y = cam.X+x/zoom, cam.Y+y/zoom
	x, y = point.RotateAroundPoint(x, y, cam.X, cam.Y, angle)
	width, height = width/zoom, height/zoom
	return
}

func (g *GUI) IsAnyHovered() bool {
	var prAng, prZoom, prX, prY = g.reset(false)
	defer func() { g.root.restore(prAng, prZoom, prX, prY) }()

	for _, c := range g.root.Containers {
		var hidden = c.Fields[f.Hidden]
		if hidden == "" && c.isHovered() {
			return true
		}
	}
	return false
}

// Works for Widgets & Containers.
func (g *GUI) IsHovered(anyId string) bool {
	var prAng, prZoom, prX, prY = g.reset(false)
	defer func() { g.root.restore(prAng, prZoom, prX, prY) }()

	var w, hasW = g.root.Widgets[anyId]
	var c, hasC = g.root.Containers[anyId]

	if hasW {
		return w.isFocused()
	}
	if hasC {
		return c.isFocused()
	}
	return false
}

// Works for Widgets & Containers.
func (g *GUI) IsFocused(widgetId string) bool {
	var prAng, prZoom, prX, prY = g.reset(false)
	defer func() { g.root.restore(prAng, prZoom, prX, prY) }()
	var w, has = g.root.Widgets[widgetId]
	if has {
		return w.isFocused()
	}
	return false
}

func (g *GUI) WidgetIdsOfContainer(containerId string) []string {
	var c = g.root.Containers[containerId]
	if c == nil {
		return nil
	}

	return collection.Clone(c.Widgets)
}
func (g *GUI) AllWidgetIds() []string {
	return collection.MapKeys(g.root.Widgets)
}
func (g *GUI) AllContainerIds() []string {
	return collection.MapKeys(g.root.Containers)
}
func (g *GUI) AllThemeIds() []string {
	return collection.MapKeys(g.root.Themes)
}
