package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph(conn devices.Connection, bufferCount int) (input *buffers.Stream[byte], err error) {
	var sampleBufferSize = int(conn.SampleBufferSize())
	input = buffers.NewStream[byte](sampleBufferSize, bufferCount)
	var complexOutput = filters.NewBytesToComplexConverter[complex64](input)
	//var resampledOutput = filters.NewResampler(complexOutput, 2400000)
	//var rationalOutput = filters.NewRationalResampler(complexOutput, 1, 2)
	var fftOutput, fftErr = filters.NewFFT[complex64](complexOutput, 1024)
	if fftErr != nil {
		err = fftErr
		return
	}
	filters.NewHistogram[complex64](fftOutput)
	return
}
