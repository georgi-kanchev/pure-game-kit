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

	box.Vars["mx"], box.Vars["my"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 0)), text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 1))
	box.Vars["mw"], box.Vars["mh"] = text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 2)), text.ToNumber[float32](text.SplitIndex(box.Rectangle, " ", 3))
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

	var rx, ry = text.Calculate(text.SplitIndex(box.Expression, " ", 0), variables), text.Calculate(text.SplitIndex(box.Expression, " ", 1), variables)
	var rw, rh = text.Calculate(text.SplitIndex(box.Expression, " ", 2), variables), text.Calculate(text.SplitIndex(box.Expression, " ", 3), variables)
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
	item.Variables["osx"], item.Variables["osy"] = 0, 0
	item.Variables["og"] = float32(box.ItemGap)
	item.Variables["mnr"] = float32(box.ItemNewRow)

	if text.SplitCount(box.ItemSize, " ") >= 2 {
		var look = varLookup(item.Variables)
		item.Variables["mw"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 0), look)
		item.Variables["mh"] = text.Calculate(text.SplitIndex(box.ItemSize, " ", 1), look)
	}
	if number.IsNaN(item.Variables["mw"]) {
		item.Variables["mw"] = 40
	}
	if number.IsNaN(item.Variables["mh"]) {
		item.Variables["mh"] = 20
	}

	item.Variables["mx"], item.Variables["my"] = bx, by

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
	vars["tlx"], vars["tly"] = vars["tx"], vars["ty"]+vars["th"]/2
	vars["trx"], vars["try"] = vars["tx"]+vars["tw"], vars["tly"]
	vars["tux"], vars["tuy"] = vars["tx"]+vars["tw"]/2, vars["ty"]
	vars["tdx"], vars["tdy"] = vars["tux"], vars["ty"]+vars["th"]
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
