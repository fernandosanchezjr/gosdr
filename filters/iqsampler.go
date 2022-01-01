package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
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
	var sequence uint64
	for {
		select {
		case raw, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQSampler").Trace("Exiting")
				return
			}
			for {
				var out = ring.Next()
				var read, err = out.Read(raw)
				if err != nil {
					log.WithField("filter", "IQSampler").WithError(err).Error("buf.Read")
					break
				}
				sequence += 1
				out.Sequence = sequence
				output <- out
				if read == rawSampleRate {
					break
				}
			}
		}
	}
}
