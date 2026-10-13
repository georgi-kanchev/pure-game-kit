package particles

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/utility/collection"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/point"
	"pure-game-kit/packages/utility/random"
)

type Particle struct {
	X, Y,
	Age, Angle, Scale,
	VelocityX, VelocityY float32
	Color          uint
	Id, FrameIndex int
	AssetId        assets.ImageId
	CustomData     map[string]any
}
type System struct{ emitters []*emitter }

func NewSystem() *System { return &System{emitters: make([]*emitter, 0, 8)} }

func (s *System) EmitFromPoint(amount int, x, y float32, update func(*Particle) (alive bool)) {
	var emitter = newEmitter(update)
	s.emitters = append(s.emitters, emitter)
	for i := range amount {
		emitter.particles = append(emitter.particles, newParticle(i, x, y))
	}
}
func (s *System) EmitFromLine(amount int, ax, ay, bx, by float32, update func(*Particle) (alive bool)) {
	var emitter = newEmitter(update)
	s.emitters = append(s.emitters, emitter)
	for i := range amount {
		var x, y = point.MoveByPercent(ax, ay, bx, by, random.Range[float32](0, 100))
		emitter.particles = append(emitter.particles, newParticle(i, x, y))
	}
}

func (s *System) Update() {
	for i := len(s.emitters) - 1; i >= 0; i-- { // iterate in reverse to not affect indices when removing
		var e = s.emitters[i]
		for j := len(e.particles) - 1; j >= 0; j-- { // same here
			if e.update != nil && !e.update(e.particles[j]) {
				e.particles = collection.RemoveUnordered(e.particles, j)
			}
		}
		if len(e.particles) == 0 {
			s.emitters = collection.RemoveUnordered(s.emitters, i)
		}
	}
}

// private ========================================================

type emitter struct {
	particles []*Particle
	update    func(*Particle) bool
}

func newParticle(id int, x, y float32) *Particle {
	return &Particle{Id: id, X: x, Y: y, Scale: 1, Color: palette.White, CustomData: make(map[string]any)}
}
func newEmitter(update func(*Particle) (alive bool)) *emitter {
	return &emitter{update: update, particles: make([]*Particle, 0, 128)}
}
