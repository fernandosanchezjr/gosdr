package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	log "github.com/sirupsen/logrus"
)

func sampleDevice(conn devices.Connection, bufferCount int, output *buffers.Stream[byte]) {
	var blockRing = buffers.NewBlockRing[byte](int(conn.SampleBufferSize()), bufferCount)
	var timestamp = buffers.NewTimestamp()
	var samplerFunc = func(samples []byte) {
		var block = blockRing.Next()
		block.WriteRaw(samples, timestamp)
		output.Send(block)
		timestamp.Increment()
	}
	if err := conn.RunSampler(samplerFunc); err != nil {
		log.WithFields(conn.Fields()).WithError(err).Trace("conn.RunSampler")
	}
}
