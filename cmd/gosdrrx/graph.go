package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/filters"
)

func createGraph(conn devices.Connection, bufferCount int) (input *buffers.Stream[byte], err error) {
	input = buffers.NewStream[byte](int(conn.SampleBufferSize()), bufferCount)
	filters.NewSplitter[byte](input)
	return
}
