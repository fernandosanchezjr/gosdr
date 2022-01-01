package main

import (
	"gioui.org/app"
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
)

func applicationLoop(manager *sdr.Manager) {
	var selectedDevice = make(chan devices.Id, 64)
	go deviceSelector(manager.DeviceChan, selectedDevice)
	go deviceController(manager, selectedDevice)
}

func main() {
	config.ParseFlags(parseFrequency)
	config.SetupLogger()
	var manager = sdr.NewManager()
	defer manager.Stop()
	go applicationLoop(manager)
	app.Main()
}
