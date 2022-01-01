package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
)

func NewIQSampler(size int, count int) (int, chan []byte, chan *buffers.IQ) {
	var iqSampleRate = size / 2
	var input = make(chan []byte, count)
	var output = make(chan *buffers.IQ, count)
	var ring = buffers.NewIQRing(iqSampleRate, count)
	go iqSamplerLoop(ring, input, output)
	return iqSampleRate, input, output
}

func iqSamplerLoop(ring *buffers.IQRing, input chan []byte, output chan *buffers.IQ) {
	for {
		select {
		case raw, ok := <-input:
			if !ok {
				close(output)
				return
			}
			var out = ring.Next()
			if _, err := out.Read(raw); err != nil {
				log.WithField("filter", "IQSampler").WithError(err).Error("buf.Read")
			} else {
				output <- out
			}
		}
	}
}
