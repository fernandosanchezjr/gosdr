package filters

import (
	"fmt"
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp/fft"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var fftId uint64

type fftTypes interface {
	complex64
}

type fftState[T fftTypes] struct {
	input       *buffers.Stream[T]
	output      *buffers.Stream[T]
	timestamp   *buffers.Timestamp
	fftBlock    *buffers.Block[T]
	window      []float32
	skipSamples int
	midPoint    int
	logger      *log.Entry
	resultType  T
}

func NewFFT[T fftTypes](
	input *buffers.Stream[T],
	fftSize int,
) (output *buffers.Stream[T], err error) {
	if fftSize > input.Size {
		err = fmt.Errorf("fft size is greater than sample size: %d > %d", fftSize, input.Size)
		return
	}
	var id = atomic.AddUint64(&fftId, 1)
	var fftBlock = buffers.NewBlock[T](fftSize)
	output = buffers.NewStream[T](fftSize, input.Count)
	var filter = &fftState[T]{
		input:       input,
		output:      output,
		fftBlock:    fftBlock,
		window:      createWindowComplex64(fftSize),
		skipSamples: input.Size / fftSize,
		midPoint:    (fftSize / 2) + (fftSize % 2),
		logger: log.WithFields(log.Fields{
			"filter": fmt.Sprintf("FFT(%T)", fftBlock.Data[0]),
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func (filter *fftState[T]) fftComplex64(input, output []complex64) {
	computeWindowComplex64(input, filter.window)
	var fftOut = fft.FFT(input)
	copy(output[:filter.midPoint], fftOut[filter.midPoint:])
	copy(output[filter.midPoint:], fftOut[:filter.midPoint])
}

func (filter *fftState[T]) blockHandler(block *buffers.Block[T]) {
	filter.logger.WithField("block", block).Trace("Stream")
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	copy(filter.fftBlock.Data, block.Data)
	var outputBlock = filter.output.Next()
	switch inBuf := any(filter.fftBlock.Data).(type) {
	case []complex64:
		switch outBuf := any(outputBlock.Data).(type) {
		case []complex64:
			filter.fftComplex64(inBuf, outBuf)
		}
	}
	filter.timestamp.Copy(outputBlock.Timestamp)
	filter.output.Send(outputBlock)
}

func (filter *fftState[T]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *fftState[T]) loop() {
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
