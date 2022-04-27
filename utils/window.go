package utils

import "github.com/racerxdl/segdsp/dsp"

func CreateWindow32(fftSize int) []float32 {
	var window64 = dsp.BlackmanHarris(fftSize, 92)
	var window32 = make([]float32, len(window64))
	for pos, value := range window64 {
		window32[pos] = float32(value)
	}
	return window32
}

func CreateWindow64(fftSize int) []float64 {
	return dsp.BlackmanHarris(fftSize, 92)
}

func ComputeWindowComplex128(input []complex128, window []float64) {
	var complexValue complex128
	for pos, value := range window {
		complexValue = input[pos]
		input[pos] = complex(real(complexValue)*value, imag(complexValue)*value)
	}
}
