package utility

type F struct {
	X, Y float32
}

type I struct {
	X, Y int32
}

type Pair[V1, V2 any] struct {
	Item1 V1
	Item2 V2
}
type Pack3[V1, V2, V3 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
}
type Pack4[V1, V2, V3, V4 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
}
type Pack5[V1, V2, V3, V4, V5 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
	Item5 V5
}
type Pack6[V1, V2, V3, V4, V5, V6 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
	Item5 V5
	Item6 V6
}
type Pack7[V1, V2, V3, V4, V5, V6, V7 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
	Item5 V5
	Item6 V6
	Item7 V7
}
type Pack8[V1, V2, V3, V4, V5, V6, V7, V8 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
	Item5 V5
	Item6 V6
	Item7 V7
	Item8 V8
}
type Pack9[V1, V2, V3, V4, V5, V6, V7, V8, V9 any] struct {
	Item1 V1
	Item2 V2
	Item3 V3
	Item4 V4
	Item5 V5
	Item6 V6
	Item7 V7
	Item8 V8
	Item9 V9
}
