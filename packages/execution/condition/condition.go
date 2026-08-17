// Utility functions that control the code flow based on some condition or time, as well as
// acting as shortcuts for simple if/else blocks.
package condition

import "pure-game-kit/packages/internal"

func JustTurnedTrue(condition bool, key int) bool {
	var prev = trueOnce[key]
	trueOnce[key] = condition
	return !prev && condition
}
func TrueEvery(seconds float32, key int) bool {
	var start, has = trueEvery[key]
	if !has || internal.Runtime > start+seconds {
		trueEvery[key] = internal.Runtime
		return has
	}
	return false
}

// private ========================================================

var trueOnce = make(map[int]bool, 128)
var trueEvery = make(map[int]float32, 128)
