package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	log "github.com/sirupsen/logrus"
	"math"
	"sync/atomic"
)

var resamplerId uint64

type resamplerState[T buffers.BlockType] struct {
	input               *buffers.Stream[T]
	output              *buffers.Stream[T]
	timestamp           *buffers.Timestamp
	currentBlock        *buffers.Block[T]
	bufferBlocks        *buffers.BlockRing[T]
	rawBufferBlocks     []*buffers.Block[T]
	lastBufferBlock     int
	lastBufferBlockRead int
	logger              *log.Entry
}

func NewResampler[T buffers.BlockType](
	input *buffers.Stream[T],
	targetRate int,
) (output *buffers.Stream[T]) {
	var id = atomic.AddUint64(&resamplerId, 1)
	output = buffers.NewStream[T](targetRate, input.Count)
	var bufferBlockCount = int(math.Ceil(float64(targetRate) / float64(input.Size)))
	var filter = &resamplerState[T]{
		input:           input,
		output:          output,
		currentBlock:    output.Next(),
		bufferBlocks:    buffers.NewBlockRing[T](input.Size, bufferBlockCount),
		rawBufferBlocks: make([]*buffers.Block[T], bufferBlockCount),
		lastBufferBlock: bufferBlockCount,
		logger: log.WithFields(log.Fields{
			"filter": "Resampler",
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func (filter *resamplerState[T]) blockHandler(block *buffers.Block[T]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	var tmpBlock = filter.bufferBlocks.Next()
	tmpBlock.WriteRaw(block.Data, block.Timestamp)
	if filter.lastBufferBlock != filter.lastBufferBlockRead {
		filter.lastBufferBlockRead += 1
		return
	}
	filter.bufferBlocks.ReverseCopy(filter.rawBufferBlocks)
	for _, tmpBlock = range filter.rawBufferBlocks {
		filter.timestamp = tmpBlock.CopyTimestamp(filter.timestamp)
		for !tmpBlock.End() {
			filter.currentBlock.Write(tmpBlock, filter.timestamp)
			if filter.currentBlock.End() {
				filter.output.Send(filter.currentBlock)
				filter.currentBlock = filter.output.Next()
				filter.timestamp.Increment()
				if (tmpBlock.Size > filter.currentBlock.Size && tmpBlock.Remainder() < filter.currentBlock.Size) ||
					(tmpBlock.Size < filter.currentBlock.Size) {
					return
				}
			}
		}
	}
	//for !block.End() {
	//	filter.currentBlock.Write(block, filter.timestamp)
	//	if filter.currentBlock.End() {
	//		filter.output.Send(filter.currentBlock)
	//		filter.currentBlock = filter.output.Next()
	//		filter.timestamp.Increment()
	//	}
	//}
}

func (filter *resamplerState[T]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *resamplerState[T]) loop() {
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
