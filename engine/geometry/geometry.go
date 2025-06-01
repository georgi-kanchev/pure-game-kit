package geometry

type Point struct{ X, Y float32 }
type Cell struct{ X, Y int32 }
type Area struct{ X, Y, Width, Height float32 }
type Chunk struct{ X, Y, Width, Height int32 }
type Line struct{ Ax, Ay, Bx, By float32 }
