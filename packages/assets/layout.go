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

func (l LayoutId) Box(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Boxes) {
		return
	}
	var rx, ry, rw, rh = boxDynamic(&layout, id, 0)
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	return (rx + rw/2) * sc, (ry + rh/2) * sc, rw * sc, rh * sc
}
func (l LayoutId) Item(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var rx, ry, rw, rh = itemDynamic(&layout, id)
	var sc = number.SquareRoot(internal.WindowWidth*internal.WindowHeight) / 512
	return (rx + rw/2) * sc, (ry + rh/2) * sc, rw * sc, rh * sc
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

	item.Variables["ow"], item.Variables["oh"] = bw, bh
	item.Variables["ov"] = 1

	// parse spacing, gap, newRow base
	var osx, osy float32
	if text.SplitCount(box.ItemSpacing, " ") >= 2 {
		osx = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 0))
		osy = text.ToNumber[float32](text.SplitIndex(box.ItemSpacing, " ", 1))
	}
	item.Variables["osx"], item.Variables["osy"] = osx, osy
	var og = float32(box.ItemGap)
	var mb = float32(box.ItemNewRow)
	item.Variables["og"], item.Variables["mnr"] = og, mb

	// default item size from box.ItemSize (evaluated with ow/oh)
	if text.SplitCount(box.ItemSize, " ") >= 2 {
		var vars = varLookup(item.Variables)
		item.Variables["mw"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 0), vars)
		item.Variables["mh"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 1), vars)
	}
	// fall back to item.Rectangle
	if number.IsNaN(item.Variables["mw"]) && text.SplitCount(item.Rectangle, " ") == 4 {
		item.Variables["mw"] = text.ToNumber[float32](text.SplitIndex(item.Rectangle, " ", 2))
	}
	if number.IsNaN(item.Variables["mh"]) && text.SplitCount(item.Rectangle, " ") == 4 {
		item.Variables["mh"] = text.ToNumber[float32](text.SplitIndex(item.Rectangle, " ", 3))
	}
	// absolute fallback
	if number.IsNaN(item.Variables["mw"]) {
		item.Variables["mw"] = 40
	}
	if number.IsNaN(item.Variables["mh"]) {
		item.Variables["mh"] = 20
	}
	var defW = item.Variables["mw"]
	var defH = item.Variables["mh"]

	// flow layout: walk items in this box to compute mx, my
	var curX = bx + osx
	var curY = by + osy
	var rowMaxH float32 = 0

	for i := 0; i < len(layout.Items); i++ {
		var it = &layout.Items[i]
		if it.BoxId != item.BoxId {
			continue
		}

		// newRow: advance Y and reset X
		if it.NewRow == 1 {
			curY += rowMaxH + mb
			if it.NewRowExpression != "" {
				var brkVars = map[string]float32{
					"ow":  bw, "oh":  bh, "ov": 1,
					"osx": osx, "osy": osy, "og": og, "mnr": mb,
					"mx": curX, "my": curY, "mw": defW, "mh": defH,
				}
				var brk = text.Calculate(it.NewRowExpression, varLookup(brkVars))
				if !number.IsNaN(brk) {
					curY += brk
				}
			}
			curX = bx + osx
			rowMaxH = 0
		}

		// compute this item's w/h for flow advance
		var itW = defW
		var itH = defH
		if text.SplitCount(it.Expression, " ") == 4 {
			var flowVars = map[string]float32{
				"ow":  bw, "oh":  bh, "ov": 1,
				"osx": osx, "osy": osy, "og": og, "mnr": mb,
				"mx": curX, "my": curY, "mw": defW, "mh": defH,
			}
			var look = varLookup(flowVars)
			var fw = text.Calculate(text.SplitIndex(it.Expression, " ", 2), look)
			if !number.IsNaN(fw) {
				itW = fw
			}
			var fh = text.Calculate(text.SplitIndex(it.Expression, " ", 3), look)
			if !number.IsNaN(fh) {
				itH = fh
			}
		}

		if i == itemId {
			item.Variables["mx"] = curX
			item.Variables["my"] = curY
			break
		}

		curX += itW + og
		if itH > rowMaxH {
			rowMaxH = itH
		}
	}

	// fallback if flow didn't set mx/my (e.g. item not found in loop)
	if number.IsNaN(item.Variables["mx"]) {
		item.Variables["mx"] = bx + osx
	}
	if number.IsNaN(item.Variables["my"]) {
		item.Variables["my"] = by + osy
	}

	var variables = varLookup(item.Variables)
	var rx = text.Calculate(text.SplitIndex(item.Expression, " ", 0), variables)
	if number.IsNaN(rx) {
		rx = item.Variables["mx"]
	}
	var ry = text.Calculate(text.SplitIndex(item.Expression, " ", 1), variables)
	if number.IsNaN(ry) {
		ry = item.Variables["my"]
	}
	var rw = text.Calculate(text.SplitIndex(item.Expression, " ", 2), variables)
	if number.IsNaN(rw) {
		rw = item.Variables["mw"]
	}
	var rh = text.Calculate(text.SplitIndex(item.Expression, " ", 3), variables)
	if number.IsNaN(rh) {
		rh = item.Variables["mh"]
	}
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
