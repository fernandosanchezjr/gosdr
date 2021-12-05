package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	rtl "github.com/jpoirier/gortlsdr"
	log "github.com/sirupsen/logrus"
)

func GetDeviceInfo(index int) (*devices.DeviceInfo, error) {
	var manufacturer, product, serial, err = rtl.GetDeviceUsbStrings(index)
	if err != nil {
		log.WithError(err).WithField("index", index).Warn("rtl.GetDeviceUsbStrings")
		return nil, err
	}
	return &devices.DeviceInfo{
		Type:         devices.RTLSDR,
		Index:        index,
		Manufacturer: manufacturer,
		ProductName:  product,
		Serial:       serial,
	}, nil
}
