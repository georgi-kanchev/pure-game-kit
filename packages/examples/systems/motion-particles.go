package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/input/mouse/button"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/random"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

func Particles() {
	var cam = graphics.NewCamera(1)
	var particles = motion.NewParticleSystem(func(p *motion.Particle) bool {
		if p.Age == 0 {
			p.VelocityX, p.VelocityY = random.Range[float32](-3, 3), random.Range[float32](0, 1)
			p.CustomData["bounces"] = float32(0)
			p.Color = palette.Cyan
		}
		p.VelocityY += time.FrameDelta() * 9.8 // gravity
		p.X, p.Y = p.X+p.VelocityX*time.FrameDelta()*60, p.Y+p.VelocityY*time.FrameDelta()*60

		if p.Y > -p.Scale/2 {
			var bounces = p.CustomData["bounces"].(float32) + 1
			p.VelocityX, p.Y = p.VelocityX/p.Age, -p.Scale/2

			if bounces < 5 {
				p.VelocityY = -random.Range[float32](5, 6) / bounces
				p.CustomData["bounces"] = bounces
			}
		}

		p.Age += time.FrameDelta()
		p.Scale = (6 - p.Age) * 5
		cam.DrawCircle(p.X, p.Y-5, p.Scale, 8, p.Color)

		if p.Age > 5 {
			p.Color = color.FadeOut(palette.Cyan, number.Map(p.Age, 5, 6, 0, 1))
		}
		return p.Age < 6
	})

	for window.KeepOpen() {
		var clx, cly = cam.PointFromEdge(0, 0.5)
		var cw, ch = cam.Size()
		cam.DrawQuad(clx, cly, cw, ch, 0, palette.DarkGray)

		particles.Update()

		if mouse.IsButtonJustPressed(button.Left) {
			var mx, my = cam.MousePosition()
			particles.EmitFromLine(30, mx-100, my, mx+100, my)
		}
		cam.DrawTextDebug(true, true, true, true)
	}
}
