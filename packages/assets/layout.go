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
	if box.Vars == nil {
		box.Vars = make(map[string]float32)
	} else {
		clear(box.Vars)
	}

	var ew = 512 * number.SquareRoot(internal.WindowWidth/internal.WindowHeight) // scales according to editor
	var eh = 512 / number.SquareRoot(internal.WindowWidth/internal.WindowHeight) // bigger windows cause bigger literals

	box.Vars["mx"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 0))
	box.Vars["my"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 1))
	box.Vars["mw"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 2))
	box.Vars["mh"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 3))
	box.Vars["mlx"], box.Vars["mly"] = box.Vars["mx"], box.Vars["my"]+box.Vars["mh"]/2
	box.Vars["mrx"], box.Vars["mry"] = box.Vars["mx"]+box.Vars["mw"], box.Vars["mly"]
	box.Vars["mux"], box.Vars["muy"] = box.Vars["mx"]+box.Vars["mw"]/2, box.Vars["my"]
	box.Vars["mdx"], box.Vars["mdy"] = box.Vars["mux"], box.Vars["my"]+box.Vars["mh"]

	box.Vars["sx"], box.Vars["sy"], box.Vars["sw"], box.Vars["sh"] = -ew/2, -eh/2, ew, eh
	box.Vars["slx"], box.Vars["sly"] = box.Vars["sx"], box.Vars["sy"]+box.Vars["sh"]/2
	box.Vars["srx"], box.Vars["sry"] = box.Vars["sx"]+box.Vars["sw"], box.Vars["sly"]
	box.Vars["sux"], box.Vars["suy"] = box.Vars["sx"]+box.Vars["sw"]/2, box.Vars["sy"]
	box.Vars["sdx"], box.Vars["sdy"] = box.Vars["sux"], box.Vars["sy"]+box.Vars["sh"]

	box.Vars["tx"], box.Vars["ty"], box.Vars["tw"], box.Vars["th"] = 0, 0, 0, 0
	box.Vars["tlx"], box.Vars["tly"], box.Vars["trx"], box.Vars["try"] = 0, 0, 0, 0
	box.Vars["tux"], box.Vars["tuy"], box.Vars["tdx"], box.Vars["tdy"] = 0, 0, 0, 0

	var variables = varLookup(box.Vars)
	if text.SplitCount(box.Targets, " ") == 4 {
		setTargetVars(layout, box.Vars, text.SplitIndex(box.Targets, " ", 0), depth+1)
		var rx = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variables)
		setTargetVars(layout, box.Vars, text.SplitIndex(box.Targets, " ", 1), depth+1)
		var ry = text.Calculate(text.SplitIndex(box.Expression, " ", 1), variables)
		setTargetVars(layout, box.Vars, text.SplitIndex(box.Targets, " ", 2), depth+1)
		var rw = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variables)
		setTargetVars(layout, box.Vars, text.SplitIndex(box.Targets, " ", 3), depth+1)
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

	if item.Variables == nil {
		item.Variables = make(map[string]float32)
	} else {
		clear(item.Variables)
	}

	item.Variables["ow"], item.Variables["oh"], item.Variables["ov"] = bw, bh, 1

	var osx, osy float32
	if text.SplitCount(box.ItemSpacing, " ") >= 2 {
		osx = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 0))
		osy = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 1))
	}
	var og, mb = box.ItemGap, box.ItemNewRow
	item.Variables["osx"], item.Variables["osy"], item.Variables["og"], item.Variables["mnr"] = osx, osy, og, mb

	if text.SplitCount(box.ItemSize, " ") >= 2 {
		var vars = varLookup(item.Variables)
		item.Variables["mw"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 0), vars)
		item.Variables["mh"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 1), vars)
	} else if text.SplitCount(item.Size, " ") >= 2 {
		item.Variables["mw"] = text.ToNumber[float32](text.SplitIndex(item.Size, " ", 0))
		item.Variables["mh"] = text.ToNumber[float32](text.SplitIndex(item.Size, " ", 1))
	} else {
		item.Variables["mw"], item.Variables["mh"] = 40, 20
	}
	var defW, defH = item.Variables["mw"], item.Variables["mh"]
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
				item.Variables["mx"], item.Variables["my"], item.Variables["mw"], item.Variables["mh"] = curX, curY, defW, defH
				curY += text.Calculate(it.NewRowExpression, varLookup(item.Variables))
			}
			curX = bx + osx
			rowMaxH = 0
		}

		item.Variables["mx"], item.Variables["my"], item.Variables["mw"], item.Variables["mh"] = curX, curY, defW, defH
		var itW, itH = defW, defH
		if text.SplitCount(it.Expression, " ") == 4 {
			var vars = varLookup(item.Variables)
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

	item.Variables["mx"], item.Variables["my"] = targetMX+(bx+bw-maxX)*alignX, targetMY+(by+bh-maxY)*alignY
	item.Variables["mw"], item.Variables["mh"] = defW, defH

	var variables = varLookup(item.Variables)
	var rx = text.Calculate(text.SplitIndex(item.Expression, " ", 0), variables)
	var ry = text.Calculate(text.SplitIndex(item.Expression, " ", 1), variables)
	var rw = text.Calculate(text.SplitIndex(item.Expression, " ", 2), variables)
	var rh = text.Calculate(text.SplitIndex(item.Expression, " ", 3), variables)
	return rx, ry, rw, rh
}

func setTargetVars(layout *internal.Layout, vars map[string]float32, tar string, depth int) {
	if tar != "" {
		var targetId = text.ToNumber[int](tar)
		if targetId >= 0 && targetId < len(layout.Boxes) {
			vars["tx"], vars["ty"], vars["tw"], vars["th"] = boxDynamic(layout, targetId, depth)
		} else {
			vars["tx"], vars["ty"], vars["tw"], vars["th"] = 0, 0, 0, 0
		}
	} else {
		vars["tx"], vars["ty"], vars["tw"], vars["th"] = 0, 0, 0, 0
	}
	vars["tlx"], vars["tly"], vars["trx"], vars["try"] = vars["tx"], vars["ty"]+vars["th"]/2, vars["tx"]+vars["tw"], vars["tly"]
	vars["tux"], vars["tuy"], vars["tdx"], vars["tdy"] = vars["tx"]+vars["tw"]/2, vars["ty"], vars["tux"], vars["ty"]+vars["th"]
}
func varLookup(vars map[string]float32) func(string) float32 {
	return func(v string) float32 {
		var value, has = vars[v]
		if !has {
			return number.NaN()
		}
		return value
	}
}
