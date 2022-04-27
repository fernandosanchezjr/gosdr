package filters

import (
	"fmt"
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var fftId uint64

type fftState struct {
	input     *buffers.Stream[complex128]
	output    *buffers.Stream[complex128]
	timestamp *buffers.Timestamp
	midPoint  int
	size      int
	logger    *log.Entry
}

func NewFFT(
	input *buffers.Stream[complex128],
	fftSize int,
) (output *buffers.Stream[complex128], err error) {
	if fftSize > input.Size {
		err = fmt.Errorf("fft size is greater than sample size: %d > %d", fftSize, input.Size)
		return
	}
	var id = atomic.AddUint64(&fftId, 1)
	output = buffers.NewStream[complex128](fftSize, input.Count)
	var filter = &fftState{
		input:    input,
		output:   output,
		size:     fftSize,
		midPoint: (fftSize / 2) + (fftSize % 2),
		logger: log.WithFields(log.Fields{
			"filter": fmt.Sprintf("FFT(complex12) size %d", fftSize),
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func (filter *fftState) blockHandler(block *buffers.Block[complex128]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	var outputBlock = filter.output.Next()
	utils.FFTComplex128(block.Data[:filter.size], outputBlock.Data, filter.midPoint)
	filter.timestamp.Copy(outputBlock.Timestamp)
	filter.output.Send(outputBlock)
}

func (filter *fftState) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *fftState) loop() {
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
