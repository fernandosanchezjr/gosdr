package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
	"github.com/fernandosanchezjr/gosdr/units"
)

func createGraph(conn devices.Connection, bufferCount int) chan []byte {
	var iqSampleRate, iqSamplerInput, iqSamplerOutput = filters.NewIQSampler(int(conn.GetSampleRate()), bufferCount)
	var resampleRate, resamplerOutput = filters.NewIQRationalResampler(
		iqSampleRate,
		bufferCount,
		2,
		16,
		iqSamplerOutput,
	)
	var filterOutput = filters.NewIQLowpassFilter(
		resampleRate,
		bufferCount,
		1.0,
		50_000,
		1000,
		resamplerOutput,
	)
	var fftOutput = filters.NewIQFFT(resampleRate, bufferCount, filterOutput)
	filters.NewIQHistogram(resampleRate, bufferCount, conn, units.Hertz(resampleRate), fftOutput)
	return iqSamplerInput
}
