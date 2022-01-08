package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQFilter(
	sampleRate int,
	decimationRate int,
	bufferCount int,
	input chan *buffers.IQ,
) (int, chan *buffers.IQ) {
	var output = make(chan *buffers.IQ, bufferCount)
	var outputSize = sampleRate / decimationRate
	var outputRing = buffers.NewIQRing(outputSize, bufferCount)
	go iqFilterLoop(outputRing, input, output)
	return outputSize, output
}

func iqFilterLoop(
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
) {
	log.WithField("filter", "Decimator").Debug("Starting")
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "Decimator").Debug("Exiting")
				runtime.GC()
				return
			}
			var out = outputRing.Next()
			out.Copy(in)
			output <- out
		}
	}
}
