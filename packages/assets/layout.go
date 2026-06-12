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
	var rx, ry, rw, rh = dynamic(&layout, id, zoom, nil, 0)
	return rx + rw/2, ry + rh/2, rw, rh
}

func (l LayoutId) Item(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var rx, ry, rw, rh = itemDynamic(&layout, id, zoom)
	return rx + rw/2, ry + rh/2, rw, rh
}

// private ========================================================

func dynamic(layout *internal.Layout, boxId int, zoom float32, resolving map[int]bool, depth int) (x, y, w, h float32) {
	if depth > 8 {
		return
	}
	if resolving == nil {
		resolving = make(map[int]bool)
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

	// My rect (from stored rectangle — these are top-left coordinates)
	box.Vars["mx"], box.Vars["my"] = text.ToNumber[float32](rect[0]), text.ToNumber[float32](rect[1])
	box.Vars["mw"], box.Vars["mh"] = text.ToNumber[float32](rect[2]), text.ToNumber[float32](rect[3])
	box.Vars["mlx"], box.Vars["mly"] = box.Vars["mx"], box.Vars["my"]+box.Vars["mh"]/2
	box.Vars["mrx"], box.Vars["mry"] = box.Vars["mx"]+box.Vars["mw"], box.Vars["mly"]
	box.Vars["mux"], box.Vars["muy"] = box.Vars["mx"]+box.Vars["mw"]/2, box.Vars["my"]
	box.Vars["mdx"], box.Vars["mdy"] = box.Vars["mux"], box.Vars["my"]+box.Vars["mh"]

	// Screen rect (origin at top-left; formulas reference these for layout math)
	box.Vars["sx"], box.Vars["sy"] = -(internal.WindowWidth*zoom)/2, -(internal.WindowHeight*zoom)/2
	box.Vars["sw"], box.Vars["sh"] = internal.WindowWidth*zoom, internal.WindowHeight*zoom
	box.Vars["slx"], box.Vars["sly"] = box.Vars["sx"], box.Vars["sy"]+box.Vars["sh"]/2
	box.Vars["srx"], box.Vars["sry"] = box.Vars["sx"]+box.Vars["sw"], box.Vars["sly"]
	box.Vars["sux"], box.Vars["suy"] = box.Vars["sx"]+box.Vars["sw"]/2, box.Vars["sy"]
	box.Vars["sdx"], box.Vars["sdy"] = box.Vars["sux"], box.Vars["sy"]+box.Vars["sh"]

	// Default target values (zero rect)
	box.Vars["tx"], box.Vars["ty"], box.Vars["tw"], box.Vars["th"] = 0, 0, 0, 0
	box.Vars["tlx"], box.Vars["tly"] = 0, 0
	box.Vars["trx"], box.Vars["try"] = 0, 0
	box.Vars["tux"], box.Vars["tuy"] = 0, 0
	box.Vars["tdx"], box.Vars["tdy"] = 0, 0

	if len(tars) == 4 {
		var dims = [4]struct {
			key string
			tar string
		}{
			{"tx", tars[0]},
			{"ty", tars[1]},
			{"tw", tars[2]},
			{"th", tars[3]},
		}
		for _, d := range dims {
			if d.tar == "" {
				continue
			}
			var tid = text.ToNumber[int](d.tar)
			if resolving[tid] || tid < 0 || tid >= len(layout.Boxes) {
				continue
			}
			resolving[boxId] = true
			var ttx, tty, ttw, tth = dynamic(layout, tid, zoom, resolving, depth+1)
			delete(resolving, boxId)
			switch d.key {
			case "tx":
				box.Vars["tx"] = ttx
			case "ty":
				box.Vars["ty"] = tty
			case "tw":
				box.Vars["tw"] = ttw
			case "th":
				box.Vars["th"] = tth
			}
		}

		// Recompute derived t-vars
		box.Vars["tlx"] = box.Vars["tx"]
		box.Vars["tly"] = box.Vars["ty"] + box.Vars["th"]/2
		box.Vars["trx"] = box.Vars["tx"] + box.Vars["tw"]
		box.Vars["try"] = box.Vars["tly"]
		box.Vars["tux"] = box.Vars["tx"] + box.Vars["tw"]/2
		box.Vars["tuy"] = box.Vars["ty"]
		box.Vars["tdx"] = box.Vars["tux"]
		box.Vars["tdy"] = box.Vars["ty"] + box.Vars["th"]
	}

	var variables = func(variable string) float32 {
		var value, has = box.Vars[variable]
		if !has {
			return number.NaN()
		}
		return value
	}
	var rx = text.Calculate(expr[0], variables)
	var ry = text.Calculate(expr[1], variables)
	var rw = text.Calculate(expr[2], variables)
	var rh = text.Calculate(expr[3], variables)
	return rx, ry, rw, rh
}

func itemDynamic(layout *internal.Layout, itemId int, zoom float32) (x, y, w, h float32) {
	var item = &layout.Items[itemId]
	var box = &layout.Boxes[item.BoxId]

	// Resolve the owning box (raw top-left result)
	var bx, by, bw, bh = dynamic(layout, int(item.BoxId), zoom, nil, 0)

	if item.Variables == nil {
		item.Variables = make(map[string]float32)
	} else {
		clear(item.Variables)
	}

	// Item formula variables: ow, oh, ov, osx, osy, og, mnr, mx, my, mw, mh
	item.Variables["ow"], item.Variables["oh"] = bw, bh
	item.Variables["ov"] = 1
	item.Variables["osx"], item.Variables["osy"] = 0, 0
	item.Variables["og"] = float32(box.ItemGap)
	item.Variables["mnr"] = float32(box.ItemNewRow)

	// Default item size from box-level itSz — may contain expressions like "oh/2"
	var itSz = text.Split(box.ItemSize, " ")
	if len(itSz) >= 2 {
		var look = func(v string) float32 {
			var val, has = item.Variables[v]
			if has {
				return val
			}
			return number.NaN()
		}
		item.Variables["mw"] = text.Calculate(itSz[0], look)
		item.Variables["mh"] = text.Calculate(itSz[1], look)
	}
	if number.IsNaN(item.Variables["mw"]) {
		item.Variables["mw"] = 40
	}
	if number.IsNaN(item.Variables["mh"]) {
		item.Variables["mh"] = 20
	}

	// Cursor position starts at box top-left (sequential layout would accumulate from item 0)
	item.Variables["mx"], item.Variables["my"] = bx, by

	var expr = text.Split(item.Expression, " ")

	var variables = func(variable string) float32 {
		var value, has = item.Variables[variable]
		if !has {
			return number.NaN()
		}
		return value
	}

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
