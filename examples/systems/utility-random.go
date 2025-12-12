package example

import (
	"pure-game-kit/debug"
	"pure-game-kit/utility/random"
)

func Randoms() {
	debug.Print("int: ", random.CombineSeeds(10, 20, 30))
	debug.Print("uint32: ", random.CombineSeeds(uint32(5), 99))
	debug.Print("float32: ", random.CombineSeeds(float32(1.2), 3.4))
	debug.Print("byte: ", random.CombineSeeds(byte(23), byte(100)))

	debug.Print("random float: ", random.Range(-5.5, 2.6))
	debug.Print("random int: ", random.Range(100, 200))
	debug.Print("random byte: ", random.Range(byte(20), byte(45)))

	debug.Print("chance 30%: ", random.HasChance(30))
	debug.Print("chance 75%: ", random.HasChance(75))
	debug.Print("chance 0%: ", random.HasChance(0))
	debug.Print("chance 100%: ", random.HasChance(100))

	var arr = []int{1, 2, 3, 4, 5}
	debug.Print("shuffle: ", random.Shuffle(arr))
	debug.Print("pick: ", random.Pick("red", "blue", "green"))

	var arr2 = []string{"apple", "banana", "cherry"}
	debug.Print("pick from: ", random.PickFrom(arr2))
}
