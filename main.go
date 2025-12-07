package main

import (
	"pure-game-kit/data/assets"
	example "pure-game-kit/examples/systems"
	"pure-game-kit/graphics"
	"pure-game-kit/gui"
	d "pure-game-kit/gui/dynamic"
	f "pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	"pure-game-kit/window"
)

func main() {
	var cam = graphics.NewCamera(1)
	var _, _, box = assets.LoadDefaultAtlasUI()
	var hud = gui.NewFromXMLs(gui.NewElementsXML(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("label", f.Color, "0 0 0 0", f.Width, d.OwnerWidth+"-40", f.Height, "100", f.GapY, "50",
			f.TextAlignmentX, "0", f.TextAlignmentY, "0.5", f.TextColor, "0 0 0 255",
			f.TextLineHeight, "80"),
		gui.Theme("button", f.Color, "220 220 220 255", f.Width, d.OwnerWidth+"-40",
			f.Height, "100", f.GapX, "20", f.GapY, "0",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[2], f.TextAlignmentX, "0", f.TextAlignmentY, "0.3", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "70", f.ButtonThemeIdHover, "button-hover", f.ButtonThemeIdPress, "button-press",
			f.TextEmbeddedColor1, "170 170 170 255"),
		gui.Theme("button-hover", f.Color, "255 255 255 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40",
			f.BoxEdgeBottom, "40", f.AssetId, box[5], f.TextAlignmentX, "0", f.TextAlignmentY, "0.3",
			f.TextColor, "127 127 127 255", f.TextLineHeight, "70", f.GapX, "20", f.GapY, "0",
			f.TextEmbeddedColor1, "220 220 220 255"),
		gui.Theme("button-press", f.Color, "200 200 200 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[4], f.TextAlignmentX, "0", f.TextAlignmentY, "0.6", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "70", f.GapX, "20", f.GapY, "0",
			f.TextEmbeddedColor1, "170 170 170 255"),
		// ======================================================
		gui.Container("title", d.CameraCenterX+"-800", d.CameraTopY+"+20", "1600", "300",
			f.ThemeId, "button", f.GapX, "20", f.GapY, "20"),
		gui.Visual("background", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		// ======================================================
		gui.Visual("description", f.Text, "`pure-game-kit `- simple 2D game engine\n`~  ~  ~  Examples  ~  ~  ~",
			f.AssetId, "", f.Color, "0 0 0 0", f.Width, "1500", f.Height, d.OwnerHeight+"-90",
			f.OffsetX, "15", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5", f.TextColor, "255 255 255 255",
			f.TextLineHeight, "100", f.TextColorOutline, "0 0 0 255", f.TextThicknessOutline, "0.95",
			f.TextEmbeddedColor1, "255 0 0 255", f.TextEmbeddedColor2, "0 255 255 255",
			f.TextEmbeddedColor3, "255 255 255 255"),
		gui.Container("menu", d.CameraCenterX+"-800", d.CameraTopY+"+270", "1600", d.CameraHeight+"-290",
			f.ThemeId, "button", f.GapX, "20", f.GapY, "20"),
		gui.Visual("bg", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		//=================================================================
		gui.Visual("gfx", f.ThemeId, "label", f.Text, "Graphics:", f.GapY, "0", f.NewRow, ""),
		gui.Button("minimal graphics", f.Text, " Minimal Render `(graphics-minimal.go)", f.NewRow, ""),
		gui.Button("boxes graphics", f.Text, " Boxes `(graphics-boxes.go)", f.NewRow, ""),
		gui.Button("texts graphics", f.Text, " Texts `(graphics-texts.go)", f.NewRow, ""),
		gui.Button("guis", f.Text, " Graphical User Interfaces (GUIs) `(guis.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("input", f.ThemeId, "label", f.Text, "Input:", f.NewRow, ""),
		gui.Button("mouse input", f.Text, " Mouse `(input-mouse.go)", f.NewRow, ""),
		gui.Button("keyboard input", f.Text, " Keboard `(input-keyboard.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("geometry", f.ThemeId, "label", f.Text, "Geometry:", f.NewRow, ""),
		gui.Button("line geometry", f.Text, " Lines `(geometry-lines.go)", f.NewRow, ""),
		gui.Button("shape geometry", f.Text, " Shapes `(geometry-shapes.go)", f.NewRow, ""),
		gui.Button("shape grids geometry", f.Text, " Chunks `(geometry-shapes-grids.go)", f.NewRow, ""),
		gui.Button("shape collision geometry", f.Text, " Collisions `(geometry-shapes-collisions.go)", f.NewRow, ""),
		gui.Button("pathfind around geometry", f.Text, " Pathfinding `(geometry-pathfinding.go)", f.NewRow, ""),
		gui.Button("path following geometry", f.Text, " Path Following `(geometry-path-following.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("data", f.ThemeId, "label", f.Text, "Assets:", f.NewRow, ""),
		gui.Button("tiled scenes", f.Text, " Tiled Scenes `(assets-tiled.go)", f.NewRow, ""),
		gui.Button("default font asset", f.Text, " Default Font `(assets-default.go)", f.NewRow, ""),
		gui.Button("default icons asset", f.Text, " Default Icons `(assets-default.go)", f.NewRow, ""),
		gui.Button("default cursors asset", f.Text, " Default Cursors `(assets-default.go)", f.NewRow, ""),
		gui.Button("default input asset", f.Text, " Default Input `(assets-default.go)", f.NewRow, ""),
		gui.Button("default ui asset", f.Text, " Default User Interface (UI) `(assets-default.go)", f.NewRow, ""),
		gui.Button("default retro atlas asset", f.Text, " Default Retro Atlas `(assets-default.go)", f.NewRow, ""),
		gui.Button("default patterns asset", f.Text, " Default Patterns `(assets-default.go)", f.NewRow, ""),
		gui.Button("default texture asset", f.Text, " Default Texture `(assets-default.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("other", f.ThemeId, "label", f.Text, "Other:", f.NewRow, ""),
		gui.Button("animation sequences", f.Text, " Animation Sequences `(motion-animations.go)", f.NewRow, ""),
		gui.Button("tweens", f.Text, " Tweens `(motion-tweens.go)", f.NewRow, ""),
		gui.Button("flows", f.Text, " Flows `(execution-flows.go)", f.NewRow, ""),
	))
	assets.LoadDefaultFont()
	assets.LoadDefaultSoundsUI()

	hud.Scale = 1.01 // removes tearing artifacts

	var buttons = map[string]func(){
		"minimal graphics": example.MinimalRender,
		"boxes graphics":   example.Boxes,
		"texts graphics":   example.Texts,
		"guis":             example.GUIs,
		//=================================================================
		"mouse input":    example.Mouse,
		"keyboard input": example.Keyboard,
		//=================================================================
		"line geometry":            example.Lines,
		"shape geometry":           example.Shapes,
		"shape grids geometry":     example.ShapesGrids,
		"shape collision geometry": example.Collisions,
		"pathfind around geometry": example.Pathfinding,
		"path following geometry":  example.PathFollowing,
		//=================================================================
		"tiled scenes":              example.Tiled,
		"default font asset":        example.DefaultAssetFont,
		"default icons asset":       example.DefaultAssetIcons,
		"default cursors asset":     example.DefaultAssetCursors,
		"default input asset":       example.DefaultAssetInput,
		"default ui asset":          example.DefaultAssetUI,
		"default retro atlas asset": example.DefaultAssetRetro,
		"default patterns asset":    example.DefaultAssetPatterns,
		"default texture asset":     example.DefaultAssetTexture,
		//=================================================================
		"animation sequences": example.Animations,
		"tweens":              example.Tweens,
		"flows":               example.Flows,
	}

	for window.KeepOpen() {
		window.Title = "pure-game-kit: hub"
		cam.SetScreenAreaToWindow()

		for k, v := range buttons {
			if hud.IsButtonJustClicked(k, cam) {
				window.Title = "pure-game-kit: " + k
				mouse.SetCursor(0)
				v()
			}
		}

		hud.UpdateAndDraw(cam)
	}
}
