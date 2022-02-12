package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQRationalResampler(
	sampleRate int,
	bufferCount int,
	interpolation int,
	decimation int,
	input chan *buffers.IQ,
	quit chan struct{},
) (int, chan *buffers.IQ) {
	var resampler = dsp.MakeRationalResampler(interpolation, decimation)
	var resampleRate = resampler.PredictOutputSize(sampleRate)
	var outputRing = buffers.NewIQRing(resampleRate, bufferCount)
	var output = make(chan *buffers.IQ, bufferCount-1)
	go iqRationalResamplerLoop(sampleRate, resampler, outputRing, input, output, quit)
	log.WithFields(log.Fields{
		"sampleRate":   sampleRate,
		"resampleRate": resampleRate,
	}).Debug("IQRationalResampler")
	return resampleRate, output
}

func iqRationalResamplerLoop(
	sampleRate int,
	resampler *dsp.RationalResampler,
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
	quit chan struct{},
) {
	log.WithField("filter", "IQRationalResampler").Debug("Starting")
	var inBuffer = buffers.NewIQ(sampleRate)
	for {
		select {
		case <-quit:
			close(output)
			log.WithField("filter", "IQRationalResampler").Debug("Exiting")
			runtime.GC()
			return
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQRationalResampler").Debug("Exiting")
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
