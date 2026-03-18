package example

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/random"
	"pure-game-kit/window"
)

func Batches() {
	var cam = graphics.NewCamera(1)
	var count = 20_000
	var points = make([][2]float32, 0, count*4)

	condition.CallAfter(0, func() {
		var w, h = cam.Size()

		for range count {
			var cx, cy = random.Range(-w/2, w/2), random.Range(-h/2, h/2)
			var baseRotation = random.Range[float32](0, 6.28)

			for j := range 3 {
				var sector = (float32(j) * 2.0 * 3.14159) / 3.0
				var angle = baseRotation + sector + random.Range[float32](-0.5, 0.5)
				var dist = random.Range[float32](10, 50)
				var x = cx + number.Cosine(angle)*dist
				var y = cy + number.Sine(angle)*dist
				points = collection.Add(points, [2]float32{x, y})
			}

			points = append(points, [2]float32{number.NaN(), number.NaN()})
		}
	})

	for window.KeepOpen() {
		cam.DrawShapes(palette.Red, points...)
		cam.DrawTextFPS()
	}
}
