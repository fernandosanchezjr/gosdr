package main

import (
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages"
	sdrPage "github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages/sdr"
)

func main() {
	config.ParseFlags()
	config.SetupLogger()
	var sdrManager = sdr.NewManager()
	go func() {
		w := app.NewWindow(app.Title("GOSDR"))
		if err := loop(w, sdrManager); err != nil {
			log.WithError(err).Fatal("Exiting")
		}
		sdrManager.Stop()
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, sdrManager *sdr.Manager) error {
	var th = themes.NewTheme()
	var ops op.Ops
	var router = pages.NewRouter(th, sdrManager)

	router.Register(0, sdrPage.New(&router))

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				router.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		case deviceEvent := <-sdrManager.DeviceChan:
			log.WithFields(deviceEvent.Fields()).Info("Received device event")
			switch deviceEvent.EventType {
			case sdr.DeviceRemoved:
				router.State.RemoveDevice(deviceEvent.Id)
			case sdr.DeviceAdded:
				router.State.AddDevice(deviceEvent.Id)
			}
			w.Invalidate()
		}
	}
}
