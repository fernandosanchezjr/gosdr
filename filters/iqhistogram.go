package filters

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/units"
	"github.com/racerxdl/segdsp/tools"
	"image/color"
	"math"
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

func generateFrequencyMap(sampleRate int, conn devices.Connection, bandwidth units.Hertz) []float64 {
	var center = conn.GetCenterFrequency()
	var lower = center - (bandwidth / 2)
	var upper = lower + bandwidth
	return getFrequencies(sampleRate, lower, upper)
}

func NewIQHistogram(
	sampleRate int,
	bufferCount int,
	conn devices.Connection,
	bandwidth units.Hertz,
	input chan *buffers.IQ,
	quit chan struct{},
) {
	var ring = buffers.NewIQRing(sampleRate, bufferCount)
	generateFrequencyMap(sampleRate, conn, bandwidth)
	go iqHistogramWindowLoop(sampleRate, ring, input, quit)
}

func getPower(sample complex64) float64 {
	if imag(sample) == 0.0 {
		return 0.0
	}
	var value = 10.0 * math.Log10(float64(tools.ComplexAbsSquared(sample)))
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0.0
	}
	return value
}

func getNormalizedValue(value float64, min float64, powerRange float64) float64 {
	switch value {
	case 0.0:
		return 150.0
	default:
		return 150.0 - ((value-min)/powerRange)*150.0
	}
}

func normalizePower(input []complex64, histogram []float64) {
	var power, min, max, powerRange float64
	for i, value := range input {
		power = getPower(value)
		min = math.Min(min, power)
		max = math.Max(max, power)
		histogram[i] = power
	}
	powerRange = max - min
	for i, value := range histogram {
		histogram[i] = getNormalizedValue(value, min, powerRange)
	}
}

func generatePath(sampleRate int, gtx layout.Context, histogram []float64) clip.PathSpec {
	var path = &clip.Path{}
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(-1.0, 151.0))
	for pos, value := range histogram {
		path.LineTo(f32.Pt(float32(pos), float32(value)))
	}
	path.LineTo(f32.Pt(float32(sampleRate+1), float32(histogram[len(histogram)-1])))
	path.LineTo(f32.Pt(float32(sampleRate+1), 151.0))
	path.LineTo(f32.Pt(-1.0, 151.0))
	path.Close()
	var pathSpec = path.End()
	return pathSpec
}

func drawHistogram(sampleRate int, gtx layout.Context, histogram []float64) {
	paint.Fill(gtx.Ops, color.NRGBA{A: 0xff})
	paint.ColorOp{
		Color: color.NRGBA{R: 0xcc, G: 0xcc, A: 0xee},
	}.Add(gtx.Ops)

	var path = generatePath(sampleRate, gtx, histogram)

	paint.FillShape(
		gtx.Ops,
		color.NRGBA{R: 0xcc, G: 0xcc, A: 0xff},
		clip.Stroke{
			Path:  path,
			Width: 1.0,
		}.Op(),
	)

	paint.FillShape(
		gtx.Ops,
		color.NRGBA{A: 0xff},
		clip.Outline{
			Path: path,
		}.Op(),
	)
}

func iqHistogramWindowLoop(sampleRate int, ring *buffers.IQRing, input chan *buffers.IQ, quit chan struct{}) {
	var window = app.NewWindow(app.Title("GOSDR IQ Histogram"))
	var width, height = unit.Dp(float32(sampleRate)), unit.Dp(float32(150))
	var ops op.Ops
	var shouldQuit = true
	var closed bool
	var histogram = make([]float64, sampleRate)
	window.Option(app.MaxSize(width, height), app.MinSize(width, height))
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
				gtx := layout.NewContext(&ops, e)
				drawHistogram(sampleRate, gtx, histogram)
				e.Frame(gtx.Ops)
			}
		case <-quit:
			if closed {
				continue
			}
			shouldQuit = false
			window.Perform(system.ActionClose)
			window.Invalidate()
			closed = true
		case in, ok := <-input:
			if !ok {
				if !closed {
					shouldQuit = false
					window.Perform(system.ActionClose)
					window.Invalidate()
					closed = true
				}
				continue
			}
			var out = ring.Next()
			out.Copy(in)
			normalizePower(out.Data(), histogram)
			window.Invalidate()
		}
	}
}
