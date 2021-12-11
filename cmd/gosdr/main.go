package main

import (
	"flag"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
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
	flag.Parse()
	var scanner = sdr.NewScanner()
	go func() {
		w := app.NewWindow(app.Title("GOSDR"))
		if err := loop(w, scanner); err != nil {
			log.WithError(err).Fatal("Exiting")
		}
		scanner.Close()
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, scanner *sdr.Scanner) error {
	var th = themes.NewTheme()
	var ops op.Ops
	var router = pages.NewRouter()

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
		case devices := <-scanner.DeviceChan:
			log.WithField("count", len(devices)).Println("Device change notification")
			router.SetDevices(devices)
			w.Invalidate()
		}
	}
}
