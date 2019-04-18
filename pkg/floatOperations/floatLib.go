package floatlib

import "math"

func FloatIsWhole(f float64) bool {
	const epsilon = 1e-9
	if _, frac := math.Modf(f); frac < epsilon || frac > 1.0-epsilon {
		return true
	}
	return false
}
