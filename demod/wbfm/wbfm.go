package wbfm

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/utils"
)

type WBFMDemodulator struct {
	Frequency  utils.Hertz
	Connection devices.Connection
}

func NewWBFMDemodulator(frequency utils.Hertz) *WBFMDemodulator {
	return &WBFMDemodulator{
		Frequency: frequency,
	}
}

func (wb *WBFMDemodulator) ConnectionsRequired() int {
	return 1
}

func (wb *WBFMDemodulator) CenterFrequency() utils.Hertz {
	return wb.Frequency
}

func (wb *WBFMDemodulator) UseConnection(connection devices.Connection) error {
	if freqErr := connection.SetCenterFrequency(wb.Frequency); freqErr != nil {
		return freqErr
	}
	wb.Connection = connection
	return nil
}

func (wb *WBFMDemodulator) Connections() []devices.Connection {
	return []devices.Connection{wb.Connection}
}

func (wb *WBFMDemodulator) Start() error {
	return nil
}

func (wb *WBFMDemodulator) Stop() error {
	return nil
}

func (wb *WBFMDemodulator) End() error {
	if wb.Connection != nil && wb.Connection.IsOpen() {
		if err := wb.Connection.Close(); err != nil {
			return err
		}
		wb.Connection = nil
	}
	return nil
}
