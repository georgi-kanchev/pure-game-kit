package assets

import (
	"math"
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
	var rx, ry, rw, rh = dynamic(&layout, id, nil, 0)
	return (rx + rw/2) * zoom, (ry + rh/2) * zoom, rw * zoom, rh * zoom
}

func (l LayoutId) Item(id int, zoom float32) (x, y, width, height float32) {
	var layout, has = internal.Layouts[uint32(l)]
	if !has || id < 0 || id >= len(layout.Items) {
		return
	}
	var rx, ry, rw, rh = itemDynamic(&layout, id)
	return (rx + rw/2) * zoom, (ry + rh/2) * zoom, rw * zoom, rh * zoom
}

// private ========================================================

func dynamic(layout *internal.Layout, boxId int, resolving map[int]bool, depth int) (x, y, w, h float32) {
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

	// My rect (from stored rectangle — editor-viewport pixels; scale to game resolution)
	// Editor uses base=512, so the ratio from editor to game is sqrt(W·H)/512.
	var rectScale = float32(math.Sqrt(float64(internal.WindowWidth*internal.WindowHeight))) / 512
	box.Vars["mx"], box.Vars["my"] = text.ToNumber[float32](rect[0])*rectScale, text.ToNumber[float32](rect[1])*rectScale
	box.Vars["mw"], box.Vars["mh"] = text.ToNumber[float32](rect[2])*rectScale, text.ToNumber[float32](rect[3])*rectScale
	box.Vars["mlx"], box.Vars["mly"] = box.Vars["mx"], box.Vars["my"]+box.Vars["mh"]/2
	box.Vars["mrx"], box.Vars["mry"] = box.Vars["mx"]+box.Vars["mw"], box.Vars["mly"]
	box.Vars["mux"], box.Vars["muy"] = box.Vars["mx"]+box.Vars["mw"]/2, box.Vars["my"]
	box.Vars["mdx"], box.Vars["mdy"] = box.Vars["mux"], box.Vars["my"]+box.Vars["mh"]

	// Screen rect (formulas reference these for layout math)
	box.Vars["sx"], box.Vars["sy"] = -internal.WindowWidth/2, -internal.WindowHeight/2
	box.Vars["sw"], box.Vars["sh"] = internal.WindowWidth, internal.WindowHeight
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
		if tars[0] != "" {
			var tid = text.ToNumber[int](tars[0])
			if !resolving[tid] && tid >= 0 && tid < len(layout.Boxes) {
				resolving[boxId] = true
				box.Vars["tx"], _, _, _ = dynamic(layout, tid, resolving, depth+1)
				delete(resolving, boxId)
			}
		}
		if tars[1] != "" {
			var tid = text.ToNumber[int](tars[1])
			if !resolving[tid] && tid >= 0 && tid < len(layout.Boxes) {
				resolving[boxId] = true
				_, box.Vars["ty"], _, _ = dynamic(layout, tid, resolving, depth+1)
				delete(resolving, boxId)
			}
		}
		if tars[2] != "" {
			var tid = text.ToNumber[int](tars[2])
			if !resolving[tid] && tid >= 0 && tid < len(layout.Boxes) {
				resolving[boxId] = true
				_, _, box.Vars["tw"], _ = dynamic(layout, tid, resolving, depth+1)
				delete(resolving, boxId)
			}
		}
		if tars[3] != "" {
			var tid = text.ToNumber[int](tars[3])
			if !resolving[tid] && tid >= 0 && tid < len(layout.Boxes) {
				resolving[boxId] = true
				_, _, _, box.Vars["th"] = dynamic(layout, tid, resolving, depth+1)
				delete(resolving, boxId)
			}
		}

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

func itemDynamic(layout *internal.Layout, itemId int) (x, y, w, h float32) {
	var item = &layout.Items[itemId]
	var box = &layout.Boxes[item.BoxId]

	// Resolve the owning box (raw top-left result)
	var bx, by, bw, bh = dynamic(layout, int(item.BoxId), nil, 0)

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
