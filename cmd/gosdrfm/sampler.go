package main

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	log "github.com/sirupsen/logrus"
)

func sampleDevice(conn devices.Connection) {
	var ring = buffers.NewByteRing(int(conn.SampleBufferSize()), 32)
	var samplerFunc = func(samples []byte) {
		var buffer = ring.Next()
		copy(buffer, samples)
	}
	if err := conn.RunSampler(samplerFunc); err != nil {
		log.WithFields(conn.Fields()).WithError(err).Error("conn.RunSampler")
	}
}
