package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph(conn devices.Connection, bufferCount int) (input *buffers.Stream[byte], err error) {
	input = buffers.NewStream[byte](int(conn.SampleBufferSize()), bufferCount)
	var demodComplexOutput = filters.NewBytesToComplexConverter[complex128](input)
	var demodFFTOutput, demodFFTErr = filters.NewFFT(demodComplexOutput, demodComplexOutput.Size)
	if demodFFTErr != nil {
		err = demodFFTErr
		return
	}
	filters.NewHistogram[complex128](demodFFTOutput, 1024)
	return
}
