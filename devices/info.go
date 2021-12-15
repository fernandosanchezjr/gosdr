package devices

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type DeviceType int

const (
	RTLSDR DeviceType = iota
)

type Id struct {
	Type   DeviceType
	Serial string
}

type Info struct {
	Type         DeviceType
	Index        int
	Manufacturer string
	ProductName  string
	Serial       string
}

func (dt DeviceType) String() string {
	switch dt {
	case RTLSDR:
		return "RTL-SDR"
	default:
		return "Unknown"
	}
}

func (di *Info) Fields() logrus.Fields {
	return logrus.Fields{
		"type":         di.Type,
		"index":        di.Index,
		"manufacturer": di.Manufacturer,
		"productName":  di.ProductName,
		"serial":       di.Serial,
	}
}

func (di *Info) Equals(other *Info) bool {
	return di.Type == other.Type && di.Index == other.Index && di.Manufacturer == other.Manufacturer &&
		di.ProductName == other.ProductName && di.Serial == other.Serial
}

func (di *Info) String() string {
	return fmt.Sprintf("(%d) %s - %s: %s", di.Index, di.Manufacturer, di.ProductName, di.Serial)
}

func (di *Info) Id() Id {
	return Id{
		Type:   di.Type,
		Serial: di.Serial,
	}
}

func (i Id) Fields() logrus.Fields {
	return logrus.Fields{
		"type":   i.Type,
		"serial": i.Serial,
	}
}
