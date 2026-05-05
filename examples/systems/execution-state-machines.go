package example

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/execution/flow"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/text"
	"pure-game-kit/packages/utility/time"
	"pure-game-kit/packages/window"
)

const jumpStrength, gravityStrength, moveSpeed, size = 600.0, 600.0, 500.0, 100.0

type player struct {
	x, y     float32
	behavior *flow.StateMachine
	state    string
}

func StateMachines() {
	var view = graphics.NewView(1)
	var player = player{behavior: flow.NewStateMachine()}
	player.behavior.GoToState(player.ground)

	for window.KeepOpen() {
		var clx, cly = view.PointFromEdge(0, 0.5)
		var cw, ch = view.Size()
		view.DrawQuad(clx, cly, cw, ch, 0, palette.DarkGray)

		player.behavior.UpdateCurrentState()

		var height = condition.If(player.state == "crouch", float32(size/2), size)
		view.DrawQuad(player.x-size/2, player.y-height, size, height, 0, palette.White)

		var tlx, tly = view.PointFromEdge(0, 0)
		view.DrawText(text.New("State: ", player.state), tlx, tly, 100)
	}
}

func (p *player) move(speedMultiplier float32) {
	var dt = time.FrameDelta() * moveSpeed * speedMultiplier
	if keyboard.IsKeyPressed(key.A) {
		p.x -= dt
	}
	if keyboard.IsKeyPressed(key.D) {
		p.x += dt
	}
}
func (p *player) jump(strengthMultiplier float32) {
	p.move(0.8)

	var timer = p.behavior.StateTimer()
	var velocity = number.Map(timer, 0, 1, jumpStrength*strengthMultiplier, 0)
	p.y -= velocity * time.FrameDelta()

	if timer > 1 {
		p.behavior.GoToState(p.fall)
	}
}

//=================================================================
// states

func (p *player) ground() {
	p.state = "ground"
	p.move(1)
	p.y = 0

	if keyboard.IsKeyPressed(key.W) {
		p.behavior.GoToState(p.normalJump)
	}
	if keyboard.IsKeyPressed(key.S) {
		p.behavior.GoToState(p.crouch)
	}
}
func (p *player) crouch() {
	p.state = "crouch"
	p.move(0.2)

	if keyboard.IsKeyPressed(key.W) {
		p.behavior.GoToState(p.superJump)
	}
	if !keyboard.IsKeyPressed(key.S) {
		p.behavior.GoToState(p.ground)
	}
}
func (p *player) normalJump() {
	p.state = "jump"
	p.jump(1)
}
func (p *player) superJump() {
	p.state = "super jump"
	p.jump(1.5)
}
func (p *player) fall() {
	p.state = "fall"
	p.move(0.6)

	var velocity = number.Map(p.behavior.StateTimer(), 0, 1, 0, gravityStrength)
	p.y += velocity * time.FrameDelta()

	if p.y > 0 {
		p.behavior.GoToState(p.ground)
	}
}
