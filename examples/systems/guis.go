package example

import (
	"fmt"
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui"
	d "pure-kit/engine/gui/dynamic"
	p "pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func GUIs() {
	window.TargetFrameRate = 0
	window.IsVSynced = false

	var cam = graphics.NewCamera(1)
	var atlas, icons = assets.LoadDefaultAtlasIcons(true)
	var _, ids, box = assets.LoadDefaultAtlasUI(true)
	var hud = gui.NewElements(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("label", p.Color, "0 0 0 0", p.Width, "300", p.Height, "100", p.GapX, "20", p.GapY, "20",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.TextAlignmentX, "0", p.TextAlignmentY, "0.5", p.TextColor, "0 0 0 255",
			p.TextLineHeight, "80", p.TooltipId, "tooltip"),
		// ======================================================
		gui.Theme("button", p.Color, "220 220 220 255", p.Width, "300", p.Height, "100", p.GapX, "20", p.GapY, "20",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[2], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.3", p.TextColor, "80 80 80 255",
			p.TextLineHeight, "70", p.ButtonThemeIdHover, "button-hover", p.ButtonThemeIdPress, "button-press",
			p.TooltipId, "tooltip", p.SliderStep, "0.1", p.SliderHandleAssetId, box[14],
			p.SliderStepAssetId, ids[49], p.DraggableSpriteColor, "0 0 255", p.DraggableSpriteScale, "0.8"),
		gui.Theme("button-hover", p.Color, "255 255 255 255", p.Width, "300", p.Height, "100",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40",
			p.BoxEdgeBottom, "40", p.AssetId, box[5], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.3",
			p.TextColor, "127 127 127 255", p.TextLineHeight, "70", p.GapX, "20", p.GapY, "20",
			p.TooltipId, "tooltip"),
		gui.Theme("button-press", p.Color, "200 200 200 255", p.Width, "300", p.Height, "100",
			p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[4], p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.6", p.TextColor, "80 80 80 255",
			p.TextLineHeight, "70", p.GapX, "20", p.GapY, "20", p.TooltipId, "tooltip"),
		// ======================================================
		gui.Theme("checkbox-on", p.Color, "220 220 220 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[9], p.TextColor, "80 80 80 255", p.TextLineHeight, "70",
			p.ButtonThemeIdHover, "checkbox-on-hover", p.ButtonThemeIdPress, "checkbox-on-press",
			p.TooltipId, "tooltip", p.TooltipText, "Currently on!",
			p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ ", p.TextAlignmentX, "0.6", p.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-hover", p.Color, "255 255 255 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[9], p.TextColor, "80 80 80 255", p.TextLineHeight, "70",
			p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ ", p.TextAlignmentX, "0.6", p.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-on-press", p.Color, "200 200 200 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[9], p.TextColor, "80 80 80 255", p.TextLineHeight, "70",
			p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ ", p.TextAlignmentX, "0.6", p.TextAlignmentY, "0.53"),
		gui.Theme("checkbox-off", p.Color, "220 220 220 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.ButtonThemeIdHover, "checkbox-off-hover", p.ButtonThemeIdPress, "checkbox-off-press",
			p.AssetId, box[9], p.TooltipId, "tooltip", p.TooltipText, "Currently off!"),
		gui.Theme("checkbox-off-hover", p.Color, "255 255 255 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[9]),
		gui.Theme("checkbox-off-press", p.Color, "200 200 200 255", p.Width, "100", p.Height, "100", p.GapX, "20",
			p.GapY, "20", p.BoxEdgeLeft, "40", p.BoxEdgeRight, "40", p.BoxEdgeTop, "40", p.BoxEdgeBottom, "40",
			p.AssetId, box[9]),
		// ======================================================
		gui.Container("panel", d.CameraLeftX+"+10", d.CameraBottomY+"-650", d.CameraWidth+"-20", "650",
			p.ThemeId, "button", p.GapX, "40", p.GapY, "20"),
		gui.Visual("background", p.FillContainer, "", p.AssetId, box[8], p.Color, "200 200 200 255"),
		// ======================================================
		gui.Visual("name-label", p.ThemeId, "label", p.Text, "Name", p.TooltipText, "Wow, tooltip for labels!"),
		gui.InputField("name", p.Width, "500", p.AssetId, box[9], p.Text, "Kenney",
			p.InputFieldPlaceholder, "Your name..."),
		gui.Visual("class-label", p.ThemeId, "label", p.Text, "Class", p.NewRow, ""),
		gui.InputField("class", p.Width, "500", p.AssetId, box[9], p.Text, "Cool"),
		gui.Visual("stepper-label", p.ThemeId, "label", p.Text, "Stepper", p.NewRow, ""),
		gui.Button("step-left", p.Width, "100", p.TextEmbeddedAssetId1, "arrow-left", p.Text, "^^ ",
			p.TooltipText, "Press this button to do absolutely nothing.", p.ButtonHotkey, "A"),
		gui.Visual("stepper", p.AssetId, box[9], p.Text, "10/10", p.TextAlignmentY, "0.5", p.GapX, "0",
			p.TextColor, "150 150 150 255"),
		gui.Button("step-right", p.Width, "100", p.TextEmbeddedAssetId1, icons[212], p.Text, "^^", p.GapX, "0"),
		gui.Visual("checkbox-label", p.ThemeId, "label", p.Text, "Checkbox", p.NewRow, ""),
		gui.Checkbox("checkbox", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on"),
		gui.Menu("dropdown", p.NewRow, "", p.AssetId, box[9], p.Text, "^^ List selection", p.Width, "820",
			p.TextAlignmentX, "0.05", p.TextAlignmentY, "0.5", p.TextColor, "150 150 150 255",
			p.TextEmbeddedAssetId1, "arrow-down", p.MenuContainerId, "menu"),
		gui.Visual("sliders-label", p.ThemeId, "label", p.Text, "Sliders", p.TextLineHeight, "60", p.NewRow, ""),
		gui.Slider("slider1", p.AssetId, box[10], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Visual("slider1-0", p.ThemeId, "label", p.Text, "0", p.TextLineHeight, "50", p.NewRow, "50",
			p.Width, "50", p.Height, "50", p.TextAlignmentX, "0.5"),
		gui.Visual("slider1-100", p.ThemeId, "label", p.Text, "100", p.TextLineHeight, "50", p.GapX, "710",
			p.Width, "80", p.Height, "50", p.TextAlignmentX, "0.5"),
		gui.Slider("slider2", p.AssetId, box[10], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Visual("slider2-0", p.ThemeId, "label", p.Text, "0", p.TextLineHeight, "50", p.NewRow, "50",
			p.Width, "50", p.Height, "50", p.TextAlignmentX, "0.5"),
		gui.Visual("slider2-100", p.ThemeId, "label", p.Text, "100", p.TextLineHeight, "50", p.GapX, "710",
			p.Width, "80", p.Height, "50", p.TextAlignmentX, "0.5"),
		gui.Visual("divider", p.AssetId, box[12], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40", p.TooltipText, "Tooltips for dividers?! WHAT?"),
		gui.Button("x", p.Width, "170", p.Height, "140", p.TextEmbeddedAssetId1, icons[250], p.Text, "^^ ",
			p.NewRow, "70", p.TextColor, "255 255 255 255", p.Color, "200 0 0 255",
			p.TooltipText, "This is a pretty squarish X button."),
		gui.Button("v", p.Width, "630", p.Height, "140", p.TextEmbeddedAssetId1, icons[249], p.Text, "^^ Accept ",
			p.TextColor, "255 255 255 255", p.Color, "0 200 0 255"),
		gui.Checkbox("outside1", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on", p.NewRow, "",
			p.CheckboxGroup, "radio"),
		gui.Checkbox("outside2", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on", p.NewRow, "",
			p.CheckboxGroup, "radio"),
		gui.Checkbox("outside3", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on", p.NewRow, "",
			p.CheckboxGroup, "radio"),
		gui.Checkbox("outside4", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on", p.NewRow, "",
			p.CheckboxGroup, "radio"),
		gui.Checkbox("outside5", p.ThemeId, "checkbox-off", p.CheckboxThemeId, "checkbox-on", p.NewRow, "",
			p.CheckboxGroup, "radio"),
		gui.Visual("divider-2", p.AssetId, box[12], p.BoxEdgeTop, "0", p.BoxEdgeBottom, "0", p.NewRow, "",
			p.Width, "820", p.Height, "40"),
		gui.Draggable("slot-1", p.NewRow, "50", p.Width, "200", p.Height, "200", p.AssetId, box[0],
			p.DraggableSpriteId, icons[50]),
		gui.Draggable("slot-2", p.Width, "200", p.Height, "200", p.AssetId, box[0], p.DraggableSpriteId, icons[100]),
		gui.Draggable("slot-3", p.Width, "200", p.Height, "200", p.AssetId, box[0], p.TextColor, "test"),
		gui.Draggable("slot-4", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-5", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-6", p.NewRow, "", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-7", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-8", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-9", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		gui.Draggable("slot-10", p.Width, "200", p.Height, "200", p.AssetId, box[0]),
		// ======================================================
		gui.Container("menu", "", "", "820", "300", p.ThemeId, "button", p.Hidden, "+", p.GapX, "10", p.GapY, "10"),
		gui.Visual("menu-bgr", p.FillContainer, "", p.AssetId, box[0], p.Color, "200 200 200 255"),
		gui.Button("menu-1", p.GapX, "0", p.Width, "1000", p.Text, "Monday"),
		gui.Button("menu-2", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Tuesday"),
		gui.Button("menu-3", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Wednesday"),
		gui.Button("menu-4", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Thursday",
			p.TooltipText, "It's thursday, wohooo!"),
		gui.Button("menu-5", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Friday"),
		gui.Visual("weekend-label", p.ThemeId, "label", p.Text, "Weekend", p.NewRow, ""),
		gui.Button("menu-6", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Saturday"),
		gui.Button("menu-7", p.GapX, "0", p.NewRow, "", p.Width, d.OwnerWidth, p.Text, "Sunday"),
		// ======================================================
		gui.Container("tooltips", "", "", "", "", p.ThemeId, "button", p.Hidden, "+"),
		gui.Tooltip("tooltip", p.AssetId, box[7], p.Width, "700", p.TextAlignmentX, "0.5", p.TextAlignmentY, "0.5"),
	)
	cam.Angle = 45

	assets.LoadDefaultFont()
	assets.SetTextureAtlasTile(atlas, "arrow-left", 14, 9, 1, 1, 0, true)
	assets.SetTextureAtlasTile(atlas, "arrow-down", 14, 9, 1, 1, 1, false)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.DrawGrid(2, 100, 100, color.Darken(color.Gray, 0.5))

		cam.DragAndZoom()

		var grab = hud.DragOnGrab()
		if grab != "" {
			fmt.Printf("grab: %v\n", grab)
		}

		var from, to = hud.DragOnDrop()
		if from != "" && to != "" {
			var fromId = hud.Property(from, p.DraggableSpriteId)
			var toId = hud.Property(to, p.DraggableSpriteId)
			hud.SetProperty(from, p.DraggableSpriteId, toId)
			hud.SetProperty(to, p.DraggableSpriteId, fromId)
		}

		hud.UpdateAndDraw(cam)
	}
}
