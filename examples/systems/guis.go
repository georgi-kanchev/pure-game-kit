package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GUIs() {
	var cam = graphics.NewCamera(1)
	var _, _, b = assets.LoadDefaultAtlasUI()
	var menu = gui.New(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("button", p.Color, "255 255 255 255", p.Width, "300", p.Height, "100", p.GapX, "20", p.GapY, "20",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, b[5], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.3", p.TextColor, "80 80 80 255",
			p.TextLineHeight, "70", p.ButtonAssetIdPress, b[4],
		),
		gui.Container("top", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "600",
			p.ThemeId, "button", p.GapX, "30", p.GapY, "30",
		),
		gui.WidgetVisual("background", p.FillContainer, "", p.AssetId, b[8], p.Color, "200 200 200 255"),
		gui.WidgetButton("btn1"),
		gui.WidgetButton("btn2", p.Width, "100"),
		gui.WidgetButton("btn3"),
		gui.WidgetButton("btn4", p.NewRow, ""),
		gui.WidgetButton("btn5", p.Text, "text", p.Width, "200"),
		gui.WidgetButton("btn6"),
		gui.WidgetVisual("btn7", p.Text, "long text", p.TextColor, "255 0 0 255"),
		gui.WidgetButton("btn8", p.Text, "BUTTON", p.NewRow, ""),
		gui.WidgetButton("btn9", p.Color, "200 50 50 255"),
	)

	cam.Angle = 45

	assets.LoadDefaultFont()

	window.IsVSynced = false
	window.TargetFrameRate = 255

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		menu.Draw(&cam)
	}
}
