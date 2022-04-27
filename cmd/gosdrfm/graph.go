package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

const histogramFFTSize = 1024

func createGraph(conn devices.Connection, bufferCount int) (input *buffers.Stream[byte], err error) {
	input = buffers.NewStream[byte](int(conn.SampleBufferSize()), bufferCount)
	var demodulatorInput = input.Clone()
	var histogramInput = input.Clone()
	filters.NewSplitter(input, []*buffers.Stream[byte]{demodulatorInput, histogramInput})
	// demodulator section
	var demodComplexOutput = filters.NewBytesToComplexConverter[complex128](demodulatorInput)
	var demodFFTOutput, demodFFTErr = filters.NewFFT(demodComplexOutput, demodComplexOutput.Size)
	if demodFFTErr != nil {
		err = demodFFTErr
		return
	}
	filters.NewNullSink[complex128](demodFFTOutput)
	// histogram section
	var histogramResampledOutput, histogramResampleErr = filters.NewDownsampler[byte](
		histogramInput,
		histogramFFTSize*2,
	)
	if histogramResampleErr != nil {
		err = histogramResampleErr
		return
	}
	var histogramComplexOutput = filters.NewBytesToComplexConverter[complex128](histogramResampledOutput)
	var histogramFFTOutput, histogramFFTErr = filters.NewFFT(histogramComplexOutput, histogramFFTSize)
	if demodFFTErr != nil {
		err = histogramFFTErr
		return
	}
	filters.NewHistogram[complex128](histogramFFTOutput)
	return
}
