package keyboard

import i "pure-game-kit/packages/internal"

func Input() string { return i.Input }

func IsKeyHeld(key int, delay float32) bool { return i.Keys[key] && i.KeyDurs[key] > delay }
func IsKeyPressed(key int) bool             { return i.Keys[key] }
func IsKeyJustPressed(key int) bool         { return i.Keys[key] && !i.KeysPrev[key] }
func IsKeyJustReleased(key int) bool        { return !i.Keys[key] && i.KeysPrev[key] }

func IsAnyKeyPressed() bool      { return i.AnyKey }
func IsAnyKeyJustPressed() bool  { return i.AnyKey && !i.AnyKeyPrev }
func IsAnyKeyJustReleased() bool { return !i.AnyKey && i.AnyKeyPrev }

func IsComboJustPressed(keys ...int) bool { return combo(keys) && IsKeyJustPressed(keys[len(keys)-1]) }
func IsComboHeld(delay float32, keys ...int) bool {
	return combo(keys) && IsKeyHeld(keys[len(keys)-1], delay)
}

//=================================================================

func combo(keys []int) bool {
	if i.KeyCount != len(keys) {
		return false // pressed key count doesn't match the combo, exit early
	}

	for _, k := range keys {
		if k < 0 || k >= len(i.Keys) || !i.Keys[k] {
			return false
		}
	}
	return true
}
