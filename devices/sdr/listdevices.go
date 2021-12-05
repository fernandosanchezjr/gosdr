package sdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/rtlsdr"
)

func ListDevices() []*devices.DeviceInfo {
	return rtlsdr.ListDevices()
}
