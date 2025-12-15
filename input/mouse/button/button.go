// Contains all mouse buttons as constants, as well as a way to convert them to & from their name.
package button

const Left, Right, Middle, Extra1, Extra2 = 0, 1, 2, 3, 4

func ToName(button int) string {
	var value, has = buttonNames[button]
	if has {
		return value.(string)
	}
	return ""
}
func FromName(name string) int {
	var value, has = buttonNames[name]
	if has {
		return value.(int)
	}
	return -1
}

//=================================================================
// private

var buttonNames = map[any]any{ // this monstrosity is the shortest & simplest way to map buttons to their names
	"Left": Left, Left: "Left", "Right": Right, Right: "Right", "Middle": Middle, Middle: "Middle",
	"Extra1": Extra1, Extra1: "Extra1", "Extra2": Extra2, Extra2: "Extra2",
}
