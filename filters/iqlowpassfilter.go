package filters

//
//import (
//	"github.com/fernandosanchezjr/gosdr/buffers"
//	"github.com/racerxdl/segdsp/dsp"
//	log "github.com/sirupsen/logrus"
//	"runtime"
//)
//
//func NewIQLowpassFilter(
//	sampleRate int,
//	lastBufferBlock int,
//	gain float64,
//	cutFrequency float64,
//	transitionWidth float64,
//	input chan *buffers.IQ,
//	quit chan struct{},
//) chan *buffers.IQ {
//	var filter = dsp.MakeFirFilter(dsp.MakeLowPass(gain, float64(sampleRate), cutFrequency, transitionWidth))
//	var outputRing = buffers.NewIQRing(sampleRate, lastBufferBlock)
//	var output = make(chan *buffers.IQ, lastBufferBlock-1)
//	go iqLowpassFilterLoop(sampleRate, filter, outputRing, input, output, quit)
//	log.WithFields(log.Fields{
//		"sampleRate": sampleRate,
//	}).Debug("IQLowpassFilter")
//	return output
//}
//
//func iqLowpassFilterLoop(
//	sampleRate int,
//	filter *dsp.FirFilter,
//	outputRing *buffers.IQRing,
//	input chan *buffers.IQ,
//	output chan *buffers.IQ,
//	quit chan struct{},
//) {
//	log.WithField("filter", "IQLowpassFilter").Debug("Starting")
//	var tmpBuffer = buffers.NewIQ(sampleRate)
//	for {
//		select {
//		case <-quit:
//			close(output)
//			log.WithField("filter", "IQLowpassFilter").Debug("Exiting")
//			runtime.GC()
//			return
//		case in, ok := <-input:
//			if !ok {
//				close(output)
//				log.WithField("filter", "IQLowpassFilter").Debug("Exiting")
//				runtime.GC()
//				return
//			}
//			tmpBuffer.Copy(in)
//			var out = outputRing.Next()
//			filter.FilterBuffer(tmpBuffer.Data(), out.Data())
//			out.Sequence = tmpBuffer.Sequence
//			output <- out
//		}
//	}
//}
