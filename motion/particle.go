package motion

import "pure-game-kit/utility/color/palette"

type Particle struct {
	X, Y,
	Age, Angle, Scale,
	VelocityX, VelocityY float32
	Color          uint
	Id, FrameIndex int
	AssetId        string
	CustomData     map[string]any
}

func newParticle(id int, x, y float32) *Particle {
	return &Particle{
		Id: id, X: x, Y: y, Scale: 1, Color: palette.White,
		CustomData: make(map[string]any),
	}
}
