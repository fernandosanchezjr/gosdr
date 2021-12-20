package wbfm

import "github.com/fernandosanchezjr/gosdr/devices"

type WBFMDemodulator struct {
	Connection devices.Connection
}

func NewWBFMDemodulator() *WBFMDemodulator {
	return &WBFMDemodulator{}
}

func (wb *WBFMDemodulator) ConnectionsRequired() int {
	return 1
}

func (wb *WBFMDemodulator) UseConnection(connection devices.Connection) error {
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
