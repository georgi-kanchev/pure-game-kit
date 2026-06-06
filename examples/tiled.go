package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Tiled() {
	window.Create("example - tiled", false, true)
	var view = graphics.NewView(3)
	var atlasId, layerIds = assets.LoadTiledLayers("examples/data/map.tmx")

	var layers []graphics.Object
	for _, id := range layerIds {
		var layer = graphics.NewTilemap(1, atlasId, id)
		layers = append(layers, layer)
	}

	var shapes = layers[1].TilemapShapes()
	var objs []graphics.Object
	for _, v := range shapes {
		var obj = graphics.Object{Shape: v}
		obj.Effects.Tint = palette.White
		obj.Effects.FillColor = palette.Red
		objs = append(objs, obj)
	}

	var cellShapes = layers[3].TilemapShapes()
	var cellObjs []graphics.Object
	for _, v := range cellShapes {
		var obj = graphics.Object{Shape: v}
		obj.Effects.Tint = palette.White
		obj.Effects.FillColor = palette.Red
		cellObjs = append(cellObjs, obj)
	}

	for window.KeepOpen() {
		view.MouseDragAndZoomSmoothly()

		view.DrawGrid(1, 16, 16, palette.DarkGray)

		for _, l := range layers {
			view.DrawObjects(&l)
		}
		for _, o := range objs {
			view.DrawObjects(&o)
		}
		for _, c := range cellObjs {
			view.DrawObjects(&c)
		}
	}
}
