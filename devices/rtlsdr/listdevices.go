package rtlsdr

import (
	rtlsdr "github.com/fernandosanchezjr/gortlsdr"
	"github.com/fernandosanchezjr/gosdr/devices"
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
