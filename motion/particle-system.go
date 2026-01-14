package motion

import (
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/random"
)

type ParticleSystem struct {
	particles []*Particle
	update    func(*Particle) bool
}

func NewParticleSystem(update func(*Particle) (alive bool)) *ParticleSystem {
	return &ParticleSystem{update: update, particles: make([]*Particle, 0, 128)}
}

//=================================================================

func (ps *ParticleSystem) EmitFromPoint(amount int, x, y float32) {
	for i := range amount {
		ps.particles = append(ps.particles, newParticle(i, x, y))
	}
}
func (ps *ParticleSystem) EmitFromLine(amount int, ax, ay, bx, by float32) {
	for i := range amount {
		var x, y = point.MoveByPercent(ax, ay, bx, by, random.Range[float32](0, 100))
		ps.particles = append(ps.particles, newParticle(i, x, y))
	}
}

//=================================================================

func (ps *ParticleSystem) Update() {
	for i := len(ps.particles) - 1; i >= 0; i-- { // iterate in reverse to not affect indices when removing
		if !ps.update(ps.particles[i]) {
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}
}
