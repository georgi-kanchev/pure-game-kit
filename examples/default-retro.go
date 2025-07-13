package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/window"
)

func DefaultRetro() {
	var camera = graphics.NewCamera(1)
	var assetId, tileIds = assets.LoadDefaultAtlasRetro()
	var node = graphics.NewNode(assetId)
	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		node.Fit(&camera)

		var mx, my = node.MousePosition(&camera)
		var index = number.Indexes2DToIndex1D(int(my/9), int(mx/9), 26, 21)
		var w, h = 9 * node.ScaleX, 9 * node.ScaleY
		var mmx, mmy = node.PointToCamera(&camera, number.Snap(mx-4.5, 9), number.Snap(my-4.5, 9))
		index = number.LimitInt(index, 0, len(tileIds)-1)
		fmt.Printf("tileIds[index]: %v\n", tileIds[index])

		camera.DrawNodes(&node)
		camera.DrawFrame(mmx, mmy, w, h, 0, 4, color.Red)
	}
}
