package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var nullSinkId uint64

func NewNullSink[T buffers.BlockType](input *buffers.Stream[T]) {
	go nullSinkLoop[T](input)
}

func nullSinkLoop[T buffers.BlockType](input *buffers.Stream[T]) {
	var closed bool
	var id = atomic.AddUint64(&nullSinkId, 1)
	var logger = log.WithFields(log.Fields{
		"filter": "NullSink",
		"id":     id,
	})
	var handler = func(block *buffers.Block[T]) {
		logger.WithField("block", block).Trace("Stream")
	}
	logger.Debug("Starting")
	for {
		if closed = input.Receive(handler); closed {
			input.Done()
			logger.Debug("Exiting")
			return
		}
	}
}
