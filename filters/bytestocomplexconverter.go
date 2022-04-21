package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var bytesToComplexConverterId uint64

type bytesToComplexConverterState[O complex64 | complex128] struct {
	input      *buffers.Stream[byte]
	output     *buffers.Stream[O]
	timestamp  *buffers.Timestamp
	logger     *log.Entry
	resultType O
}

func NewBytesToComplexConverter[O complex64 | complex128](
	input *buffers.Stream[byte],
) (output *buffers.Stream[O]) {
	var newSize = input.Size / 2
	var id = atomic.AddUint64(&bytesToComplexConverterId, 1)
	output = buffers.NewStream[O](newSize, input.Count)
	var filter = &bytesToComplexConverterState[O]{
		input:  input,
		output: output,
		logger: log.WithFields(log.Fields{
			"filter": "BytesToComplexConverter",
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func float32Converter[O complex64 | complex128](data []byte, out *buffers.Block[O]) {
	var outPos int
	for i := 0; i < len(data); i += 2 {
		out.Data[outPos] = O(complex(utils.ConvertByte[float32](data[i]), utils.ConvertByte[float32](data[i+1])))
		outPos++
	}
}

func float64Converter[O complex64 | complex128](data []byte, out *buffers.Block[O]) {
	var outPos int
	for i := 0; i < len(data); i += 2 {
		out.Data[outPos] = O(complex(utils.ConvertByte[float64](data[i]), utils.ConvertByte[float64](data[i+1])))
		outPos++
	}
}

func (filter *bytesToComplexConverterState[O]) blockHandler(block *buffers.Block[byte]) {
	filter.logger.WithField("block", block).Trace("Stream")
	var outBlock = filter.output.Next()
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	filter.timestamp.Copy(outBlock.Timestamp)
	switch any(filter.resultType).(type) {
	case complex64:
		float32Converter[O](block.Data, outBlock)
	case complex128:
		float64Converter[O](block.Data, outBlock)
	}
	filter.output.Send(outBlock)
}

func (filter *bytesToComplexConverterState[O]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *bytesToComplexConverterState[O]) loop() {
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
