package main

import "flag"

var deviceSerial string
var deviceIndex int

func init() {
	flag.IntVar(&deviceIndex, "index", -1, "device index")
	flag.StringVar(&deviceSerial, "serial", deviceSerial, "device serial")
}
