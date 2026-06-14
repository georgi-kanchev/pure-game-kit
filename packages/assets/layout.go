package assets

import (
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
	internal.Layouts[id] = layout
	return LayoutId(id)
}

func (l LayoutId) Unload() {
	delete(internal.Layouts, uint32(l))
}

func (l LayoutId) BoxArea(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Boxes) {
		return
	}
	var rx, ry, rw, rh = boxDynamic(&layout, id, 0)
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	return (rx + rw/2) * sc, (ry + rh/2) * sc, rw * sc, rh * sc
}
func (l LayoutId) ItemArea(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var rx, ry, rw, rh = itemDynamic(&layout, id)
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	return (rx + rw/2) * sc, (ry + rh/2) * sc, rw * sc, rh * sc
}
func (l LayoutId) ItemMask(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var ownerId = layout.Items[id].BoxId
	var ox, oy, ow, oh = l.BoxArea(int(ownerId), zoom)
	return ox - ow/2, oy - oh/2, ow, oh
}

// private ========================================================

func boxDynamic(layout *internal.Layout, boxId int, depth int) (x, y, w, h float32) {
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

	var variables = box.Vars.Lookup()
	if text.SplitCount(box.Targets, " ") == 4 {
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 0), depth+1)
		var rx = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variables)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 1), depth+1)
		var ry = text.Calculate(text.SplitIndex(box.Expression, " ", 1), variables)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 2), depth+1)
		var rw = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variables)
		setTargetVars(layout, &box.Vars, text.SplitIndex(box.Targets, " ", 3), depth+1)
		var rh = text.Calculate(text.SplitIndex(box.Expression, " ", 3), variables)
		return rx, ry, rw, rh
	}

	var rx = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variables)
	var ry = text.Calculate(text.SplitIndex(box.Expression, " ", 1), variables)
	var rw = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variables)
	var rh = text.Calculate(text.SplitIndex(box.Expression, " ", 3), variables)
	return rx, ry, rw, rh
}
func itemDynamic(layout *internal.Layout, itemId int) (x, y, w, h float32) {
	var item = &layout.Items[itemId]
	var box = &layout.Boxes[item.BoxId]
	var bx, by, bw, bh = boxDynamic(layout, int(item.BoxId), 0)

	item.Vars = internal.Vars{}

	item.Vars.Ow, item.Vars.Oh, item.Vars.Ov = bw, bh, 1

	var osx, osy float32
	if text.SplitCount(box.ItemSpacing, " ") >= 2 {
		osx = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 0))
		osy = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 1))
	}
	var og, mb = box.ItemGap, box.ItemNewRow
	item.Vars.Osx, item.Vars.Osy, item.Vars.Og, item.Vars.Mnr = osx, osy, og, mb

	if text.SplitCount(box.ItemSize, " ") >= 2 {
		var vars = item.Vars.Lookup()
		item.Vars.Mw = text.Calculate(text.SplitIndex(box.ItemSize, " ", 0), vars)
		item.Vars.Mh = text.Calculate(text.SplitIndex(box.ItemSize, " ", 1), vars)
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

	for i := 0; i < len(layout.Items); i++ {
		var it = &layout.Items[i]
		if it.BoxId != item.BoxId {
			continue
		}

		if it.NewRow == 1 {
			curY += rowMaxH + mb
			if it.NewRowExpression != "" {
				item.Vars.Mx, item.Vars.My, item.Vars.Mw, item.Vars.Mh = curX, curY, defW, defH
				curY += text.Calculate(it.NewRowExpression, item.Vars.Lookup())
			}
			curX = bx + osx
			rowMaxH = 0
		}

		item.Vars.Mx, item.Vars.My, item.Vars.Mw, item.Vars.Mh = curX, curY, defW, defH
		var itW, itH = defW, defH
		if text.SplitCount(it.Expression, " ") == 4 {
			var vars = item.Vars.Lookup()
			itW = text.Calculate(text.SplitIndex(it.Expression, " ", 2), vars)
			itH = text.Calculate(text.SplitIndex(it.Expression, " ", 3), vars)
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

	var variables = item.Vars.Lookup()
	var rx = text.Calculate(text.SplitIndex(item.Expression, " ", 0), variables)
	var ry = text.Calculate(text.SplitIndex(item.Expression, " ", 1), variables)
	var rw = text.Calculate(text.SplitIndex(item.Expression, " ", 2), variables)
	var rh = text.Calculate(text.SplitIndex(item.Expression, " ", 3), variables)
	return rx, ry, rw, rh
}

func setTargetVars(layout *internal.Layout, vars *internal.Vars, tar string, depth int) {
	if tar != "" {
		var targetId = text.ToNumber[int](tar)
		if targetId >= 0 && targetId < len(layout.Boxes) {
			vars.Tx, vars.Ty, vars.Tw, vars.Th = boxDynamic(layout, targetId, depth)
		} else {
			vars.Tx, vars.Ty, vars.Tw, vars.Th = 0, 0, 0, 0
		}
	} else {
		vars.Tx, vars.Ty, vars.Tw, vars.Th = 0, 0, 0, 0
	}
	vars.Tlx, vars.Tly, vars.Trx, vars.Try = vars.Tx, vars.Ty+vars.Th/2, vars.Tx+vars.Tw, vars.Tly
	vars.Tux, vars.Tuy, vars.Tdx, vars.Tdy = vars.Tx+vars.Tw/2, vars.Ty, vars.Tux, vars.Ty+vars.Th
}
