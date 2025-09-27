package example

import (
	"fmt"
	"pure-kit/engine/utility/flag"
)

const ( // ....DCBA
	A byte = 1 << iota // ....---+
	B                  // ....--+-
	C                  // ....-+--
	D                  // ....+---
)

func Flags() {
	var value byte                // ....----
	value = flag.TurnOn(value, B) // ....--+-
	value = flag.TurnOn(value, A) // ....--++

	fmt.Printf("After turning on B and A: %04b\n", value) // 0011
	fmt.Println("A is on?", flag.IsOn(value, A))          // true
	fmt.Println("C is on?", flag.IsOn(value, C))          // false

	value = flag.TurnOff(value, A) // ....--+-

	fmt.Printf("After turning off A: %04b\n", value) // 0010

	value = flag.Toggle(value, C) // ....-++-
	value = flag.Toggle(value, B) // ....-+--

	fmt.Printf("After toggling toggling C and B: %04b\n", value) // 0100
}
