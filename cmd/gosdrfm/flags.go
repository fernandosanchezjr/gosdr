package main

import (
	"flag"
	"github.com/fernandosanchezjr/gosdr/units"
)

var rtlSDR bool
var deviceIndex int
var deviceSerial string
var requestedFrequency = "96.5MHz"
var frequency units.Hertz
var agc bool
var autoGain bool
var gain float64 = 1.0
var ppm int

func init() {
	flag.BoolVar(&rtlSDR, "rtl", rtlSDR, "use RTL-SDR devices")
	flag.IntVar(&deviceIndex, "index", deviceIndex, "device index")
	flag.StringVar(&deviceSerial, "serial", deviceSerial, "device serial")
	flag.StringVar(&requestedFrequency, "frequency", requestedFrequency, "tuner frequency in Hz")
	flag.BoolVar(&agc, "agc", agc, "enable AGC")
	flag.BoolVar(&autoGain, "auto-gain", autoGain, "enable auto-gain")
	flag.Float64Var(&gain, "gain", gain, "gain level in dB")
	flag.IntVar(&ppm, "frequency-correction", ppm, "frequency correction in PPM")
}

func parseFrequency() error {
	var frequencyErr error
	frequency, frequencyErr = units.ParseHertz(requestedFrequency)
	return frequencyErr
}
