package main

import (
	"flag"
	"fmt"
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
	"os"
)

var rtlSDR bool
var deviceIndex int
var deviceSerial string
var requestedFrequency = "96.5MHz"
var frequency utils.Hertz
var agc bool
var autoGain bool
var gain float64 = 1.0
var ppm int

func init() {
	flag.BoolVar(&rtlSDR, "rtl", rtlSDR, "use RTL-SDR devices")
	flag.IntVar(&deviceIndex, "index", deviceIndex, "device index")
	flag.StringVar(&deviceSerial, "serial", deviceSerial, "device serial")
	flag.StringVar(&requestedFrequency, "frequency", requestedFrequency, "tuner frequency in Hz")
	flag.BoolVar(&agc, "name", agc, "enable AGC")
	flag.BoolVar(&autoGain, "auto-gain", autoGain, "enable auto-gain")
	flag.Float64Var(&gain, "gain", gain, "gain level in dB")
	flag.IntVar(&ppm, "frequency-correction", ppm, "frequency correction in PPM")
}

func parseFrequency() error {
	var frequencyErr error
	frequency, frequencyErr = utils.ParseHertz(requestedFrequency)
	return frequencyErr
}

var customFlagHandlers = []func() error{
	parseFrequency,
}

func parseFlags() {
	flag.Parse()
	var failed = !flag.Parsed()
	for _, handler := range customFlagHandlers {
		if err := handler(); err != nil {
			log.WithError(err).Error("Argument error")
			failed = true
			break
		}
	}
	if failed {
		printUsage()
	}
}

func printUsage() {
	fmt.Println("gosdrfm")
	flag.PrintDefaults()
	os.Exit(1)
}
