package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"github.com/fernandosanchezjr/gosdr/utils"
)

func main() {
	parseFlags()
	var selectedDevice = make(chan devices.Id, 64)
	var manager = sdr.NewManager()
	defer manager.Stop()
	go deviceSelector(manager.DeviceChan, selectedDevice)
	go deviceController(manager, selectedDevice)
	utils.Wait()
}
