package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
)

func blockHandler(*buffers.Block) {

}

func createGraph2(conn devices.Connection, bufferCount int) *buffers.Stream {
	var inputStream = buffers.NewStream(bufferCount)

	go func() {
		var closed bool
		for {
			closed = inputStream.Receive(blockHandler)
			if closed {
				inputStream.Done()
				return
			}
		}
	}()

	return inputStream
}
