package devices

import (
	"github.com/fernandosanchezjr/gosdr/utils"
	log "github.com/sirupsen/logrus"
)

type Connection interface {
	Close() error
	IsOpen() bool
	Refresh() error
	Fields() log.Fields
	GetInfo() *Info
	GetAGC() bool
	SetAGC(enabled bool) error
	GetAutoGain() bool
	SetAutoGain(enabled bool) error
	GetTunerGain() float32
	SetTunerGain(gain float32) error
	GetFrequencyCorrection() int
	SetFrequencyCorrection(ppm int) error
	Reset() error
	GetCenterFrequency() utils.Hertz
	SetCenterFrequency(centerFrequency utils.Hertz) error
}
