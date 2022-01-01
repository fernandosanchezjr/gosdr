package utils

import "math"

func ComplexMagnitude(value complex128) (magnitude float64) {
	magnitude = math.Sqrt(math.Pow(real(value), 2.0) + math.Pow(imag(value), 2.0))
	if math.IsNaN(magnitude) {
		return 0.0
	}
	return
}
