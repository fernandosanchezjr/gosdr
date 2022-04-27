package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var bytesToComplexConverterId uint64

type bytesToComplexTypes interface {
	complex64 | complex128
}

type bytesToComplexConverterState[T bytesToComplexTypes] struct {
	input      *buffers.Stream[byte]
	output     *buffers.Stream[T]
	timestamp  *buffers.Timestamp
	window32   []float32
	window64   []float64
	logger     *log.Entry
	resultType T
}

func NewBytesToComplexConverter[T bytesToComplexTypes](
	input *buffers.Stream[byte],
) (output *buffers.Stream[T]) {
	var newSize = input.Size / 2
	var id = atomic.AddUint64(&bytesToComplexConverterId, 1)
	output = buffers.NewStream[T](newSize, input.Count)
	var filter = &bytesToComplexConverterState[T]{
		input:    input,
		output:   output,
		window32: utils.CreateWindow32(newSize),
		window64: utils.CreateWindow64(newSize),
		logger: log.WithFields(log.Fields{
			"filter": "Bytes To Complex Converter",
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func float32Converter[T bytesToComplexTypes](data []byte, out *buffers.Block[T], window []float32) {
	var outPos int
	var windowValue float32
	for i := 0; i < len(data); i += 2 {
		windowValue = window[outPos]
		out.Data[outPos] = T(complex(
			utils.ConvertByte[float32](data[i])*windowValue,
			utils.ConvertByte[float32](data[i+1])*windowValue,
		))
		outPos++
	}
}

func float64Converter[T bytesToComplexTypes](data []byte, out *buffers.Block[T], window []float64) {
	var outPos int
	var windowValue float64
	for i := 0; i < len(data); i += 2 {
		windowValue = window[outPos]
		out.Data[outPos] = T(complex(
			utils.ConvertByte[float64](data[i])*windowValue,
			utils.ConvertByte[float64](data[i+1])*windowValue,
		))
		outPos++
	}
}

func (filter *bytesToComplexConverterState[T]) blockHandler(block *buffers.Block[byte]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	var outBlock = filter.output.Next()
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	filter.timestamp.Copy(outBlock.Timestamp)
	switch any(filter.resultType).(type) {
	case complex64:
		float32Converter[T](block.Data, outBlock, filter.window32)
	case complex128:
		float64Converter[T](block.Data, outBlock, filter.window64)
	}
	filter.output.Send(outBlock)
}

func (filter *bytesToComplexConverterState[T]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *bytesToComplexConverterState[T]) loop() {
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
