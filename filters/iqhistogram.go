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

func NewIQHistogram(
	sampleRate int,
	bufferCount int,
	conn devices.Connection,
	bandwidth units.Hertz,
	input chan *buffers.IQ,
	quit chan struct{},
) {
	var output = make(chan *bytes.Buffer, bufferCount-1)
	var ring = buffers.NewIQRing(sampleRate, bufferCount)
	var bufferRing = buffers.NewBufferRing(bufferCount)
	var center = conn.GetCenterFrequency()
	var lower = center - (bandwidth / 2)
	var upper = lower + bandwidth
	var histogramFrequencies = getFrequencies(sampleRate, lower, upper)
	go iqHistogramLoop(histogramFrequencies, ring, bufferRing, input, output, quit)
	go iqHistogramWindowLoop(output, quit)
}

func calculatePower(sample complex64) float64 {
	var power = 10.0 * math.Log10(float64(tools.ComplexAbsSquared(sample)))
	if math.IsInf(power, 0) || math.IsNaN(power) {
		return 0.0
	}
	return power
}

func histogramIQtoFloat(input []complex64, histogram []float64) {
	for i, value := range input {
		histogram[i] = calculatePower(value)
	}
}

func frequencyFormatter(v interface{}) string {
	return units.Hertz(v.(float64)).String()
}

func iqHistogramLoop(
	frequencies []float64,
	ring *buffers.IQRing,
	bufferRing *buffers.BufferRing,
	input chan *buffers.IQ,
	output chan *bytes.Buffer,
	quit chan struct{},
) {
	log.WithField("filter", "IQHistogram").Debug("Starting")
	var histogram = make([]float64, len(frequencies))
	for {
		select {
		case <-quit:
			close(output)
			log.WithField("filter", "IQHistogram").Debug("Exiting")
			runtime.GC()
			return
		case in, ok := <-input:
			if !ok {
				close(output)
				log.WithField("filter", "IQHistogram").Debug("Exiting")
				runtime.GC()
				return
			}
			var out = ring.Next()
			out.Copy(in)
			histogramIQtoFloat(out.Data(), histogram)
			var outBuf = bufferRing.Next()
			graph := chart.Chart{
				Series: []chart.Series{
					chart.ContinuousSeries{
						XValues:         frequencies,
						YValues:         histogram,
						XValueFormatter: frequencyFormatter,
					},
				},
				YAxis: chart.YAxis{
					Name:      "",
					NameStyle: chart.Style{},
					Style:     chart.Style{},
					Zero:      chart.GridLine{},
					AxisType:  0,
					Ascending: false,
					Range: &chart.ContinuousRange{
						Min:        -160,
						Max:        60,
						Domain:     len(histogram),
						Descending: false,
					},
					TickStyle:      chart.Style{},
					GridMajorStyle: chart.Style{},
					GridMinorStyle: chart.Style{},
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

func iqHistogramWindowLoop(input chan *bytes.Buffer, quit chan struct{}) {
	var window = app.NewWindow(app.Title("GOSDR IQ Histogram"))
	var chartImage image.Image
	var decodeErr error
	var ops op.Ops
	var shouldQuit = true
	var closed bool
	var sized bool
	for {
		select {
		case e := <-window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				if shouldQuit {
					close(quit)
				}
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
			continue
		case <-quit:
			if closed {
				continue
			}
			shouldQuit = false
			window.Perform(system.ActionClose)
			window.Invalidate()
			closed = true
			continue
		case in, ok := <-input:
			if !ok {
				if !closed {
					shouldQuit = false
					window.Perform(system.ActionClose)
					go window.Invalidate()
					closed = true
				}
				continue
			}
			chartImage, _, decodeErr = image.Decode(in)
			if decodeErr != nil {
				log.WithError(decodeErr).Error("image.Decode")
			} else {
				if !sized {
					var size = chartImage.Bounds().Max
					var x = unit.Dp(float32(size.X))
					var y = unit.Dp(float32(size.Y))
					window.Option(app.MaxSize(x, y), app.MinSize(x, y))
					sized = true
				}
				window.Invalidate()
			}
		}
	}
}
