package keyboard

import "pure-game-kit/packages/internal"

func Input() string { return internal.Input }

func IsKeyPressed(key int) bool {
	return key >= 0 && key < len(internal.Keys) && internal.Keys[key]
}

func IsKeyHeld(key int) bool {
	return IsKeyPressed(key) && internal.KeysPrev[key]
}

func IsKeyJustPressed(key int) bool {
	return IsKeyPressed(key) && !internal.KeysPrev[key]
}

func IsKeyJustReleased(key int) bool {
	return key >= 0 && key < len(internal.Keys) && !internal.Keys[key] && internal.KeysPrev[key]
}

func IsAnyKeyPressed() bool      { return internal.AnyKey }
func IsAnyKeyJustPressed() bool  { return internal.AnyKey && !internal.AnyKeyPrev }
func IsAnyKeyJustReleased() bool { return !internal.AnyKey && internal.AnyKeyPrev }

func IsComboJustPressed(keys ...int) bool {
	return combo(keys) && IsKeyJustPressed(keys[len(keys)-1])
}

func IsComboHeld(keys ...int) bool {
	return combo(keys) && IsKeyHeld(keys[len(keys)-1])
}

func combo(keys []int) bool {
	if internal.KeyCount != len(keys) {
		return false // pressed key count doesn't match the combo, exit early
	}

	for _, k := range keys {
		if k < 0 || k >= len(internal.Keys) || !internal.Keys[k] {
			return false
		}
	}
	return true
}
