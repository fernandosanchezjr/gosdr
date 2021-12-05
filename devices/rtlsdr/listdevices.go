package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/jpoirier/gortlsdr"
)

func ListDevices() []*devices.DeviceInfo {
	var result []*devices.DeviceInfo
	var count = rtlsdr.GetDeviceCount()
	for i := 0; i < count; i++ {
		if info, err := GetDeviceInfo(i); err != nil {
			continue
		} else {
			result = append(result, info)
		}
	}
	return result
}
