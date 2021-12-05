package main

import (
	"flag"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	page "github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages/deviceselection"
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
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	router := page.NewRouter()
	router.Register(0, deviceselection.New(&router))

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
