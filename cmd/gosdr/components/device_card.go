package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
)

const (
	connectText    = "CONNECT"
	disconnectText = "DISCONNECT"
)

func deviceTitle(gtx C, th *themes.Theme, device *devices.Info) D {
	return HorizontalList(gtx, 1,
		layout.Flexed(1,
			material.H6(
				th.Theme,
				fmt.Sprint(device.Manufacturer, " ", device.ProductName),
			).Layout,
		),
	)
}

func deviceSubTitle(gtx C, th *themes.Theme, device *devices.Info) D {
	return HorizontalList(gtx, 1,
		layout.Flexed(1,
			material.Subtitle1(
				th.Theme,
				fmt.Sprintf("Serial %s \u2014 Index %d", device.Serial, device.Index),
			).Layout,
		),
	)
}

func connectLabel(connected bool) string {
	if connected {
		return disconnectText
	} else {
		return connectText
	}
}

func DeviceCard(
	gtx C,
	th *themes.Theme,
	manager *sdr.Manager,
	device *devices.Info,
	connectButton *widget.Clickable,
) D {
	var connectClicked bool
	var id = device.Id()
	var isConnected = manager.IsConnected(id)
	for connectButton.Clicked() {
		connectClicked = true
	}
	if connectClicked {
		if isConnected {
			manager.Close(id)
		} else {
			go func() {
				if _, openErr := manager.Open(id); openErr != nil {
					log.WithFields(device.Fields()).WithError(openErr).Error("Open")
				}
			}()
		}
	}
	var widgets = []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return deviceTitle(gtx, th, device)
		}),
		layout.Rigid(func(gtx C) D {
			return deviceSubTitle(gtx, th, device)
		}),
		layout.Rigid(func(gtx C) D {
			return layout.Inset{
				Top:    unit.Dp(10),
				Right:  unit.Dp(0),
				Bottom: unit.Dp(10),
				Left:   unit.Dp(0),
			}.Layout(gtx, func(gtx C) D {
				var button = material.Button(th.Theme, connectButton, connectLabel(isConnected))
				if !isConnected {
					button.Background = th.Primary.Dark.Bg
				}
				return button.Layout(gtx)
			})
		}),
	}
	return Card(gtx, th, func(gtx C) D {
		return VerticalList(gtx, widgets...)
	})
}
