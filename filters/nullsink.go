package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var nullSinkId uint64

func NewNullSink(iqInput chan *buffers.IQ, floatInput chan []float64) {
	if iqInput != nil {
		go nullSinkLoop(iqInput, floatInput)
	}
}

func nullSinkLoop(iqInput chan *buffers.IQ, floatInput chan []float64) {
	var id = atomic.AddUint64(&nullSinkId, 1)
	for {
		select {
		case in, ok := <-iqInput:
			if !ok {
				log.WithFields(log.Fields{
					"filter": "NullSink",
					"id":     id,
				}).Trace("Exiting")
				return
			}
			log.WithFields(log.Fields{
				"filter":   "NullSink",
				"id":       id,
				"sequence": in.Sequence,
				"type":     "IQ",
			}).Trace("Sample received")
		case _, ok := <-floatInput:
			if !ok {
				log.WithFields(log.Fields{
					"filter": "NullSink",
					"id":     id,
				}).Trace("Exiting")
				return
			}
			log.WithFields(log.Fields{
				"filter": "NullSink",
				"id":     id,
				"type":   "Float",
			}).Trace("Sample received")
		}
	}
}
