// Contains all keyboard keys as constants, as well as a way to convert them to & from their name.
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

func ToName(key int) string {
	return keyToName[key]
}
func FromName(name string) int {
	return nameToKey[name]
}

//=================================================================
// private

var nameToKey = map[string]int{
	"Space": Space, "Apostrophe": Apostrophe, "Comma": Comma, "Minus": Minus, "Dot": Dot, "Slash": Slash,

	"Number0": Number0, "Number1": Number1, "Number2": Number2, "Number3": Number3, "Number4": Number4,
	"Number5": Number5, "Number6": Number6, "Number7": Number7, "Number8": Number8, "Number9": Number9,

	"Semicolon": Semicolon, "Equal": Equal,

	"A": A, "B": B, "C": C, "D": D, "E": E, "F": F, "G": G, "H": H, "I": I, "J": J, "K": K, "L": L, "M": M,
	"N": N, "O": O, "P": P, "Q": Q, "R": R, "S": S, "T": T, "U": U, "V": V, "W": W, "X": X, "Y": Y, "Z": Z,

	"LeftBracket": LeftBracket, "BackSlash": BackSlash, "RightBracket": RightBracket, "Grave": Grave,

	"Escape": Escape, "Enter": Enter, "Tab": Tab, "Backspace": Backspace, "Insert": Insert,
	"Delete": Delete, "RightArrow": RightArrow, "LeftArrow": LeftArrow, "DownArrow": DownArrow,
	"UpArrow": UpArrow, "PageUp": PageUp, "PageDown": PageDown, "Home": Home, "End": End,
	"CapsLock": CapsLock, "ScrollLock": ScrollLock, "NumLock": NumLock, "PrintScreen": PrintScreen, "Pause": Pause,

	"F1": F1, "F2": F2, "F3": F3, "F4": F4, "F5": F5, "F6": F6,
	"F7": F7, "F8": F8, "F9": F9, "F10": F10, "F11": F11, "F12": F12,

	"Numpad0": Numpad0, "Numpad1": Numpad1, "Numpad2": Numpad2, "Numpad3": Numpad3, "Numpad4": Numpad4,
	"Numpad5": Numpad5, "Numpad6": Numpad6, "Numpad7": Numpad7, "Numpad8": Numpad8, "Numpad9": Numpad9,

	"NumpadDot": NumpadDot, "NumpadDivide": NumpadDivide, "NumpadMultiply": NumpadMultiply,
	"NumpadSubtract": NumpadSubtract, "NumpadAdd": NumpadAdd, "NumpadEnter": NumpadEnter, "NumpadEqual": NumpadEqual,

	"LeftShift": LeftShift, "LeftControl": LeftControl, "LeftAlt": LeftAlt, "LeftSuper": LeftSuper,
	"RightShift": RightShift, "RightControl": RightControl, "RightAlt": RightAlt, "RightSuper": RightSuper,
	"Menu": Menu,
}
var keyToName = map[int]string{
	Space: "Space", Apostrophe: "Apostrophe", Comma: "Comma", Minus: "Minus", Dot: "Dot", Slash: "Slash",

	Number0: "Number0", Number1: "Number1", Number2: "Number2", Number3: "Number3", Number4: "Number4",
	Number5: "Number5", Number6: "Number6", Number7: "Number7", Number8: "Number8", Number9: "Number9",

	Semicolon: "Semicolon", Equal: "Equal",

	A: "A", B: "B", C: "C", D: "D", E: "E", F: "F", G: "G", H: "H", I: "I", J: "J", K: "K", L: "L", M: "M",
	N: "N", O: "O", P: "P", Q: "Q", R: "R", S: "S", T: "T", U: "U", V: "V", W: "W", X: "X", Y: "Y", Z: "Z",

	LeftBracket: "LeftBracket", BackSlash: "BackSlash", RightBracket: "RightBracket", Grave: "Grave",

	Escape: "Escape", Enter: "Enter", Tab: "Tab", Backspace: "Backspace", Insert: "Insert",
	Delete: "Delete", RightArrow: "RightArrow", LeftArrow: "LeftArrow", DownArrow: "DownArrow",
	UpArrow: "UpArrow", PageUp: "PageUp", PageDown: "PageDown", Home: "Home", End: "End",
	CapsLock: "CapsLock", ScrollLock: "ScrollLock", NumLock: "NumLock", PrintScreen: "PrintScreen", Pause: "Pause",

	F1: "F1", F2: "F2", F3: "F3", F4: "F4", F5: "F5", F6: "F6",
	F7: "F7", F8: "F8", F9: "F9", F10: "F10", F11: "F11", F12: "F12",

	Numpad0: "Numpad0", Numpad1: "Numpad1", Numpad2: "Numpad2", Numpad3: "Numpad3", Numpad4: "Numpad4",
	Numpad5: "Numpad5", Numpad6: "Numpad6", Numpad7: "Numpad7", Numpad8: "Numpad8", Numpad9: "Numpad9",

	NumpadDot: "NumpadDot", NumpadDivide: "NumpadDivide", NumpadMultiply: "NumpadMultiply",
	NumpadSubtract: "NumpadSubtract", NumpadAdd: "NumpadAdd", NumpadEnter: "NumpadEnter", NumpadEqual: "NumpadEqual",

	LeftShift: "LeftShift", LeftControl: "LeftControl", LeftAlt: "LeftAlt", LeftSuper: "LeftSuper",
	RightShift: "RightShift", RightControl: "RightControl", RightAlt: "RightAlt", RightSuper: "RightSuper",
	Menu: "Menu",
}
