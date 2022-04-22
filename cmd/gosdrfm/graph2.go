package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
	"github.com/fernandosanchezjr/gosdr/units"
)

func createGraph(conn devices.Connection, bufferCount int) (input *buffers.Stream[byte], err error) {
	var sampleBufferSize = int(conn.SampleBufferSize())
	input = buffers.NewStream[byte](sampleBufferSize, bufferCount)
	var complexOutput = filters.NewBytesToComplexConverter[complex64](input)
	var resampledOutput = filters.NewResampler(complexOutput, int(units.Sps(200000).NearestSize(512)))
	var fftOutput, fftErr = filters.NewFFT[complex64](resampledOutput, 1024)
	if fftErr != nil {
		err = fftErr
		return
	}
	filters.NewHistogram[complex64](fftOutput)
	return
}
