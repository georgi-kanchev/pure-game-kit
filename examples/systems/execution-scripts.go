package example

import (
	"fmt"
	"pure-kit/engine/execution/script"
	"pure-kit/engine/utility/text"
)

func Scripts() {
	var scr = script.New()

	scr.AddFunction("ToLower", func(value string) string {
		return text.LowerCase(value)
	})
	scr.ExecuteCode(`
function Greet(name)
    print("Hello, " .. name .. "!")
end

function Add(a, b)
    return a + b
end

function ToLowerAndHi(value)
	return ToLower(value) .. " and Hi!"
end

function IsEven(num)
    return num % 2 == 0
end`)
	scr.ExecuteFunction("Greet", "tosho")
	fmt.Printf("Add 3 + 4: %v\n", scr.ExecuteFunction("Add", 3, 4))
	fmt.Printf("2 is even: %v\n", scr.ExecuteFunction("IsEven", 2))
	fmt.Printf("lowercase: %v\n", scr.ExecuteFunction("ToLower", "Hello, World!"))
	fmt.Printf("ToLowerAndHi: %v\n", scr.ExecuteFunction("ToLowerAndHi", "TESTING"))
}
