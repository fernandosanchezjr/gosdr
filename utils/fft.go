package utils

import (
	"github.com/argusdusty/gofft"
)

func FFTComplex128(input, output []complex128, midPoint int) {
	if err := gofft.FFT(input); err != nil {
		panic(err)
	}
	copy(output[:midPoint], input[midPoint:])
	copy(output[midPoint:], input[:midPoint])
}
