package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
	"github.com/fernandosanchezjr/gosdr/units"
)

const (
	fftSize = 1024
)

func createGraph(conn devices.Connection, bufferCount int, quit chan struct{}) chan []byte {
	var iqSampleRate, iqSamplerInput, iqSamplerOutput = filters.NewIQSampler(
		int(conn.GetSampleRate()),
		bufferCount,
		quit,
	)
	var resampleRate, resamplerOutput = filters.NewIQRationalResampler(
		iqSampleRate,
		bufferCount,
		2,
		16,
		iqSamplerOutput,
		quit,
	)
	var filterOutput = filters.NewIQLowpassFilter(
		resampleRate,
		bufferCount,
		1.0,
		50_000,
		1000,
		resamplerOutput,
		quit,
	)
	var fftOutput = filters.NewIQFFT(fftSize, resampleRate, bufferCount, filterOutput, quit)
	filters.NewIQHistogram(fftSize, bufferCount, conn, units.Hertz(resampleRate), fftOutput, quit)
	return iqSamplerInput
}
