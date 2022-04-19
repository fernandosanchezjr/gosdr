package utils

import "math"

type Complexes interface {
	complex64 | complex128
}

func ComplexMagnitude(value complex128) (magnitude float64) {
	magnitude = math.Sqrt(math.Pow(real(value), 2.0) + math.Pow(imag(value), 2.0))
	if math.IsNaN(magnitude) {
		return 0.0
	}
	return
}
