package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"github.com/sirupsen/logrus"
)

func deviceSelector(events chan sdr.DeviceEvent, deviceIds chan devices.Id) {
	for {
		select {
		case event := <-events:
			if event.EventType != sdr.DeviceAdded {
				continue
			}
			if !rtlSDR && event.Id.Type == devices.RTLSDR {
				continue
			}
			if deviceSerial != "" && deviceSerial == event.Id.Serial {
				logrus.WithFields(event.Id.Fields()).Info("Selected device by serial")
				deviceIds <- event.Id
				continue
			}
			if deviceIndex == event.Index {
				logrus.WithFields(event.Id.Fields()).WithField("index", deviceIndex).Info(
					"Selected device by index",
				)
				deviceIds <- event.Id
				continue
			}
		}
	}
}

func closeConnection(conn devices.Connection) {
	if err := conn.Close(); err != nil {
		logrus.WithFields(conn.Fields()).WithError(err).Error("conn.Close")
	}
}

func deviceController(manager *sdr.Manager, deviceIds chan devices.Id) {
	for {
		select {
		case id := <-deviceIds:
			var conn, connErr = manager.Open(id)
			if connErr != nil {
				logrus.WithFields(id.Fields()).WithError(connErr).Error("manager.Open")
				continue
			}
			if agcErr := conn.SetAGC(agc); agcErr != nil {
				logrus.WithFields(conn.Fields()).WithError(agcErr).Error("conn.SetAGC")
				closeConnection(conn)
				continue
			}
			if autoGainErr := conn.SetAutoGain(autoGain); autoGainErr != nil {
				logrus.WithFields(conn.Fields()).WithError(autoGainErr).Error("conn.SetAutoGain")
				closeConnection(conn)
				continue
			}
			if gainErr := conn.SetTunerGain(float32(gain)); gainErr != nil {
				logrus.WithFields(conn.Fields()).WithError(gainErr).Error("conn.SetTunerGain")
				closeConnection(conn)
				continue
			}
			if ppmErr := conn.SetFrequencyCorrection(ppm); ppmErr != nil {
				logrus.WithFields(conn.Fields()).WithError(ppmErr).Error("conn.SetFrequencyCorrection")
				closeConnection(conn)
				continue
			}
			if resetErr := conn.Reset(); resetErr != nil {
				logrus.WithFields(conn.Fields()).WithError(resetErr).Error("conn.Reset")
				closeConnection(conn)
				continue
			}
			if freqErr := conn.SetCenterFrequency(frequency); freqErr != nil {
				logrus.WithFields(conn.Fields()).WithError(freqErr).Error("conn.SetCenterFrequency")
				closeConnection(conn)
				continue
			}
			logrus.WithFields(conn.Fields()).Info("Opened device")
		}
	}
}