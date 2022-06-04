package main

import "flag"

var rtlSDR bool = true
var bufferSeconds = 30

func init() {
	flag.BoolVar(&rtlSDR, "rtl", rtlSDR, "use RTL-SDR devices")
	flag.IntVar(&bufferSeconds, "buffer-time-seconds", bufferSeconds, "buffer time in seconds")
}
