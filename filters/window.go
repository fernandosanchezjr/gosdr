package filters

func computeWindow(input []complex64, window []float64) {
	for pos, value := range window {
		var floatValue = float32(value)
		var complexValue = input[pos]
		input[pos] = complex(real(complexValue)*floatValue, imag(complexValue)*floatValue)
	}
}
