// Helper functions for bit-masks.
//
// They are useful for storing up to 64 flags (bool values) in a single integer
// where each bit represents each flag (on/off).
package flag

import "golang.org/x/exp/constraints"

func IsOn[T constraints.Integer](allFlags, flag T) bool {
	return allFlags&flag != 0
}

func TurnOn[T constraints.Integer](allFlags, flag T) T {
	return allFlags | flag
}

func Toggle[T constraints.Integer](allFlags, flag T) T {
	return allFlags ^ flag
}

func TurnOff[T constraints.Integer](allFlags, flag T) T {
	return allFlags &^ flag
}

func FromBit[T constraints.Integer](bitPosition int) T {
	return 1 << bitPosition
}
