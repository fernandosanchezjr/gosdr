package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/racerxdl/segdsp/dsp/fft"
	log "github.com/sirupsen/logrus"
)

func NewIQFFT(sampleRate int, bufferCount int, input chan *buffers.IQ) chan *buffers.IQ {
	var output = make(chan *buffers.IQ, bufferCount)
	var ring = buffers.NewIQRing(sampleRate, bufferCount)
	go iqFFTLoop(ring, sampleRate, input, output)
	return output
}

func computeFFT(midPoint int, input []complex64, output []complex64) {
	copy(input, fft.FFT(input))
	copy(output[0:midPoint], input[midPoint:])
	copy(output[midPoint+1:], input[0:midPoint-1])
}

func iqFFTLoop(outputRing *buffers.IQRing, size int, input chan *buffers.IQ, output chan *buffers.IQ) {
	var midPoint = size / 2
	var tmp = buffers.NewIQ(size)
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQFFT").Trace("Exiting")
				return
			}
			var out = outputRing.Next()
			tmp.Copy(in)
			computeFFT(midPoint, tmp.Data(), out.Data())
			out.Sequence = tmp.Sequence
			output <- out
		}
	}
}
