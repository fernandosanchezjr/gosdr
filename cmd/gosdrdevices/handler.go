package main

import (
	"flag"
	"fmt"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"os"
	"sort"
)

func handleInvocation() {
	switch flag.Arg(0) {
	case "index":
		handleIndex()
	case "serial":
		handleSerial()
	case "listSerials":
		handleListSerials()
	case "count":
		handleCount()
	default:
		os.Exit(1)
	}
}

func handleIndex() {
	if deviceSerial == "" {
		os.Exit(2)
	}
	for _, device := range sdr.ListDevices() {
		if device.Serial == deviceSerial {
			fmt.Println(device.Index)
			os.Exit(0)
		}
	}
	os.Exit(3)
}

func handleSerial() {
	if deviceIndex < 0 {
		os.Exit(4)
	}
	for _, device := range sdr.ListDevices() {
		if device.Index == deviceIndex {
			fmt.Println(device.Serial)
			os.Exit(0)
		}
	}
	os.Exit(5)
}

func handleListSerials() {
	devices := sdr.ListDevices()

	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].Index < devices[j].Index
	})

	for _, device := range devices {
		fmt.Println(device.Serial)
	}
}

func handleCount() {
	fmt.Println(len(sdr.ListDevices()))
}
