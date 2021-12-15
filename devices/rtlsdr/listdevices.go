package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/jpoirier/gortlsdr"
)

func ListDevices() []*devices.Info {
	var result []*devices.Info
	var count = rtlsdr.GetDeviceCount()
	for i := 0; i < count; i++ {
		if info, err := GetInfo(i); err != nil {
			continue
		} else if info.Valid() {
			result = append(result, info)
		}
	}
	return result
}
