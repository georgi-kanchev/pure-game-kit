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
	// example.Randoms()
	// example.StorageBinary()
	example.StorageYAML()
	example.StorageXML()
	example.StorageJSON()

	var cam = graphics.NewCamera(1)
	var _, _, box = assets.LoadDefaultAtlasUI()
	var hud = gui.NewFromXMLs(gui.NewElementsXML(
		gui.Container("themes", "", "", "", ""),
		gui.Theme("label", f.Color, "0 0 0 0", f.Width, d.OwnerWidth+"-40", f.Height, "50", f.GapY, "20",
			f.TextAlignmentX, "0", f.TextAlignmentY, "0.5", f.TextColor, "0 0 0 255",
			f.TextLineHeight, "40"),
		gui.Theme("button", f.Color, "220 220 220 255", f.Width, d.OwnerWidth+"-40",
			f.Height, "55", f.GapX, "20", f.GapY, "0",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[2], f.TextAlignmentX, "0", f.TextAlignmentY, "0.3", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "35", f.ButtonThemeIdHover, "button-hover", f.ButtonThemeIdPress, "button-press"),
		gui.Theme("button-hover", f.Color, "255 255 255 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40",
			f.BoxEdgeBottom, "40", f.AssetId, box[5], f.TextAlignmentX, "0", f.TextAlignmentY, "0.3",
			f.TextColor, "127 127 127 255", f.TextLineHeight, "35", f.GapX, "20", f.GapY, "0"),
		gui.Theme("button-press", f.Color, "200 200 200 255", f.Width, "300", f.Height, "100",
			f.BoxEdgeLeft, "40", f.BoxEdgeRight, "40", f.BoxEdgeTop, "40", f.BoxEdgeBottom, "40",
			f.AssetId, box[4], f.TextAlignmentX, "0", f.TextAlignmentY, "0.6", f.TextColor, "80 80 80 255",
			f.TextLineHeight, "35", f.GapX, "20", f.GapY, "0"),
		// ======================================================
		gui.Container("title", d.CameraCenterX+"-400", d.CameraTopY+"+20", "800", "150",
			f.ThemeId, "button", f.GapY, "20"),
		gui.Visual("background", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		// ======================================================
		gui.Visual("description", f.Text, "pure-game-kit - simple 2D game engine\nExamples",
			f.AssetId, "", f.Color, "0 0 0 0", f.Width, d.OwnerWidth+"-30", f.Height, d.OwnerHeight+"-50",
			f.OffsetX, "15", f.TextAlignmentX, "0.5", f.TextAlignmentY, "0.5", f.TextColor, "255 255 255 255",
			f.TextLineHeight, "50", f.TextColorOutline, "0 0 0 255", f.TextThicknessOutline, "0.8"),
		gui.Container("menu", d.CameraCenterX+"-400", d.CameraTopY+"+150", "800", d.CameraHeight+"-180",
			f.ThemeId, "button", f.GapX, "20", f.GapY, "20"),
		gui.Visual("bg", f.FillContainer, "", f.AssetId, box[8], f.Color, "200 200 200 255"),
		//=================================================================
		gui.Visual("gfx", f.ThemeId, "label", f.Text, "Graphics:", f.GapY, "0", f.NewRow, ""),
		gui.Button("minimal graphics", f.Text, " Minimal Render (graphics-minimal.go)", f.NewRow, ""),
		gui.Button("boxes", f.Text, " Boxes (graphics-boxes.go)", f.NewRow, ""),
		gui.Button("texts", f.Text, " Texts (graphics-texts.go)", f.NewRow, ""),
		gui.Button("effects", f.Text, " Effects (graphics-effects.go)", f.NewRow, ""),
		gui.Button("guis", f.Text, " Graphical User Interfaces (GUIs) (guis.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("input", f.ThemeId, "label", f.Text, "Input:", f.NewRow, ""),
		gui.Button("mouse input", f.Text, " Mouse (input-mouse.go)", f.NewRow, ""),
		gui.Button("keyboard input", f.Text, " Keboard (input-keyboard.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("geometry", f.ThemeId, "label", f.Text, "Geometry:", f.NewRow, ""),
		gui.Button("line geometry", f.Text, " Lines (geometry-lines.go)", f.NewRow, ""),
		gui.Button("shape geometry", f.Text, " Shapes (geometry-shapes.go)", f.NewRow, ""),
		gui.Button("shape grids geometry", f.Text, " Chunks (geometry-shapes-grids.go)", f.NewRow, ""),
		gui.Button("shape collision geometry", f.Text, " Collisions (geometry-shapes-collisions.go)", f.NewRow, ""),
		gui.Button("pathfind around geometry", f.Text, " Pathfinding (geometry-pathfinding.go)", f.NewRow, ""),
		gui.Button("path following geometry", f.Text, " Path Following (geometry-path-following.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("motion", f.ThemeId, "label", f.Text, "Motion:", f.NewRow, ""),
		gui.Button("animation sequences", f.Text, " Animation Sequences (motion-animations.go)", f.NewRow, ""),
		gui.Button("tweens", f.Text, " Tweens (motion-tweens.go)", f.NewRow, ""),
		gui.Button("particles", f.Text, " Particles (motion-particles.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("execution", f.ThemeId, "label", f.Text, "Execution:", f.NewRow, ""),
		gui.Button("state machines", f.Text, " State Machines (execution-state-machines.go)", f.NewRow, ""),
		//=================================================================
		gui.Visual("data", f.ThemeId, "label", f.Text, "Assets:", f.NewRow, ""),
		gui.Button("tiled scenes", f.Text, " Tiled Scenes (assets-tiled.go)", f.NewRow, ""),
		gui.Button("default font asset", f.Text, " Default Font (assets-default.go)", f.NewRow, ""),
		gui.Button("default icons asset", f.Text, " Default Icons (assets-default.go)", f.NewRow, ""),
		gui.Button("default cursors asset", f.Text, " Default Cursors (assets-default.go)", f.NewRow, ""),
		gui.Button("default input asset", f.Text, " Default Input (assets-default.go)", f.NewRow, ""),
		gui.Button("default ui asset", f.Text, " Default User Interface (UI) (assets-default.go)", f.NewRow, ""),
		gui.Button("default retro atlas asset", f.Text, " Default Retro Atlas (assets-default.go)", f.NewRow, ""),
		gui.Button("default patterns asset", f.Text, " Default Patterns (assets-default.go)", f.NewRow, ""),
		gui.Button("default texture asset", f.Text, " Default Texture (assets-default.go)", f.NewRow, ""),
	))
	assets.LoadDefaultFont()
	assets.LoadDefaultSoundsUI()

	hud.Scale = 2.01 // removes tearing artifacts

	var buttons = map[string]func(){
		"minimal graphics": example.MinimalRender,
		"boxes":            example.Boxes,
		"texts":            example.Texts,
		"guis":             example.GUIs,
		"effects":          example.Effects,
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
		"animation sequences": example.Animations,
		"tweens":              example.Tweens,
		"particles":           example.Particles,
		//=================================================================
		"state machines": example.StateMachines,
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
	}

	for window.KeepOpen() {
		window.Title = "pure-game-kit: hub"
		cam.SetScreenAreaToWindow()

		hud.UpdateAndDraw(cam)

		for k, v := range buttons {
			if hud.IsButtonJustClicked(k) {
				window.Title = "pure-game-kit: " + k
				mouse.SetCursor(0)
				v()
			}
		}
	}
}
