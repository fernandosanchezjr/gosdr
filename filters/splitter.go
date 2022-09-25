package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var splitterId uint64

type splitterState[T buffers.BlockType] struct {
	input     *buffers.Stream[T]
	outputs   []*buffers.Stream[T]
	timestamp *buffers.Timestamp
	logger    *log.Entry
}

func NewSplitter[T buffers.BlockType](input *buffers.Stream[T], outputs ...*buffers.Stream[T]) {
	var id = atomic.AddUint64(&splitterId, 1)
	var filter = &splitterState[T]{
		input:   input,
		outputs: outputs,
		logger: log.WithFields(log.Fields{
			"filter": "Splitter",
			"id":     id,
		}),
	}
	go filter.loop()
}

func (filter *splitterState[T]) blockHandler(block *buffers.Block[T]) {
	var out *buffers.Block[T]
	filter.logger.WithField("block", block).Trace("Stream in")
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	for _, output := range filter.outputs {
		out = output.Next()
		out.WriteRaw(block.Data, filter.timestamp)
		filter.timestamp.Increment()
		output.Send(out)
	}
}

func (filter *splitterState[T]) close() {
	for _, output := range filter.outputs {
		output.Close()
	}
	filter.input.Done()
	filter.input = nil
	filter.outputs = nil
}

func (filter *splitterState[T]) loop() {
	var closed bool
	filter.logger.Debug("Starting")
	for {
		if closed = filter.input.Receive(filter.blockHandler); closed {
			filter.logger.Debug("Exiting")
			filter.close()
			return
		}
	}
}
