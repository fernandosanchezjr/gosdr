package demod

import "github.com/fernandosanchezjr/gosdr/devices"

type Demodulator interface {
	ConnectionsRequired() int
	UseConnection(connection devices.Connection) error
	Connections() []devices.Connection
	Start() error
	Stop() error
	End() error
}
