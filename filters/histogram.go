package filters

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"github.com/fernandosanchezjr/gosdr/buffers"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"image/color"
	"math"
	"os"
	"sync/atomic"
)

var histogramId uint64

type histogramTypes interface {
	complex64 | complex128
}

type histogramState[T histogramTypes] struct {
	input         *buffers.Stream[T]
	histogramRing *buffers.BlockRing[float64]
	histogramChan chan []float64
	histogram     []float64
	height        float64
	pixelWidth    float32
	peakWidth     int
	size          int
	logger        *log.Entry
}

func NewHistogram[T histogramTypes](
	input *buffers.Stream[T],
	size int,
) {
	var id = atomic.AddUint64(&histogramId, 1)
	var histogramBlock = buffers.NewBlock[T](input.Size)
	var peakWidth = input.Size / size
	var filter = &histogramState[T]{
		input:         input,
		histogramRing: buffers.NewBlockRing[float64](size, input.Count),
		histogramChan: make(chan []float64, input.Count),
		histogram:     make([]float64, input.Size),
		height:        360,
		peakWidth:     peakWidth,
		size:          size,
		logger: log.WithFields(log.Fields{
			"filter": fmt.Sprintf("Histogram(%T)", histogramBlock.Data[0]),
			"id":     id,
		}),
	}
	go filter.loop()
	go filter.drawingLoop()
	return
}

func getNormalizedValue(value float64, min float64, powerRange float64) float64 {
	return 150.0 - ((value-min)/powerRange)*150.0
}

func normalizeHistogram(histogram []float64, min float64, powerRange float64, height float64) {
	for i, value := range histogram {
		histogram[i] = (getNormalizedValue(value, min, powerRange) / 150.0) * height
	}
}

func calculateNormalizedPower64(input []complex64, histogram []float64, height float64) {
	var power, min, max float64
	for i, value := range input {
		power = utils.GetPower64(value)
		min = math.Min(min, power)
		max = math.Max(max, power)
		histogram[i] = power
	}
	normalizeHistogram(histogram, min, max-min, height)
}

func calculateNormalizedPower128(input []complex128, histogram []float64, height float64) {
	var power, min, max float64
	for i, value := range input {
		power = utils.GetPower128(value)
		min = math.Min(min, power)
		max = math.Max(max, power)
		histogram[i] = power
	}
	normalizeHistogram(histogram, min, max-min, height)
}

func (filter *histogramState[T]) blockHandler(block *buffers.Block[T]) {
	filter.logger.WithField("block", block).Trace("Stream in")
	switch inBuf := any(block.Data).(type) {
	case []complex64:
		calculateNormalizedPower64(inBuf, filter.histogram, filter.height)
	case []complex128:
		calculateNormalizedPower128(inBuf, filter.histogram, filter.height)
	}
	var output = filter.histogramRing.Next()
	var minPower float64
	var histogramPos = 0
	for i := range output.Data {
		minPower = filter.height
		for j := 0; j < filter.peakWidth; j++ {
			minPower = math.Min(minPower, filter.histogram[histogramPos])
			histogramPos += 1
		}
		output.Data[i] = minPower
	}
	filter.histogramChan <- output.Data
}

func (filter *histogramState[T]) close() {
	filter.input.Done()
	close(filter.histogramChan)
	filter.input = nil
}

func (filter *histogramState[T]) loop() {
	var closed bool
	filter.logger.Debug("Starting")
	for {
		if closed = filter.input.Receive(filter.blockHandler); closed {
			filter.logger.Debug("Exiting")
			filter.close()
			return
		}
	}
}

func generatePath(sampleRate int, gtx layout.Context, histogram []float64, pixelWidth float32) clip.PathSpec {
	var path = &clip.Path{}
	path.Begin(gtx.Ops)
	path.MoveTo(f32.Pt(-1.0, float32(histogram[0])))
	for pos, value := range histogram {
		path.LineTo(f32.Pt(float32(pos)*pixelWidth, float32(value)))
	}
	path.LineTo(f32.Pt((float32(sampleRate)*pixelWidth)+1, float32(histogram[len(histogram)-1])))
	var pathSpec = path.End()
	return pathSpec
}

func drawHistogram(sampleRate int, gtx layout.Context, histogram []float64, pixelWidth float32) {
	paint.Fill(gtx.Ops, color.NRGBA{A: 0xff})
	paint.ColorOp{
		Color: color.NRGBA{R: 0xcc, G: 0xcc, A: 0xee},
	}.Add(gtx.Ops)
	var path = generatePath(sampleRate, gtx, histogram, pixelWidth)
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

func (filter *histogramState[T]) drawingLoop() {
	var window = app.NewWindow(app.Title("GOSDR IQ Histogram"))
	var width, height = unit.Dp(float32(filter.size)), unit.Dp(float32(filter.height))
	var histogram = make([]float64, filter.size)
	var ops op.Ops
	var inputClosing bool
	window.Option(app.Size(width, height))
	for {
		select {
		case e := <-window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				if !inputClosing {
					os.Exit(0)
				}
				return
			case system.FrameEvent:
				if filter.input == nil {
					continue
				}
				ops.Reset()
				filter.height = float64(e.Size.Y)
				filter.pixelWidth = float32(e.Size.X) / float32(filter.size)
				gtx := layout.NewContext(&ops, e)
				drawHistogram(filter.input.Size, gtx, histogram, filter.pixelWidth)
				e.Frame(gtx.Ops)
			}
		case histogramData, ok := <-filter.histogramChan:
			if !ok && !inputClosing {
				inputClosing = true
				window.Perform(system.ActionClose)
				continue
			}
			copy(histogram, histogramData)
			window.Invalidate()
		}
	}
}
