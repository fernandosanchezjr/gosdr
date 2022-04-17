package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	log "github.com/sirupsen/logrus"
)

//func sampleDevice(conn devices.Connection, output chan []byte) {
//	var ring = buffers.NewByteRing(int(conn.SampleBufferSize()), 16)
//	var samplerFunc = func(samples []byte) {
//		var buffer = ring.Next()
//		copy(buffer, samples)
//		output <- buffer
//	}
//	if err := conn.RunSampler(samplerFunc); err != nil {
//		log.WithFields(conn.Fields()).WithError(err).Trace("conn.RunSampler")
//	}
//}

func sampleDevice2(conn devices.Connection, bufferCount int, stream *buffers.Stream) {
	var blockRing = buffers.NewBlockRing(int(conn.SampleBufferSize()), bufferCount)
	var timestamp = buffers.NewTimestamp()
	var samplerFunc = func(samples []byte) {
		var block = blockRing.Next()
		block.WriteBytes(samples, timestamp)
		log.WithField("ts", timestamp).Info("Block read")
		stream.Send(block)
		timestamp.Increment()
	}
	if err := conn.RunSampler(samplerFunc); err != nil {
		log.WithFields(conn.Fields()).WithError(err).Trace("conn.RunSampler")
	}
}
