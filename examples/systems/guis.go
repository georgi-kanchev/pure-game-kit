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
	var atlas, icons = assets.LoadDefaultAtlasIcons()
	var _, _, box = assets.LoadDefaultAtlasUI()
	var menu = gui.New(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("button", p.Color, "220 220 220 255", p.Width, "300", p.Height, "100", p.GapX, "20", p.GapY, "20",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[2], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.3", p.TextColor, "80 80 80 255",
			p.TextLineHeight, "70", p.ButtonHoverThemeId, "button-hover", p.ButtonPressThemeId, "button-press",
		),
		gui.Theme("button-hover", p.Color, "255 255 255 255", p.Width, "300", p.Height, "100",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40",
			p.BoxEdgeBottom, "40", p.AssetId, box[5], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.3",
			p.TextColor, "127 127 127 255", p.TextLineHeight, "70", p.GapX, "20", p.GapY, "20",
		),
		gui.Theme("button-press", p.Color, "200 200 200 255", p.Width, "300", p.Height, "100",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[4], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.6", p.TextColor, "80 80 80 255",
			p.TextLineHeight, "70", p.GapX, "20", p.GapY, "20",
		),
		gui.Theme("label", p.Color, "0 0 0 0", p.Width, "300", p.Height, "100", p.GapX, "20", p.GapY, "20",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.TextAlignmentX, "0", p.TextAlignmentY, "0.5", p.TextColor, "0 0 0 255",
			p.TextLineHeight, "80",
		),
		// ======================================================
		gui.Container("top", d.CameraLeftX+"+10", d.CameraTopY+"+10", d.CameraWidth+"-20", "1175",
			p.ThemeId, "button", p.GapX, "50", p.GapY, "50"),
		gui.Visual("background", p.FillContainer, "", p.AssetId, box[8], p.Color, "200 200 200 255"),
		gui.Visual("name-label", p.ThemeId, "label", p.Text, "Name"),
		gui.Visual("name", p.Width, "500", p.AssetId, box[9], p.Text, "Kenney", p.TextAlignmentX, "0.1",
			p.TextAlignmentY, "0.5", p.TextColor, "150 150 150 255"),
		gui.Visual("stepper-label", p.ThemeId, "label", p.Text, "Stepper", p.NewRow, ""),
		gui.Button("step-left", p.Width, "100", p.TextEmbeddedAssetId1, "arrow-left", p.Text, "^^ "),
		gui.Visual("stepper", p.AssetId, box[9], p.Text, "10/10", p.TextAlignmentY, "0.5", p.GapX, "0",
			p.TextColor, "150 150 150 255"),
		gui.Button("step-right", p.Width, "100", p.TextEmbeddedAssetId1, icons[212], p.Text, "^^", p.GapX, "0"),
		gui.Visual("checkbox-label", p.ThemeId, "label", p.Text, "Checkbox", p.NewRow, ""),
		gui.Button("checkbox", p.Width, "100", p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ ", p.AssetId, box[9],
			p.TextAlignmentX, "0.6", p.TextAlignmentY, "0.53"),
		gui.Button("dropdown", p.NewRow, "", p.AssetId, box[9], p.Text, "^^ List selection", p.Width, "820",
			p.TextAlignmentX, "0.05", p.TextAlignmentY, "0.5", p.TextColor, "150 150 150 255",
			p.TextEmbeddedAssetId1, "arrow-down"),
		gui.Visual("sliders-label", p.ThemeId, "label", p.Text, "Sliders", p.TextLineHeight, "60", p.NewRow, ""),
		gui.Visual("slider1", p.AssetId, box[10], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Visual("slider1-0", p.ThemeId, "label", p.Text, "0", p.TextLineHeight, "50", p.NewRow, "10",
			p.Width, "50", p.TextAlignmentX, "0.5"),
		gui.Visual("slider1-100", p.ThemeId, "label", p.Text, "100", p.TextLineHeight, "50", p.GapX, "710",
			p.Width, "80", p.TextAlignmentX, "0.5"),
		gui.Visual("slider2", p.AssetId, box[10], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Visual("slider2-0", p.ThemeId, "label", p.Text, "0", p.TextLineHeight, "50", p.NewRow, "10",
			p.Width, "50", p.TextAlignmentX, "0.5"),
		gui.Visual("slider2-100", p.ThemeId, "label", p.Text, "100", p.TextLineHeight, "50", p.GapX, "710",
			p.Width, "80", p.TextAlignmentX, "0.5"),
		gui.Visual("divider", p.AssetId, box[12], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Button("x", p.Width, "170", p.Height, "140", p.TextEmbeddedAssetId1, icons[250], p.Text, "^^ ",
			p.NewRow, "70", p.TextColor, "255 255 255 255", p.Color, "200 0 0 255"),
		gui.Button("v", p.Width, "630", p.Height, "140", p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ Accept ",
			p.TextColor, "255 255 255 255", p.Color, "0 200 0 255"),
	)

	cam.Angle = 45

	assets.LoadDefaultFont()

	assets.SetTextureAtlasTile(atlas, "arrow-left", 14, 9, 1, 1, 0, true)
	assets.SetTextureAtlasTile(atlas, "arrow-down", 14, 9, 1, 1, 1, false)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		if menu.ButtonClickedOnce("dropdown", cam) {
			menu.SetProperty("dropdown", p.GapX, menu.Property("dropdown", p.GapX)+"+20")
		}

		menu.Draw(cam)
	}
}
