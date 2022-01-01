package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph(conn devices.Connection, bufferCount int) chan []byte {
	var iqSampleRate, iqSamplerInput, iqSamplerOutput = filters.NewIQSampler(int(conn.SampleBufferSize()), bufferCount)
	var multiplexer = make([]chan *buffers.IQ, 4)
	for i := 0; i < 4; i++ {
		var output = make(chan *buffers.IQ, bufferCount)
		multiplexer[i] = output
		filters.NewNullSink(output, nil)
	}
	filters.NewIQMultiplexer(iqSampleRate, bufferCount, iqSamplerOutput, multiplexer...)

	return iqSamplerInput
}
