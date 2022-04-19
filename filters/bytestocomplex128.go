package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var bytesToComplexConverterId uint64

type bytesToComplexConverterState[O complex64 | complex128] struct {
	id         uint64
	input      *buffers.Stream[byte]
	output     *buffers.Stream[O]
	ring       *buffers.BlockRing[O]
	timestamp  *buffers.Timestamp
	resultType O
}

func NewBytesToComplexConverter[O complex64 | complex128](
	input *buffers.Stream[byte],
	bufferSize,
	bufferCount int,
) *buffers.Stream[O] {
	filter := &bytesToComplexConverterState[O]{
		id:     atomic.AddUint64(&bytesToComplexConverterId, 1),
		input:  input,
		output: buffers.NewStream[O](bufferCount),
		ring:   buffers.NewBlockRing[O](bufferSize/2, bufferCount),
	}
	go filter.loop()
	return filter.output
}

func float32Converter[O complex64 | complex128](data []byte, out *buffers.Block[O]) {
	var outPos int
	for i := 0; i < len(data); i += 2 {
		out.Data[outPos] = O(complex(utils.ConvertByte(data[i]), utils.ConvertByte(data[i+1])))
		outPos++
	}
}

func float64Converter[O complex64 | complex128](data []byte, out *buffers.Block[O]) {
	var outPos int
	for i := 0; i < len(data); i += 2 {
		out.Data[outPos] = O(complex(float64(utils.ConvertByte(data[i])), float64(utils.ConvertByte(data[i+1]))))
		outPos++
	}
}

func (filter *bytesToComplexConverterState[O]) blockHandler(block *buffers.Block[byte]) {
	log.WithFields(log.Fields{
		"filter": "BytesToComplexConverter",
		"id":     filter.id,
		"block":  block,
	}).Trace("Stream")
	var out = filter.ring.Next()
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	filter.timestamp.Copy(out.Timestamp)
	switch any(filter.resultType).(type) {
	case complex64:
		float32Converter[O](block.Data, out)
	case complex128:
		float64Converter[O](block.Data, out)
	}
	filter.output.Send(out)
}

func (filter *bytesToComplexConverterState[O]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
	filter.ring = nil
}

func (filter *bytesToComplexConverterState[O]) loop() {
	var closed bool
	var logger = log.WithFields(log.Fields{
		"filter": "BytesToComplexConverter",
		"id":     filter.id,
	})
	logger.Debug("Starting")
	for {
		if closed = filter.input.Receive(filter.blockHandler); closed {
			logger.Debug("Exiting")
			filter.close()
			return
		}
	}
}
