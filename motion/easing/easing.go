/*
Easing functions. Similarly to math's Sine and Cosine, they interpolate a number according
to a progress between (but not limited to) 0 and 1.
Useful for describing a smooth movement or path for a numeric value.
*/
package easing

import "pure-game-kit/utility/number"

func Linear(progress float32) float32 {
	return progress
}

func SineIn(progress float32) float32 {
	return 1 - number.Cosine(progress*pi/2)
}
func SineOut(progress float32) float32 {
	return number.Sine(progress * pi / 2)
}
func SineInOut(progress float32) float32 {
	return -0.5 * (number.Cosine(pi*progress) - 1)
}

func CircIn(progress float32) float32 {
	return 1 - number.SquareRoot(1-number.Power(progress, 2))
}
func CircOut(progress float32) float32 {
	return number.SquareRoot(1 - number.Power(progress-1, 2))
}
func CircInOut(progress float32) float32 {
	var p2 = progress * 2
	if p2 < 1 {
		return -0.5 * (number.SquareRoot(1-number.Power(p2, 2)) - 1)
	}
	p2 -= 2
	return 0.5 * (number.SquareRoot(1-number.Power(p2, 2)) + 1)
}

func CubicIn(progress float32) float32 {
	return progress * progress * progress
}
func CubicOut(progress float32) float32 {
	progress--
	return progress*progress*progress + 1
}
func CubicInOut(progress float32) float32 {
	p2 := progress * 2
	if p2 < 1 {
		return 0.5 * p2 * p2 * p2
	}
	p2 -= 2
	return 0.5 * (p2*p2*p2 + 2)
}

func QuadIn(progress float32) float32 {
	return progress * progress
}
func QuadOut(progress float32) float32 {
	return -progress * (progress - 2)
}
func QuadInOut(progress float32) float32 {
	var p2 = progress * 2
	if p2 < 1 {
		return 0.5 * p2 * p2
	}
	p2--
	return -0.5 * (p2*(p2-2) - 1)
}

func ExpoIn(progress float32) float32 {
	if progress == 0 {
		return 0
	}
	return number.Power(2, 10*(progress-1))
}
func ExpoOut(progress float32) float32 {
	if progress == 1 {
		return 1
	}
	return 1 - number.Power(2, -10*progress)
}
func ExpoInOut(progress float32) float32 {
	if progress == 0 {
		return 0
	}
	if progress == 1 {
		return 1
	}
	p2 := progress * 2
	if p2 < 1 {
		return 0.5 * number.Power(2, 10*(p2-1))
	}
	return 0.5 * (2 - number.Power(2, -10*(p2-1)))
}

func BackIn(progress float32) float32 {
	const s float32 = 1.70158
	return progress * progress * ((s+1)*progress - s)
}
func BackOut(progress float32) float32 {
	const s float32 = 1.70158
	progress--
	return progress*progress*((s+1)*progress+s) + 1
}
func BackInOut(progress float32) float32 {
	const s float32 = 1.70158 * 1.525
	var p2 = progress * 2
	if p2 < 1 {
		return 0.5 * (p2 * p2 * ((s+1)*p2 - s))
	}
	p2 -= 2
	return 0.5 * (p2*p2*((s+1)*p2+s) + 2)
}

func BounceIn(progress float32) float32 {
	return 1 - BounceOut(1-progress)
}
func BounceOut(progress float32) float32 {
	if progress < 1/2.75 {
		return 7.5625 * progress * progress
	} else if progress < 2/2.75 {
		progress -= 1.5 / 2.75
		return 7.5625*progress*progress + 0.75
	} else if progress < 2.5/2.75 {
		progress -= 2.25 / 2.75
		return 7.5625*progress*progress + 0.9375
	}
	progress -= 2.625 / 2.75
	return 7.5625*progress*progress + 0.984375
}
func BounceInOut(progress float32) float32 {
	if progress < 0.5 {
		return BounceIn(progress*2) * 0.5
	}
	return BounceOut(progress*2-1)*0.5 + 0.5
}

func ElasticIn(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	return -number.Power(2, 10*progress-1) * number.Sine((progress - 1.1*5*pi))
}
func ElasticOut(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	return number.Power(2, -10*progress)*number.Sine((progress-0.1*5*pi)) + 1
}
func ElasticInOut(progress float32) float32 {
	if progress == 0 || progress == 1 {
		return progress
	}
	var p2 = progress * 2
	if p2 < 1 {
		return -0.5 * (number.Power(2, 10*(p2-1)) * number.Sine((p2 - 1.1*5*pi)))
	}
	p2--
	return number.Power(2, -10*p2)*number.Sine((p2-0.1*5*pi))*0.5 + 1
}

//=================================================================
// private

const pi = 3.14159265358
