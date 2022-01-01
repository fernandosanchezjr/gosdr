package filters

import (
	"github.com/fernandosanchezjr/gosdr/buffers"
)

func NewDecimator(sampleRate, decimationRate, count int, input chan *buffers.IQ) (int, chan *buffers.IQ) {
	var output = make(chan *buffers.IQ, count)
	var outputSize = sampleRate / decimationRate
	var ring = buffers.NewIQRing(outputSize, count)
	go decimatorLoop(decimationRate, ring, input, output)
	return outputSize, output
}

func decimatorLoop(decimationRate int, ring *buffers.IQRing, input, output chan *buffers.IQ) {
	for {
		select {
		case in, ok := <-input:
			if !ok {
				close(output)
				return
			}
			var out = ring.Next()
			var inputSamples = in.Data()
			var outputSamples = out.Data()
			var inputPos, outputPos = 0, 0
			for inputPos < len(inputSamples) {
				outputSamples[outputPos] = inputSamples[inputPos]
				inputPos += decimationRate
				outputPos += 1
			}
			output <- out
		}
	}
}
