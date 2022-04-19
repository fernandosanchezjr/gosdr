package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph2(conn devices.Connection, bufferCount int) *buffers.Stream[byte] {
	var sampleBufferSize = int(conn.SampleBufferSize())
	var inputStream = buffers.NewStream[byte](bufferCount)

	var complexOutput = filters.NewBytesToComplexConverter[complex128](inputStream, sampleBufferSize, bufferCount)

	filters.NewNullSink[complex128](complexOutput)

	return inputStream
}
