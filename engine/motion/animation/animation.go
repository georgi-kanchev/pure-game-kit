package animation

import "math"

type Sequence[T any] struct {
	Items                []T
	ItemsPerSecond       float32
	IsPlaying, IsLooping bool

	time float32
}

func NewSequence[T any](itemsPerSecond float32, items ...T) Sequence[T] {
	return Sequence[T]{Items: items, ItemsPerSecond: itemsPerSecond}
}
func NewSequenceLoop[T any](itemsPerSecond float32, items ...T) Sequence[T] {
	return Sequence[T]{Items: items, ItemsPerSecond: itemsPerSecond, IsLooping: true}
}

func (sequence *Sequence[T]) Advance(deltaTime float32) (item *T, index int) {
	if !sequence.IsPlaying {
		return &sequence.Items[0], 0
	}

	sequence.time += deltaTime

	var count = float64(len(sequence.Items))
	var duration = count / float64(sequence.ItemsPerSecond)
	var progress = float64(sequence.time) / duration

	if progress >= 1 {
		if !sequence.IsLooping {
			sequence.IsPlaying = false
		}

		sequence.time = 0
		progress = 0
	}

	var i = int(math.Min(progress*count, count-1))
	return &sequence.Items[i], i
}
