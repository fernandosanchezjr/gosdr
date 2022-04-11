package device_card

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/components"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices"
	"github.com/fernandosanchezjr/gosdr/devices/rtlsdr"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	log "github.com/sirupsen/logrus"
)

const (
	connectText    = "CONNECT"
	disconnectText = "DISCONNECT"
)

func deviceTitle(gtx components.C, th *themes.Theme, device *devices.Info) components.D {
	return components.HorizontalList(gtx, 1,
		layout.Flexed(1,
			material.H6(
				th.Theme,
				fmt.Sprint(device.Manufacturer, " ", device.ProductName),
			).Layout,
		),
	)
}

func deviceSubTitle(gtx components.C, th *themes.Theme, device *devices.Info) components.D {
	return components.HorizontalList(gtx, 1,
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
	gtx components.C,
	th *themes.Theme,
	manager *sdr.Manager,
	device *devices.Info,
	connectButton *widget.Clickable,
	samplingMode *widget.Enum,
) components.D {
	var connectClicked bool
	var id = device.Id()
	var isConnected = manager.IsConnected(id)
	handleConnectButton(connectButton, connectClicked, isConnected, manager, id)
	var widgets = []layout.FlexChild{
		layout.Rigid(func(gtx components.C) components.D {
			return deviceTitle(gtx, th, device)
		}),
		layout.Rigid(func(gtx components.C) components.D {
			return deviceSubTitle(gtx, th, device)
		}),
	}
	widgets = append(widgets, generateDeviceDetails(isConnected, manager, id, th, samplingMode)...)
	widgets = append(widgets, layout.Rigid(func(gtx components.C) components.D {
		return layout.Inset{
			Top:    unit.Dp(10),
			Right:  unit.Dp(0),
			Bottom: unit.Dp(10),
			Left:   unit.Dp(0),
		}.Layout(gtx, func(gtx components.C) components.D {
			var button = material.Button(th.Theme, connectButton, connectLabel(isConnected))
			if !isConnected {
				button.Background = th.Primary.Dark.Bg
			}
			return button.Layout(gtx)
		})
	}))
	return components.Card(gtx, th, func(gtx components.C) components.D {
		return components.VerticalList(gtx, widgets...)
	})
}

func generateDeviceDetails(
	isConnected bool,
	manager *sdr.Manager,
	id devices.Id,
	th *themes.Theme,
	samplingMode *widget.Enum,
) []layout.FlexChild {
	var deviceDetails []layout.FlexChild
	if isConnected {
		if device, deviceErr := manager.Open(id); deviceErr != nil {
			log.WithFields(device.Fields()).WithError(deviceErr).Error("Open")
		} else if device != nil {
			switch d := device.(type) {
			case *rtlsdr.Connection:
				deviceDetails = rtlSDRCardBody(th, samplingMode, d)
				break
			default:
			}
		}
	}
	return deviceDetails
}

func handleConnectButton(
	connectButton *widget.Clickable,
	connectClicked bool,
	isConnected bool,
	manager *sdr.Manager,
	id devices.Id,
) {
	for connectButton.Clicked() {
		connectClicked = true
	}
	if connectClicked {
		if isConnected {
			manager.CloseAsync(id)
		} else {
			manager.OpenAsync(id)
		}
	}
}
