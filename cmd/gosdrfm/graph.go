package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph(conn devices.Connection) chan []byte {
	var iqSampleRate, iqSamplerInput, iqSamplerOutput = filters.NewIQSampler(int(conn.SampleBufferSize()), 16)
	var _, decimatorOutput = filters.NewDecimator(iqSampleRate, 128, 16, iqSamplerOutput)
	filters.NewNullIQSink(decimatorOutput)
	return iqSamplerInput
}
