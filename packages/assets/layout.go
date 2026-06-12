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
	return dynamic(&layout, id, zoom)
}
func (l LayoutId) Item(id int) (x, y, width, height float32) {
	return
}

// private ========================================================

func dynamic(layout *internal.Layout, boxId int, zoom float32) (x, y, w, h float32) {
	var box = &layout.Boxes[boxId]
	var rect = text.Split(box.Rectangle, " ")
	var expr = text.Split(box.Expression, " ")
	var tars = text.Split(box.Targets, " ")

	if box.Vars == nil {
		box.Vars = make(map[string]float32)
	} else {
		clear(box.Vars)
	}

	box.Vars["mx"], box.Vars["my"] = text.ToNumber[float32](rect[0]), text.ToNumber[float32](rect[1])
	box.Vars["mw"], box.Vars["mh"] = text.ToNumber[float32](rect[2]), text.ToNumber[float32](rect[3])
	box.Vars["mlx"], box.Vars["mly"] = box.Vars["mx"], box.Vars["my"]+box.Vars["mh"]/2
	box.Vars["mrx"], box.Vars["mry"] = box.Vars["mx"]+box.Vars["mw"], box.Vars["mly"]
	box.Vars["mux"], box.Vars["muy"] = box.Vars["mx"]+box.Vars["mw"]/2, box.Vars["my"]
	box.Vars["mdx"], box.Vars["mdy"] = box.Vars["mux"], box.Vars["my"]+box.Vars["mh"]

	box.Vars["sx"], box.Vars["sy"] = -(internal.WindowWidth*zoom)/2, -(internal.WindowHeight*zoom)/2
	box.Vars["sw"], box.Vars["sh"] = internal.WindowWidth*zoom, internal.WindowHeight*zoom
	box.Vars["slx"], box.Vars["sly"] = box.Vars["sx"], box.Vars["sy"]+box.Vars["sh"]/2
	box.Vars["srx"], box.Vars["sry"] = box.Vars["sx"]+box.Vars["sw"], box.Vars["sly"]
	box.Vars["sux"], box.Vars["suy"] = box.Vars["sx"]+box.Vars["sw"]/2, box.Vars["sy"]
	box.Vars["sdx"], box.Vars["sdy"] = box.Vars["sux"], box.Vars["sy"]+box.Vars["sh"]

	if len(tars) == 4 && (tars[0] != "" || tars[1] != "" || tars[2] != "" || tars[3] != "") {
		var tx, _, _, _ = dynamic(layout, text.ToNumber[int](tars[0]), zoom)
		var _, ty, _, _ = dynamic(layout, text.ToNumber[int](tars[1]), zoom)
		var _, _, tw, _ = dynamic(layout, text.ToNumber[int](tars[2]), zoom)
		var _, _, _, th = dynamic(layout, text.ToNumber[int](tars[3]), zoom)
		box.Vars["tx"], box.Vars["ty"], box.Vars["tw"], box.Vars["th"] = tx, ty, tw, th
		box.Vars["tlx"], box.Vars["tly"] = box.Vars["tx"], box.Vars["ty"]+box.Vars["th"]/2
		box.Vars["trx"], box.Vars["try"] = box.Vars["tx"]+box.Vars["tw"], box.Vars["tly"]
		box.Vars["tux"], box.Vars["tuy"] = box.Vars["tx"]+box.Vars["tw"]/2, box.Vars["ty"]
		box.Vars["tdx"], box.Vars["tdy"] = box.Vars["tux"], box.Vars["ty"]+box.Vars["th"]
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
	return rx + rw/2, ry + rh/2, rw, rh
}
