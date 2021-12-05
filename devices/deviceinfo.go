package devices

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type DeviceType int

const (
	RTLSDR = iota
)

type DeviceInfo struct {
	Type         DeviceType
	Index        int
	Manufacturer string
	ProductName  string
	Serial       string
}

func (di *DeviceInfo) Fields() logrus.Fields {
	return logrus.Fields{
		"index":        di.Index,
		"manufacturer": di.Manufacturer,
		"productName":  di.ProductName,
		"serial":       di.Serial,
	}
}

func (di *DeviceInfo) Equals(other *DeviceInfo) bool {
	return di.Type == other.Type && di.Index == other.Index && di.Manufacturer == other.Manufacturer &&
		di.ProductName == other.ProductName && di.Serial == other.Serial
}

func (di *DeviceInfo) String() string {
	return fmt.Sprintf("(%d) %s - %s: %s", di.Index, di.Manufacturer, di.ProductName, di.Serial)
}
