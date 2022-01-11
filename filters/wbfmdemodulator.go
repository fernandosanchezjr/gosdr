package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	log "github.com/sirupsen/logrus"
	"runtime"
)

const (
	wbfmWindowTaps      = 8192
	wbfmAttenuation     = 74
	wbfmDecimationRatio = 16
)

func NewWBFMDemodulator(
	iqSampleRate int,
	bufferCount int,
	input chan *buffers.IQ,
) (int, chan *buffers.IQ) {
	var output = make(chan *buffers.IQ, bufferCount)
	var filter = dsp.MakeFirFilter(
		dsp.MakeLowPass(20.0, float64(iqSampleRate), 50_000, 1_000),
	)
	var decimator = dsp.MakeDecimator(wbfmDecimationRatio)
	var decimatedRate = decimator.PredictOutputSize(iqSampleRate)
	log.WithField("decimatedRate", decimatedRate).Debug("Decimator")
	var outputRing = buffers.NewIQRing(decimatedRate, bufferCount)
	go wbfmDemodulatorLoop(iqSampleRate, filter, decimator, outputRing, input, output)
	return decimatedRate, output
}

func wbfmDemodulatorLoop(
	iqSampleRate int,
	filter *dsp.FirFilter,
	decimator *dsp.Decimator,
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
) {
	var sequence uint64
	var tmpIn = buffers.NewIQ(iqSampleRate)
	var tmpFiltered = make([]complex64, filter.PredictOutputSize(iqSampleRate))
	log.WithField("filter", "Decimator").Debug("Starting")
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "Decimator").Debug("Exiting")
				runtime.GC()
				return
			}
			var out = outputRing.Next()
			tmpIn.Copy(in)
			filter.FilterBuffer(tmpIn.Data(), tmpFiltered)
			decimator.WorkBuffer(tmpFiltered, out.Data())
			sequence += 1
			out.Sequence = sequence
			output <- out
		}
	}
}
