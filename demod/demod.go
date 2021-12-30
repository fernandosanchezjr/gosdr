package demod

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/units"
)

type Demodulator interface {
	ConnectionsRequired() int
	CenterFrequency() units.Hertz
	UseConnection(connection devices.Connection) error
	Connections() []devices.Connection
	Start() error
	Stop() error
	End() error
}
