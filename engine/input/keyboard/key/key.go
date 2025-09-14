package key

const (
	Space, Apostrophe, Comma, Minus, Dot, Slash = 32, 39, 44, 45, 46, 47

	Number0, Number1, Number2, Number3, Number4, Number5, Number6, Number7, Number8,
	Number9 = 48, 49, 50, 51, 52, 53, 54, 55, 56, 57

	Semicolon, Equal = 59, 61

	A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y,
	Z = 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90

	LeftBracket, BackSlash, RightBracket, Grave = 91, 92, 93, 96

	Escape, Enter, Tab, Backspace, Insert, Delete, RightArrow, LeftArrow, DownArrow, UpArrow, PageUp, PageDown, Home,
	End, CapsLock, ScrollLock, NumLock, PrintScreen,
	Pause = 256, 257, 258, 259, 260, 261, 262, 263, 264, 265, 266, 267, 268, 269, 280, 281, 282, 283, 284

	F1, F2, F3, F4, F5, F6, F7, F8, F9, F10, F11,
	F12 = 290, 291, 292, 293, 294, 295, 296, 297, 298, 299, 300, 301

	Numpad0, Numpad1, Numpad2, Numpad3, Numpad4, Numpad5, Numpad6, Numpad7, Numpad8,
	Numpad9 = 320, 321, 322, 323, 324, 325, 326, 327, 328, 329

	NumpadDot, NumpadDivide, NumpadMultiply, NumpadSubtract, NumpadAdd, NumpadEnter,
	NumpadEqual = 330, 331, 332, 333, 334, 335, 336
	LeftShift, LeftControl, LeftAlt, LeftSuper, RightShift, RightControl, RightAlt, RightSuper,
	Menu = 340, 341, 342, 343, 344, 345, 346, 347, 348
)

var keyNames = map[any]any{ // this monstrosity is the shortest & simplest way to map keys to their names
	"Space": Space, Space: "Space", "Apostrophe": Apostrophe, Apostrophe: "Apostrophe",
	"Comma": Comma, Comma: "Comma", "Minus": Minus, Minus: "Minus", "Dot": Dot, Dot: "Dot",
	"Slash": Slash, Slash: "Slash",

	"Number0": Number0, Number0: "Number0", "Number1": Number1, Number1: "Number1",
	"Number2": Number2, Number2: "Number2", "Number3": Number3, Number3: "Number3",
	"Number4": Number4, Number4: "Number4", "Number5": Number5, Number5: "Number5",
	"Number6": Number6, Number6: "Number6", "Number7": Number7, Number7: "Number7",
	"Number8": Number8, Number8: "Number8", "Number9": Number9, Number9: "Number9",

	"Semicolon": Semicolon, Semicolon: "Semicolon", "Equal": Equal, Equal: "Equal",

	"A": A, A: "A", "B": B, B: "B", "C": C, C: "C", "D": D, D: "D", "E": E, E: "E", "F": F, F: "F", "G": G, G: "G",
	"H": H, H: "H", "I": I, I: "I", "J": J, J: "J", "K": K, K: "K", "L": L, L: "L", "M": M, M: "M", "N": N, N: "N",
	"O": O, O: "O", "P": P, P: "P", "Q": Q, Q: "Q", "R": R, R: "R", "S": S, S: "S", "T": T, T: "T", "U": U, U: "U",
	"V": V, V: "V", "W": W, W: "W", "X": X, X: "X", "Y": Y, Y: "Y", "Z": Z, Z: "Z",

	"LeftBracket": LeftBracket, LeftBracket: "LeftBracket", "BackSlash": BackSlash, BackSlash: "BackSlash",
	"RightBracket": RightBracket, RightBracket: "RightBracket", "Grave": Grave, Grave: "Grave",

	"Escape": Escape, Escape: "Escape", "Enter": Enter, Enter: "Enter", "Tab": Tab, Tab: "Tab",
	"Backspace": Backspace, Backspace: "Backspace", "Insert": Insert, Insert: "Insert",
	"Delete": Delete, Delete: "Delete", "RightArrow": RightArrow, RightArrow: "RightArrow",
	"LeftArrow": LeftArrow, LeftArrow: "LeftArrow", "DownArrow": DownArrow, DownArrow: "DownArrow",
	"UpArrow": UpArrow, UpArrow: "UpArrow", "PageUp": PageUp, PageUp: "PageUp",
	"PageDown": PageDown, PageDown: "PageDown", "Home": Home, Home: "Home", "End": End, End: "End",
	"CapsLock": CapsLock, CapsLock: "CapsLock", "ScrollLock": ScrollLock, ScrollLock: "ScrollLock",
	"NumLock": NumLock, NumLock: "NumLock", "PrintScreen": PrintScreen, PrintScreen: "PrintScreen",
	"Pause": Pause, Pause: "Pause",

	"F1": F1, F1: "F1", "F2": F2, F2: "F2", "F3": F3, F3: "F3", "F4": F4, F4: "F4", "F5": F5, F5: "F5",
	"F6": F6, F6: "F6", "F7": F7, F7: "F7", "F8": F8, F8: "F8", "F9": F9, F9: "F9", "F10": F10, F10: "F10",
	"F11": F11, F11: "F11", "F12": F12, F12: "F12",

	"Numpad0": Numpad0, Numpad0: "Numpad0", "Numpad1": Numpad1, Numpad1: "Numpad1",
	"Numpad2": Numpad2, Numpad2: "Numpad2", "Numpad3": Numpad3, Numpad3: "Numpad3",
	"Numpad4": Numpad4, Numpad4: "Numpad4", "Numpad5": Numpad5, Numpad5: "Numpad5",
	"Numpad6": Numpad6, Numpad6: "Numpad6", "Numpad7": Numpad7, Numpad7: "Numpad7",
	"Numpad8": Numpad8, Numpad8: "Numpad8", "Numpad9": Numpad9, Numpad9: "Numpad9",

	"NumpadDot": NumpadDot, NumpadDot: "NumpadDot", "NumpadDivide": NumpadDivide, NumpadDivide: "NumpadDivide",
	"NumpadMultiply": NumpadMultiply, NumpadMultiply: "NumpadMultiply",
	"NumpadSubtract": NumpadSubtract, NumpadSubtract: "NumpadSubtract",
	"NumpadAdd": NumpadAdd, NumpadAdd: "NumpadAdd", "NumpadEnter": NumpadEnter, NumpadEnter: "NumpadEnter",
	"NumpadEqual": NumpadEqual, NumpadEqual: "NumpadEqual",

	"LeftShift": LeftShift, LeftShift: "LeftShift", "LeftControl": LeftControl, LeftControl: "LeftControl",
	"LeftAlt": LeftAlt, LeftAlt: "LeftAlt", "LeftSuper": LeftSuper, LeftSuper: "LeftSuper",
	"RightShift": RightShift, RightShift: "RightShift", "RightControl": RightControl, RightControl: "RightControl",
	"RightAlt": RightAlt, RightAlt: "RightAlt", "RightSuper": RightSuper, RightSuper: "RightSuper",
	"Menu": Menu, Menu: "Menu",
}

func ToName(key int) string {
	var value, has = keyNames[key]
	if has {
		return value.(string)
	}
	return ""
}

func FromName(name string) int {
	var value, has = keyNames[name]
	if has {
		return value.(int)
	}
	return 0
}
