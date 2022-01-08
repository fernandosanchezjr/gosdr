package filters

import (
	"bytes"
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/units"
	"github.com/racerxdl/segdsp/tools"
	log "github.com/sirupsen/logrus"
	chart "github.com/wcharczuk/go-chart/v2"
	"image"
	"math"
	"runtime"
)

func getFrequencies(steps int, lower, upper units.Hertz) []float64 {
	var frequencyStep = (upper - lower) / units.Hertz(steps)
	var frequencies = make([]float64, steps)
	var current = lower
	for i := range frequencies {
		frequencies[i] = current.Float()
		current += frequencyStep
	}
	return frequencies
}

func NewIQHistogram(sampleRate int, count int, conn devices.Connection, input chan *buffers.IQ) {
	var output = make(chan *bytes.Buffer, count)
	var window = app.NewWindow(app.Title("GOSDR IQ Histogram"))
	var lower, upper = conn.GetFrequencyBounds()
	var ring = buffers.NewIQRing(sampleRate, count)
	var bufferRing = buffers.NewBufferRing(count)
	var histogramFrequencies = getFrequencies(sampleRate, lower, upper)
	go iqHistogramLoop(window, histogramFrequencies, ring, bufferRing, input, output)
	go iqHistogramWindowLoop(window, output)
}

func histogramNormalize(value float64, max float64, log10 float64) (normalized float64) {
	if value == 0.0 {
		return 0.0
	}
	normalized = log10 * (value / max)
	if math.IsNaN(normalized) {
		return 0.0
	}
	return
}

func histogramIQtoFloat(log10 float64, input []complex64, histogram []float64) {
	var magnitude, maxMagnitude float64
	for i, value := range input {
		magnitude = float64(tools.ComplexAbs(value))
		histogram[i] = magnitude
		maxMagnitude = math.Max(maxMagnitude, magnitude)
	}
	for i, value := range histogram {
		histogram[i] = histogramNormalize(value, maxMagnitude, log10)
	}
}

func iqHistogramLoop(
	window *app.Window,
	frequencies []float64,
	ring *buffers.IQRing,
	bufferRing *buffers.BufferRing,
	input chan *buffers.IQ,
	output chan *bytes.Buffer,
) {
	log.WithField("filter", "IQHistogram").Debug("Starting")
	var log10 = 10.0 * math.Log10(10.0)
	var histogram = make([]float64, len(frequencies))
	for {
		select {
		case in, ok := <-input:
			if !ok {
				window.Close()
				log.WithField("filter", "IQHistogram").Debug("Exiting")
				runtime.GC()
				return
			}
			var out = ring.Next()
			out.Copy(in)
			histogramIQtoFloat(log10, out.Data(), histogram)
			var outBuf = bufferRing.Next()
			graph := chart.Chart{
				Series: []chart.Series{
					chart.ContinuousSeries{
						XValues: frequencies,
						YValues: histogram,
					},
				},
			}
			if renderErr := graph.Render(chart.PNG, outBuf); renderErr != nil {
				log.WithError(renderErr).Error("graph.Render")
			} else {
				output <- outBuf
			}
		}
	}
}

func iqHistogramWindowLoop(window *app.Window, input chan *bytes.Buffer) {
	var chartImage image.Image
	var decodeErr error
	var ops op.Ops
	for {
		select {
		case in := <-input:
			chartImage, _, decodeErr = image.Decode(in)
			if decodeErr != nil {
				log.WithError(decodeErr).Error("image.Decode")
			} else {
				var size = chartImage.Bounds().Max
				var x = unit.Dp(float32(size.X))
				var y = unit.Dp(float32(size.Y))
				window.Option(app.MaxSize(x, y), app.MinSize(x, y))
				window.Invalidate()
			}
		case e := <-window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				// A request to draw the window state.
				gtx := layout.NewContext(&ops, e)
				if chartImage != nil {
					imageOp := paint.NewImageOp(chartImage)
					imageOp.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
				}
				// Update the display.
				e.Frame(gtx.Ops)
			}
		}
	}
}
