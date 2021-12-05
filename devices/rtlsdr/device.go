package rtlsdr

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	rtl "github.com/jpoirier/gortlsdr"
	log "github.com/sirupsen/logrus"
)

type Device struct {
	Info    *devices.DeviceInfo
	context *rtl.Context
}

func Open(index int) (*Device, error) {
	var context, openErr = rtl.Open(index)
	if openErr != nil {
		log.WithError(openErr).WithField("index", index).Error("rtl.Open")
		return nil, openErr
	}
	var info, infoErr = GetDeviceInfo(index)
	if infoErr != nil {
		return nil, infoErr
	}
	var device = &Device{
		Info:    info,
		context: context,
	}
	return device, nil
}

func (d *Device) Close() error {
	var err = d.context.Close()
	if err != nil {
		log.WithError(err).WithFields(d.Info.Fields()).Error("context.Close()")
	}
	return err
}

func (d *Device) Fields() log.Fields {
	return d.Info.Fields()
}
