package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQSampler(sampleRate int, bufferCount int) (int, chan []byte, chan *buffers.IQ) {
	var iqSampleRate = sampleRate / 2
	var input = make(chan []byte, bufferCount)
	var output = make(chan *buffers.IQ, bufferCount)
	var ring = buffers.NewIQRing(iqSampleRate, bufferCount)
	go iqSamplerLoop(sampleRate, ring, input, output)
	return iqSampleRate, input, output
}

func iqSamplerLoop(rawSampleRate int, ring *buffers.IQRing, input chan []byte, output chan *buffers.IQ) {
	log.WithField("filter", "IQSampler").Debug("Starting")
	var sequence uint64
	var out = ring.Next()
	var readOffset int
	for {
		select {
		case raw, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQSampler").Debug("Exiting")
				runtime.GC()
				return
			}
			readOffset = 0
			for {
				var read, readErr = out.Read(raw[readOffset:])
				readOffset += read
				if readErr != nil {
					log.WithField("filter", "IQSampler").WithError(readErr).Error("buf.Read")
					break
				}
				if out.Full() {
					out.Sequence = sequence
					output <- out
					sequence += 1
					out = ring.Next()
				}
				if readOffset == len(raw) {
					break
				}
			}
		}
	}
}
