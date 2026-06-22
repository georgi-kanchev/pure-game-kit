package assets

import (
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/storage"
	"pure-game-kit/packages/utility/text"
)

type LayoutId uint32

func LoadLayout(xmlPath string) LayoutId {
	var layout = internal.Layout{}
	storage.FromXML(file.LoadText(xmlPath), &layout)

	if len(layout.Boxes) == 0 {
		return 0
	}

	internal.NextLayoutId++
	var id = internal.NextLayoutId

	for i := range layout.Items { // pre-calculate item range indexes for each box
		var b = int(layout.Items[i].BoxId)
		if !layout.Boxes[b].ItemRangeCalculated {
			layout.Boxes[b].ItemStart = len(layout.Items)
			layout.Boxes[b].ItemEnd = 0
			layout.Boxes[b].ItemRangeCalculated = true
		}
		if i < layout.Boxes[b].ItemStart {
			layout.Boxes[b].ItemStart = i
		}
		if i+1 > layout.Boxes[b].ItemEnd {
			layout.Boxes[b].ItemEnd = i + 1
		}
	}

	internal.Layouts[id] = &layout
	return LayoutId(id)
}

func (l LayoutId) Unload() {
	delete(internal.Layouts, uint32(l))
}

func (l LayoutId) Box(id int) (area geometry.Area) {
	var layout = internal.Layouts[uint32(l)]
	if layout == nil {
		return geometry.Area{}
	}
	var rx, ry, rw, rh, rv = boxDynamic(layout, id, 0)
	if !rv {
		return geometry.Area{} // not visible
	}
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	area = geometry.NewArea((rx+rw/2)*sc, (ry+rh/2)*sc, rw*sc, rh*sc)
	return area
}
func (l LayoutId) Item(id int, scrollX, scrollY float32) (area, mask geometry.Area) {
	var layout = internal.Layouts[uint32(l)]
	if layout == nil || id < 0 || id >= len(layout.Items) {
		return geometry.Area{}, geometry.Area{}
	}
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	var rx, ry, rw, rh, rv = itemDynamic(layout, id, scrollX, scrollY, sc)
	if !rv {
		return geometry.Area{}, geometry.Area{} // not visible
	}
	var ownerId = layout.Items[id].BoxId
	var o = l.Box(int(ownerId))
	area = geometry.NewArea((rx+rw/2)*sc, (ry+rh/2)*sc, rw*sc, rh*sc)
	mask = geometry.NewArea(o.X, o.Y, o.Width, o.Height)
	return area, mask
}

func (l LayoutId) SetVisibleItem(id int, visible bool) {
	var layout = internal.Layouts[uint32(l)]
	if layout == nil || id < 0 || id >= len(layout.Items) {
		return
	}
	var item = &layout.Items[id]
	if visible {
		item.Visible = 1
	} else {
		item.Visible = 0
	}
}

// private ========================================================

var activeVars *internal.Vars

func boxDynamic(layout *internal.Layout, boxId int, depth int) (x, y, w, h float32, vis bool) {
	if depth > 8 {
		return
	}

	var box = &layout.Boxes[boxId]
	var ew = 512 * number.SquareRoot(internal.WindowWidth/internal.WindowHeight)
	var eh = 512 / number.SquareRoot(internal.WindowWidth/internal.WindowHeight)

	box.Vars = internal.Vars{}
	box.Vars.Mx = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 0))
	box.Vars.My = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 1))
	box.Vars.Mw = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 2))
	box.Vars.Mh = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 3))
	box.Vars.Mlx, box.Vars.Mly = box.Vars.Mx, box.Vars.My+box.Vars.Mh/2
	box.Vars.Mrx, box.Vars.Mry = box.Vars.Mx+box.Vars.Mw, box.Vars.Mly
	box.Vars.Mux, box.Vars.Muy = box.Vars.Mx+box.Vars.Mw/2, box.Vars.My
	box.Vars.Mdx, box.Vars.Mdy = box.Vars.Mux, box.Vars.My+box.Vars.Mh

	box.Vars.Sx, box.Vars.Sy, box.Vars.Sw, box.Vars.Sh = -ew/2, -eh/2, ew, eh
	box.Vars.Slx, box.Vars.Sly = box.Vars.Sx, box.Vars.Sy+box.Vars.Sh/2
	box.Vars.Srx, box.Vars.Sry = box.Vars.Sx+box.Vars.Sw, box.Vars.Sly
	box.Vars.Sux, box.Vars.Suy = box.Vars.Sx+box.Vars.Sw/2, box.Vars.Sy
	box.Vars.Sdx, box.Vars.Sdy = box.Vars.Sux, box.Vars.Sy+box.Vars.Sh

	box.Vars.Tx, box.Vars.Ty, box.Vars.Tw, box.Vars.Th = 0, 0, 0, 0
	box.Vars.Tlx, box.Vars.Tly, box.Vars.Trx, box.Vars.Try = 0, 0, 0, 0
	box.Vars.Tux, box.Vars.Tuy, box.Vars.Tdx, box.Vars.Tdy = 0, 0, 0, 0

	activeVars = &box.Vars
	if text.SplitCount(box.Targets, " ") == 4 {
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 0), depth+1)
		activeVars = &box.Vars
		var rx = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variable)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 1), depth+1)
		activeVars = &box.Vars
		var ry = text.Calculate(text.SplitIndex(box.Expression, " ", 1), variable)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 2), depth+1)
		activeVars = &box.Vars
		var rw = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variable)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 3), depth+1)
		activeVars = &box.Vars
		var rh = text.Calculate(text.SplitIndex(box.Expression, " ", 3), variable)
		return rx, ry, rw, rh, box.Visible == 1
	}

	var rx = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variable)
	var ry = text.Calculate(text.SplitIndex(box.Expression, " ", 1), variable)
	var rw = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variable)
	var rh = text.Calculate(text.SplitIndex(box.Expression, " ", 3), variable)
	return rx, ry, rw, rh, box.Visible == 1
}
func itemDynamic(layout *internal.Layout, itemId int, scrollX, scrollY, sc float32) (x, y, w, h float32, vis bool) {
	var item = &layout.Items[itemId]
	var box = &layout.Boxes[item.BoxId]
	if box.Visible == 0 {
		return 0, 0, 0, 0, false
	}
	var bx, by, bw, bh, _ = boxDynamic(layout, int(item.BoxId), 0)

	item.Vars = internal.Vars{}
	item.Vars.Ow, item.Vars.Oh, item.Vars.Ov = bw, bh, 1

	var osx, osy float32
	if text.SplitCount(box.ItemSpacing, " ") >= 2 {
		osx = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 0))
		osy = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 1))
	}
	var og, mb = box.ItemGap, box.ItemNewRow
	item.Vars.Osx, item.Vars.Osy, item.Vars.Og, item.Vars.Mnr = osx, osy, og, mb

	activeVars = &item.Vars
	if text.SplitCount(box.ItemSize, " ") >= 2 {
		item.Vars.Mw = text.Calculate(text.SplitIndex(box.ItemSize, " ", 0), variable)
		item.Vars.Mh = text.Calculate(text.SplitIndex(box.ItemSize, " ", 1), variable)
	} else if text.SplitCount(item.Size, " ") >= 2 {
		item.Vars.Mw = text.ToNumber[float32](text.SplitIndex(item.Size, " ", 0))
		item.Vars.Mh = text.ToNumber[float32](text.SplitIndex(item.Size, " ", 1))
	} else {
		item.Vars.Mw, item.Vars.Mh = 40, 20
	}
	var defW, defH = item.Vars.Mw, item.Vars.Mh
	var alignX, alignY float32
	if text.SplitCount(box.ItemAlign, " ") >= 2 {
		alignX = text.ToNumber[float32](text.SplitIndex(box.ItemAlign, " ", 0))
		alignY = text.ToNumber[float32](text.SplitIndex(box.ItemAlign, " ", 1))
	}

	var curX, curY, maxX, maxY = bx + osx, by + osy, bx, by
	var rowMaxH, targetMX, targetMY float32
	for i := box.ItemStart; i < box.ItemEnd; i++ {
		var it = &layout.Items[i]
		if it.Visible == 0 {
			continue
		}

		if it.NewRow == 1 {
			curY += rowMaxH + mb
			if it.NewRowExpression != "" {
				item.Vars.Mx, item.Vars.My, item.Vars.Mw, item.Vars.Mh = curX, curY, defW, defH
				curY += text.Calculate(it.NewRowExpression, variable)
			}
			curX = bx + osx
			rowMaxH = 0
		}

		item.Vars.Mx, item.Vars.My, item.Vars.Mw, item.Vars.Mh = curX, curY, defW, defH
		var itW, itH = defW, defH
		if text.SplitCount(it.Expression, " ") == 4 {
			itW = text.Calculate(text.SplitIndex(it.Expression, " ", 2), variable)
			itH = text.Calculate(text.SplitIndex(it.Expression, " ", 3), variable)
		}

		if i == itemId {
			targetMX, targetMY = curX, curY
		}
		if curX+itW > maxX {
			maxX = curX + itW
		}
		if curY+itH > maxY {
			maxY = curY + itH
		}

		curX += itW + og
		if itH > rowMaxH {
			rowMaxH = itH
		}
	}

	item.Vars.Mx, item.Vars.My = targetMX+(bx+bw-maxX)*alignX, targetMY+(by+bh-maxY)*alignY
	item.Vars.Mw, item.Vars.Mh = defW, defH
	item.Vars.Mx, item.Vars.My = item.Vars.Mx-(max(0, maxX-(bx+bw))*scrollX), item.Vars.My-(max(0, maxY-(by+bh))*scrollY)

	box.ContentWidth, box.ContentHeight = (maxX-bx)*sc-1, (maxY-by)*sc-1

	var rx = text.Calculate(text.SplitIndex(item.Expression, " ", 0), variable)
	var ry = text.Calculate(text.SplitIndex(item.Expression, " ", 1), variable)
	var rw = text.Calculate(text.SplitIndex(item.Expression, " ", 2), variable)
	var rh = text.Calculate(text.SplitIndex(item.Expression, " ", 3), variable)
	return rx, ry, rw, rh, item.Visible == 1
}

func setTargetVars(layout *internal.Layout, vars *internal.Vars, tar string, depth int) {
	if tar != "" {
		var targetId = text.ToNumber[int](tar)
		if targetId >= 0 && targetId < len(layout.Boxes) {
			vars.Tx, vars.Ty, vars.Tw, vars.Th, _ = boxDynamic(layout, targetId, depth)
		} else {
			vars.Tx, vars.Ty, vars.Tw, vars.Th = 0, 0, 0, 0
		}
	} else {
		vars.Tx, vars.Ty, vars.Tw, vars.Th = 0, 0, 0, 0
	}
	vars.Tlx, vars.Tly, vars.Trx, vars.Try = vars.Tx, vars.Ty+vars.Th/2, vars.Tx+vars.Tw, vars.Tly
	vars.Tux, vars.Tuy, vars.Tdx, vars.Tdy = vars.Tx+vars.Tw/2, vars.Ty, vars.Tux, vars.Ty+vars.Th
}
func variable(name string) float32 {
	switch name {
	case "mx":
		return activeVars.Mx
	case "my":
		return activeVars.My
	case "mw":
		return activeVars.Mw
	case "mh":
		return activeVars.Mh
	case "mlx":
		return activeVars.Mlx
	case "mly":
		return activeVars.Mly
	case "mrx":
		return activeVars.Mrx
	case "mry":
		return activeVars.Mry
	case "mux":
		return activeVars.Mux
	case "muy":
		return activeVars.Muy
	case "mdx":
		return activeVars.Mdx
	case "mdy":
		return activeVars.Mdy
	case "sx":
		return activeVars.Sx
	case "sy":
		return activeVars.Sy
	case "sw":
		return activeVars.Sw
	case "sh":
		return activeVars.Sh
	case "slx":
		return activeVars.Slx
	case "sly":
		return activeVars.Sly
	case "srx":
		return activeVars.Srx
	case "sry":
		return activeVars.Sry
	case "sux":
		return activeVars.Sux
	case "suy":
		return activeVars.Suy
	case "sdx":
		return activeVars.Sdx
	case "sdy":
		return activeVars.Sdy
	case "tx":
		return activeVars.Tx
	case "ty":
		return activeVars.Ty
	case "tw":
		return activeVars.Tw
	case "th":
		return activeVars.Th
	case "tlx":
		return activeVars.Tlx
	case "tly":
		return activeVars.Tly
	case "trx":
		return activeVars.Trx
	case "try":
		return activeVars.Try
	case "tux":
		return activeVars.Tux
	case "tuy":
		return activeVars.Tuy
	case "tdx":
		return activeVars.Tdx
	case "tdy":
		return activeVars.Tdy
	case "ow":
		return activeVars.Ow
	case "oh":
		return activeVars.Oh
	case "ov":
		return activeVars.Ov
	case "osx":
		return activeVars.Osx
	case "osy":
		return activeVars.Osy
	case "og":
		return activeVars.Og
	case "mnr":
		return activeVars.Mnr
	}
	return 0
}
