package rtlsdr

import (
	rtl "github.com/fernandosanchezjr/gortlsdr"
	"github.com/fernandosanchezjr/gosdr/devices"
	log "github.com/sirupsen/logrus"
)

func GetInfo(index int) (*devices.Info, error) {
	var manufacturer, product, serial, err = rtl.GetDeviceUsbStrings(index)
	if err != nil {
		log.WithError(err).WithField("index", index).Warn("rtl.GetDeviceUsbStrings")
		return nil, err
	}
	return &devices.Info{
		Type:         devices.RTLSDR,
		Index:        index,
		Manufacturer: manufacturer,
		ProductName:  product,
		Serial:       serial,
	}, nil
}
