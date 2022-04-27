package utils

import (
	"github.com/racerxdl/segdsp/tools"
	"math"
)

func GetPower64(sample complex64) float64 {
	if imag(sample) == 0.0 {
		return 0.0
	}
	var value = 10.0 * math.Log10(float64(tools.ComplexAbsSquared(sample)))
	if math.IsInf(value, 1) {
		return 10.0
	} else if math.IsInf(value, -1) {
		return -10.0
	} else if math.IsNaN(value) {
		return 0
	}
	return value
}

func ComplexAbsSquared128(x complex128) float64 {
	return real(x)*real(x) + imag(x)*imag(x)
}

func GetPower128(sample complex128) float64 {
	if imag(sample) == 0.0 {
		return 0.0
	}
	var value = 10.0 * math.Log10(ComplexAbsSquared128(sample))
	if math.IsInf(value, 1) {
		return 10.0
	} else if math.IsInf(value, -1) {
		return -10.0
	} else if math.IsNaN(value) {
		return 0
	}
	return value
}
