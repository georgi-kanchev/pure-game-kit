package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/engine"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/random"
)

func main() {
	engine.Initialize("pure-game-kit", 60, 120, false, false)

	var view = graphics.NewView(1)
	// var path = "tools/sdf-font-generator/results/Roboto-Bold."
	// var font = assets.LoadFont2(path+"png", path+"xml")
	// obj.TextFont = font
	// obj.Text = "Hello, World!"

	var imgId = assets.LoadImage("examples/data/flail.PNG")
	var objs = make([]*graphics.Object, 1024)

	for i := range objs {
		var obj = graphics.NewObject(float32(random.Range(0, 1920)), float32(random.Range(0, 1080)))
		objs[i] = &obj
		objs[i].ImageId = imgId
		objs[i].Width, objs[i].Height = 32, 32
	}

	// assets.LoadDefaultFont()
	engine.Run(func() {
		view.DrawObjects(objs...)
	})
}
