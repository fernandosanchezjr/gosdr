package demod

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/utils"
)

type Demodulator interface {
	ConnectionsRequired() int
	CenterFrequency() utils.Hertz
	UseConnection(connection devices.Connection) error
	Connections() []devices.Connection
	Start() error
	Stop() error
	End() error
}
