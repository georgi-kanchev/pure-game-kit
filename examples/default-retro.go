package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"pure-kit/engine/window"
)

func DefaultRetro() {
	var camera = graphics.NewCamera(1)
	var assetId, tileIds = assets.LoadDefaultAtlasRetro()
	var node = graphics.NewNode(assetId)

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.PivotX, camera.PivotY = 0, 0
		node.Fit(&camera)
		node.ScaleX *= 0.9
		node.ScaleY *= 0.9

		var mx, my = node.MousePosition(&camera)
		var mmx, mmy = node.PointToCamera(&camera, number.Snap(mx-4.5, 9), number.Snap(my-4.5, 9))
		var imx, imy = int(mx / 9), int(my / 9)
		var index = number.LimitInt(number.Indexes2DToIndex1D(imy, imx, 26, 21), 0, len(tileIds)-1)

		camera.DrawNodes(&node)

		if !node.IsHovered(&camera) {
			continue
		}

		var info = text.New("id: ", tileIds[index], "\ncoords: ", imx, ", ", imy, "\nindex: ", index)
		camera.DrawFrame(mmx, mmy, 8*node.ScaleX, 8*node.ScaleY, 0, 6, color.Cyan)
		camera.DrawText("", info, 5, 5, 80, 10, color.Black)
		camera.DrawText("", info, 0, 0, 80, 10, color.Cyan)
	}
}
