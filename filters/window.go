package filters

import "github.com/racerxdl/segdsp/dsp"

func computeWindowComplex64(input []complex64, window []float32) {
	for pos, value := range window {
		var complexValue = input[pos]
		input[pos] = complex(real(complexValue)*value, imag(complexValue)*value)
	}
}

func createWindowComplex64(fftSize int) []float32 {
	var rawWindow = dsp.BlackmanHarris(fftSize, 92)
	var window = make([]float32, len(rawWindow))
	for i := range window {
		window[i] = float32(rawWindow[i])
	}
	return window
}
