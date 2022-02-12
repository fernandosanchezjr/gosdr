package main

import (
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	graphBufferCount = 4
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
				logrus.WithFields(event.Id.Fields()).Trace("Selected device")
				deviceIds <- event.Id
				continue
			}
			if deviceIndex == event.Index {
				logrus.WithFields(event.Id.Fields()).WithField("index", deviceIndex).Trace(
					"Selected",
				)
				deviceIds <- event.Id
				continue
			}
		}
	}
}

func closeConnection(deviceId devices.Id, manager *sdr.Manager) {
	manager.Close(deviceId)
}

func connectionQuitter(deviceId devices.Id, manager *sdr.Manager, quit chan struct{}) {
	select {
	case <-quit:
		closeConnection(deviceId, manager)
		manager.Stop()
		os.Exit(0)
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
				closeConnection(id, manager)
				continue
			}
			if autoGainErr := conn.SetAutoGain(autoGain); autoGainErr != nil {
				logrus.WithFields(conn.Fields()).WithError(autoGainErr).Error("conn.SetAutoGain")
				closeConnection(id, manager)
				continue
			}
			if gainErr := conn.SetTunerGain(float32(gain)); gainErr != nil {
				logrus.WithFields(conn.Fields()).WithError(gainErr).Error("conn.SetTunerGain")
				closeConnection(id, manager)
				continue
			}
			if ppmErr := conn.SetFrequencyCorrection(ppm); ppmErr != nil {
				logrus.WithFields(conn.Fields()).WithError(ppmErr).Error("conn.SetFrequencyCorrection")
				closeConnection(id, manager)
				continue
			}
			if freqErr := conn.SetCenterFrequency(frequency); freqErr != nil {
				logrus.WithFields(conn.Fields()).WithError(freqErr).Error("conn.SetCenterFrequency")
				closeConnection(id, manager)
				continue
			}
			if resetErr := conn.Reset(); resetErr != nil {
				logrus.WithFields(conn.Fields()).WithError(resetErr).Error("conn.Reset")
				closeConnection(id, manager)
				continue
			}
			var quit = make(chan struct{})
			var input = createGraph(conn, graphBufferCount, quit)
			manager.AddDeviceCleanup(id, func() {
				close(input)
			})
			go connectionQuitter(id, manager, quit)
			go sampleDevice(conn, input)
			logrus.WithFields(conn.Fields()).Info("Sampling SDR")
		}
	}
}
