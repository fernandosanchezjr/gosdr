package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"runtime"
)

type iqMultiplexerDestination struct {
	outputRing *buffers.IQRing
	output     chan *buffers.IQ
}

func NewIQMultiplexer(
	sampleRate int,
	bufferCount int,
	input chan *buffers.IQ,
	output ...chan *buffers.IQ,
) {
	var destinations = make([]*iqMultiplexerDestination, len(output))
	for index, out := range output {
		destinations[index] = &iqMultiplexerDestination{
			outputRing: buffers.NewIQRing(sampleRate, bufferCount),
			output:     out,
		}
	}
	go iqMultiplexerLoop(input, destinations)
}

func iqMultiplexerLoop(
	input chan *buffers.IQ,
	outputs []*iqMultiplexerDestination,
) {
	log.WithField("filter", "IQMultiplexer").Debug("Starting")
	for {
		select {
		case in, ok := <-input:
			if !ok {
				for _, destination := range outputs {
					close(destination.output)
				}
				outputs = nil
				log.WithField("filter", "IQMultiplexer").Debug("Exiting")
				runtime.GC()
				return
			}
			for _, destination := range outputs {
				var out = destination.outputRing.Next()
				out.Copy(in)
				destination.output <- out
			}
		}
	}
}
