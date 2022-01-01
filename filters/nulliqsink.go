package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
)

func NewNullIQSink(input chan *buffers.IQ) {
	go nullIQSinkLoop(input)
}

func nullIQSinkLoop(input chan *buffers.IQ) {
	for {
		select {
		case iq, ok := <-input:
			if !ok {
				return
			}
			log.WithField("size", len(iq.Data())).Trace("Null IQ Sink")
		}
	}
}
