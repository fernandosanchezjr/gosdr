package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var rationalResamplerId uint64

type rationalResamplerTypes interface {
	complex64
}

type rationalResamplerState[T rationalResamplerTypes] struct {
	input      *buffers.Stream[T]
	output     *buffers.Stream[T]
	timestamp  *buffers.Timestamp
	resampler  *dsp.RationalResampler
	logger     *log.Entry
	resultType T
}

func NewRationalResampler[T rationalResamplerTypes](
	input *buffers.Stream[T],
	interpolation int,
	decimation int,
) (output *buffers.Stream[T]) {
	var resampler = dsp.MakeRationalResampler(interpolation, decimation)
	var resampleRate = resampler.PredictOutputSize(input.Size)
	var id = atomic.AddUint64(&rationalResamplerId, 1)
	output = buffers.NewStream[T](resampleRate, input.Count)
	var filter = &rationalResamplerState[T]{
		input:     input,
		output:    output,
		resampler: resampler,
		logger: log.WithFields(log.Fields{
			"filter": "Rational Resampler",
			"id":     id,
		}),
	}
	go filter.loop()
	return
}

func (filter *rationalResamplerState[T]) blockHandler(block *buffers.Block[T]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	filter.timestamp = block.CopyTimestamp(filter.timestamp)
	var out = filter.output.Next()
	switch outBuf := any(out.Data).(type) {
	case []complex64:
		switch inBuf := any(block.Data).(type) {
		case []complex64:
			filter.resampler.WorkBuffer(inBuf, outBuf)
		}
	}
	filter.timestamp.Copy(out.Timestamp)
	filter.output.Send(out)
}

func (filter *rationalResamplerState[T]) close() {
	filter.output.Close()
	filter.input.Done()
	filter.input = nil
	filter.output = nil
}

func (filter *rationalResamplerState[T]) loop() {
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
