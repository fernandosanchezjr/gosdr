package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/fft"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQFFT(sampleRate int, bufferCount int, input chan *buffers.IQ) chan *buffers.IQ {
	var output = make(chan *buffers.IQ, bufferCount)
	var window = dsp.BlackmanHarris(sampleRate, 92)
	var outputRing = buffers.NewIQRing(sampleRate, bufferCount)
	go iqFFTLoop(sampleRate, window, outputRing, input, output)
	log.WithFields(log.Fields{
		"sampleRate": sampleRate,
	}).Debug("IQFFT")
	return output
}

func computeFFT(midPoint int, input []complex64, output []complex64) {
	//copy(output, fft.FFT(input))
	copy(input, fft.FFT(input))
	copy(output[0:midPoint], input[midPoint:])
	copy(output[midPoint+1:], input[0:midPoint-1])
}

func iqFFTLoop(
	sampleRate int,
	window []float64,
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
) {
	log.WithField("filter", "IQFFT").Debug("Starting")
	var midPoint = (sampleRate / 2) + (sampleRate % 2)
	var fftBuffer = buffers.NewIQ(sampleRate)
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQFFT").Debug("Exiting")
				runtime.GC()
				return
			}
			fftBuffer.Copy(in)
			var out = outputRing.Next()
			computeWindow(fftBuffer.Data(), window)
			computeFFT(midPoint, fftBuffer.Data(), out.Data())
			out.Sequence = fftBuffer.Sequence
			output <- out
		}
	}
}
