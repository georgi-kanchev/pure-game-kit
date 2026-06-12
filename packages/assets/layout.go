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
	var rx, ry, rw, rh = dynamic(&layout, id, 0)
	return (rx + rw/2), (ry + rh/2), rw, rh
}

func (l LayoutId) Item(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var rx, ry, rw, rh = itemDynamic(&layout, id)
	return (rx + rw/2), (ry + rh/2), rw, rh
}

// private ========================================================

func dynamic(layout *internal.Layout, boxId int, depth int) (x, y, w, h float32) {
	if depth > 8 {
		return
	}

	var box = &layout.Boxes[boxId]
	var rect = text.Split(box.Rectangle, " ")
	var expr = text.Split(box.Expression, " ")
	var tars = text.Split(box.Targets, " ")

	if box.Vars == nil {
		box.Vars = make(map[string]float32)
	} else {
		clear(box.Vars)
	}

	var rectScale float32 = 1 //float32(math.Sqrt(float64(internal.WindowWidth*internal.WindowHeight))) / 512
	box.Vars["mx"], box.Vars["my"] = text.ToNumber[float32](rect[0])*rectScale, text.ToNumber[float32](rect[1])*rectScale
	box.Vars["mw"], box.Vars["mh"] = text.ToNumber[float32](rect[2])*rectScale, text.ToNumber[float32](rect[3])*rectScale
	box.Vars["mlx"], box.Vars["mly"] = box.Vars["mx"], box.Vars["my"]+box.Vars["mh"]/2
	box.Vars["mrx"], box.Vars["mry"] = box.Vars["mx"]+box.Vars["mw"], box.Vars["mly"]
	box.Vars["mux"], box.Vars["muy"] = box.Vars["mx"]+box.Vars["mw"]/2, box.Vars["my"]
	box.Vars["mdx"], box.Vars["mdy"] = box.Vars["mux"], box.Vars["my"]+box.Vars["mh"]

	box.Vars["sx"], box.Vars["sy"] = -internal.WindowWidth/2, -internal.WindowHeight/2
	box.Vars["sw"], box.Vars["sh"] = internal.WindowWidth, internal.WindowHeight
	box.Vars["slx"], box.Vars["sly"] = box.Vars["sx"], box.Vars["sy"]+box.Vars["sh"]/2
	box.Vars["srx"], box.Vars["sry"] = box.Vars["sx"]+box.Vars["sw"], box.Vars["sly"]
	box.Vars["sux"], box.Vars["suy"] = box.Vars["sx"]+box.Vars["sw"]/2, box.Vars["sy"]
	box.Vars["sdx"], box.Vars["sdy"] = box.Vars["sux"], box.Vars["sy"]+box.Vars["sh"]

	box.Vars["tx"], box.Vars["ty"], box.Vars["tw"], box.Vars["th"] = 0, 0, 0, 0
	box.Vars["tlx"], box.Vars["tly"] = 0, 0
	box.Vars["trx"], box.Vars["try"] = 0, 0
	box.Vars["tux"], box.Vars["tuy"] = 0, 0
	box.Vars["tdx"], box.Vars["tdy"] = 0, 0

	var variables = varLookup(box.Vars)

	if len(tars) == 4 {
		setTargetVars(layout, box.Vars, tars[0], depth+1)
		var rx = text.Calculate(expr[0], variables)

		setTargetVars(layout, box.Vars, tars[1], depth+1)
		var ry = text.Calculate(expr[1], variables)

		setTargetVars(layout, box.Vars, tars[2], depth+1)
		var rw = text.Calculate(expr[2], variables)

		setTargetVars(layout, box.Vars, tars[3], depth+1)
		var rh = text.Calculate(expr[3], variables)

		return rx, ry, rw, rh
	}

	var rx = text.Calculate(expr[0], variables)
	var ry = text.Calculate(expr[1], variables)
	var rw = text.Calculate(expr[2], variables)
	var rh = text.Calculate(expr[3], variables)
	return rx, ry, rw, rh
}

func setTargetVars(layout *internal.Layout, vars map[string]float32, tar string, depth int) {
	if tar != "" {
		var tid = text.ToNumber[int](tar)
		if tid >= 0 && tid < len(layout.Boxes) {
			vars["tx"], vars["ty"], vars["tw"], vars["th"] = dynamic(layout, tid, depth)
		} else {
			vars["tx"], vars["ty"], vars["tw"], vars["th"] = 0, 0, 0, 0
		}
	} else {
		vars["tx"], vars["ty"], vars["tw"], vars["th"] = 0, 0, 0, 0
	}
	vars["tlx"] = vars["tx"]
	vars["tly"] = vars["ty"] + vars["th"]/2
	vars["trx"] = vars["tx"] + vars["tw"]
	vars["try"] = vars["tly"]
	vars["tux"] = vars["tx"] + vars["tw"]/2
	vars["tuy"] = vars["ty"]
	vars["tdx"] = vars["tux"]
	vars["tdy"] = vars["ty"] + vars["th"]
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

func itemDynamic(layout *internal.Layout, itemId int) (x, y, w, h float32) {
	var item = &layout.Items[itemId]
	var box = &layout.Boxes[item.BoxId]

	var bx, by, bw, bh = dynamic(layout, int(item.BoxId), 0)

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

	var itSz = text.Split(box.ItemSize, " ")
	if len(itSz) >= 2 {
		var look = varLookup(item.Variables)
		item.Variables["mw"] = text.Calculate(itSz[0], look)
		item.Variables["mh"] = text.Calculate(itSz[1], look)
	}
	if number.IsNaN(item.Variables["mw"]) {
		item.Variables["mw"] = 40
	}
	if number.IsNaN(item.Variables["mh"]) {
		item.Variables["mh"] = 20
	}

	item.Variables["mx"], item.Variables["my"] = bx, by

	var expr = text.Split(item.Expression, " ")
	var variables = varLookup(item.Variables)

	var rx = text.Calculate(expr[0], variables)
	if number.IsNaN(rx) {
		rx = item.Variables["mx"]
	}
	var ry = text.Calculate(expr[1], variables)
	if number.IsNaN(ry) {
		ry = item.Variables["my"]
	}
	var rw = text.Calculate(expr[2], variables)
	if number.IsNaN(rw) {
		rw = item.Variables["mw"]
	}
	var rh = text.Calculate(expr[3], variables)
	if number.IsNaN(rh) {
		rh = item.Variables["mh"]
	}
	return rx, ry, rw, rh
}
