package unit

const (
	Day = 1 << iota
	Hour
	Minute
	Second
	Millisecond
	All   = Day | Hour | Minute | Second | Millisecond
	Clock = Hour | Minute | Second
	Timer = Minute | Second | Millisecond
)
