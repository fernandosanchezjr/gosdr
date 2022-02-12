package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp"
	"github.com/racerxdl/segdsp/dsp/fft"
	log "github.com/sirupsen/logrus"
	"runtime"
)

func NewIQFFT(
	fftSize int,
	sampleRate int,
	bufferCount int,
	input chan *buffers.IQ,
	quit chan struct{},
) chan *buffers.IQ {
	var output = make(chan *buffers.IQ, bufferCount-1)
	var outputRing = buffers.NewIQRing(fftSize, bufferCount)
	var window = dsp.BlackmanHarris(fftSize, 92)
	go iqFFTLoop(fftSize, sampleRate, window, outputRing, input, output, quit)
	log.WithFields(log.Fields{
		"sampleRate": sampleRate,
	}).Debug("IQFFT")
	return output
}

func computeFFT(midPoint int, input []complex64, output []complex64) {
	copy(input, fft.FFT(input))
	copy(output[0:midPoint], input[midPoint:])
	copy(output[midPoint+1:], input[0:midPoint-1])
}

func iqFFTLoop(
	fftSize int,
	sampleRate int,
	window []float64,
	outputRing *buffers.IQRing,
	input chan *buffers.IQ,
	output chan *buffers.IQ,
	quit chan struct{},
) {
	log.WithField("filter", "IQFFT").Debug("Starting")
	var midPoint = (fftSize / 2) + (fftSize % 2)
	var fftBuf = make([]complex64, fftSize)
	var currentIQ = buffers.NewIQ(sampleRate)
	for {
		select {
		case <-quit:
			close(output)
			log.WithField("filter", "IQFFT").Debug("Exiting")
			runtime.GC()
			return
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQFFT").Debug("Exiting")
				runtime.GC()
				return
			}
			currentIQ.Copy(in)
			currentIQ.Reset()
		default:
		}
		if _, err := currentIQ.Write(fftBuf); err == nil {
			computeWindow(fftBuf, window)
			var out = outputRing.Next()
			computeFFT(midPoint, fftBuf, out.Data())
			out.Sequence = currentIQ.Sequence
			output <- out
		}
	}
}
