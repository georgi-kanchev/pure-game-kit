package graphics

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	tm "pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/utility/time/unit"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (v *View) DrawColor(color uint) {
	var x, y, w, h = v.Bounds()
	v.begin()
	batcher.QueueQuad(x, y, w, h, 0, getColor(color), v.Mask)
	batcher.Draw()
	v.end()
}
func (v *View) DrawGrid(thickness, spacingX, spacingY float32, color uint) {
	if spacingX*v.Zoom < 1 && spacingY*v.Zoom < 1 {
		return // way too dense grid - give up
	}
	v.begin()

	var renderColor = getColor(color)
	var sx, sy, sw, sh = v.area()
	var ulx, uly = v.PointFromScreen(sx, sy)
	var urx, ury = v.PointFromScreen(sx+sw, sy)
	var lrx, lry = v.PointFromScreen(sx+sw, sy+sh)
	var llx, lly = v.PointFromScreen(sx, sy+sh)
	var xs = []float32{ulx, urx, llx, lrx}
	var ys = []float32{uly, ury, lly, lry}
	var minX, maxX = xs[0], xs[0]
	var minY, maxY = ys[0], ys[0]

	for i := 1; i < 4; i++ {
		if xs[i] < minX {
			minX = xs[i]
		}
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] < minY {
			minY = ys[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}

	var left = number.RoundDown(minX/spacingX) * spacingX
	var right = number.RoundUp(maxX/spacingX) * spacingX
	var top = number.RoundDown(minY/spacingY) * spacingY
	var bottom = number.RoundUp(maxY/spacingY) * spacingY

	for x := left; x <= right; x += spacingX {
		var myThickness = thickness
		if number.DivisionRemainder(x, spacingX*10) == 0 {
			myThickness *= 3
		}
		batcher.QueueLine(x, top, x, bottom, myThickness, renderColor, v.Mask)
	}
	for y := top; y <= bottom; y += spacingY {
		var myThickness = thickness
		if number.DivisionRemainder(y, spacingY*10) == 0 {
			myThickness *= 3
		}
		batcher.QueueLine(left, y, right, y, myThickness, renderColor, v.Mask)
	}

	if top <= 0 && bottom >= 0 {
		batcher.QueueLine(left, 0, right, 0, thickness*6, renderColor, v.Mask)
	}
	if left <= 0 && right >= 0 {
		batcher.QueueLine(0, top, 0, bottom, thickness*6, renderColor, v.Mask)
	}

	batcher.Draw()
	v.end()
}

//=================================================================

func (v *View) DrawTextDebug(fps, time, assets, memory bool) {
	if condition.TrueEvery(0.15, ";;;debug") {
		debugStr = ""
		if fps {
			debugStr += text.New("FPS ", int(internal.FPS), " (", int(internal.AverageFPS), ")\n\n")
		}
		if time {
			debugStr += text.New(
				"Time: \n",
				"Running = ", tm.AsClock12(internal.Runtime, ":", unit.Hour|unit.Timer, false), "\n",
				"Frame Busy = ", number.Round(internal.FrameTime*1000, 3), "ms ",
				"(", number.Round((internal.FrameTime/internal.DeltaTime)*100), "%)\n",
				"Frame Idle = ", number.Round((internal.DeltaTime-internal.FrameTime)*1000, 3), "ms ",
				"(", number.Round(((internal.DeltaTime-internal.FrameTime)/internal.DeltaTime)*100), "%)\n",
				"Frame Total = ", number.Round((internal.DeltaTime)*1000, 3), "ms ",
				"\n\n")
		}
		if assets {
			debugStr += text.New("Assets: \n",
				"Textures = ", len(internal.Textures), "\n",
				"Fonts = ", len(internal.Fonts), "\n",
				"Sounds = ", len(internal.Sounds), "\n",
				"Music = ", len(internal.Music), "\n",
				"Tile Data = ", len(internal.TileLayers), "\n\n")
		}
		if memory {
			debugStr += debug.MemoryUsage()
		}
	}

	rl.DrawText(debugStr, 0, 0, int32(40/v.Zoom), rl.White)
}
