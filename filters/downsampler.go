package filters

import (
	"fmt"
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var downsamplerId uint64

type downsamplerState[T buffers.BlockType] struct {
	input     *buffers.Stream[T]
	output    *buffers.Stream[T]
	timestamp *buffers.Timestamp
	logger    *log.Entry
}

func NewDownsampler[T buffers.BlockType](
	input *buffers.Stream[T],
	targetRate int,
) (output *buffers.Stream[T], err error) {
	if targetRate > input.Size {
		err = fmt.Errorf("targetRate is greater than input rate: %d > %d", targetRate, input.Size)
		return
	}
	var id = atomic.AddUint64(&downsamplerId, 1)
	output = buffers.NewStream[T](targetRate, input.Count)
	var filter = &downsamplerState[T]{
		input:  input,
		output: output,
		logger: log.WithFields(log.Fields{
			"filter": "Downsampler",
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func (filter *downsamplerState[T]) blockHandler(block *buffers.Block[T]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	var out = filter.output.Next()
	out.WriteRaw(block.Data, filter.timestamp)
	filter.output.Send(out)
}

func (filter *downsamplerState[T]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *downsamplerState[T]) loop() {
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
