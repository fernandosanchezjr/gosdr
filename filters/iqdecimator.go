package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQDecimator(
	sampleRate int,
	bufferCount int,
	decimation int,
	input chan *buffers.IQ,
) (int, chan *buffers.IQ) {
	var decimator = dsp.MakeDecimator(decimation)
	var resampleRate = decimator.PredictOutputSize(sampleRate)
	var outputRing = buffers.NewIQRing(resampleRate, bufferCount-1)
	var output = make(chan *buffers.IQ, bufferCount)
	go iqDecimatorLoop(sampleRate, decimator, outputRing, input, output)
	log.WithFields(log.Fields{
		"sampleRate":    sampleRate,
		"decimatedRate": resampleRate,
	}).Debug("IQDecimator")
	return resampleRate, output
}

func iqDecimatorLoop(
	sampleRate int,
	resampler *dsp.Decimator,
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
) {
	log.WithField("filter", "IQDecimator").Debug("Starting")
	var inBuffer = buffers.NewIQ(sampleRate)
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQDecimator").Debug("Exiting")
				runtime.GC()
				return
			}
			inBuffer.Copy(in)
			var out = outputRing.Next()
			resampler.WorkBuffer(inBuffer.Data(), out.Data())
			out.Sequence = inBuffer.Sequence
			output <- out
		}
	}
}
