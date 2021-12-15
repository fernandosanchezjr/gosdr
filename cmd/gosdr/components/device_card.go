package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices"
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

func DeviceCard(gtx C, th *themes.Theme, device *devices.Info) D {
	var widgets = []layout.FlexChild{
		layout.Rigid(func(gtx C) D {
			return deviceTitle(gtx, th, device)
		}),
		layout.Rigid(func(gtx C) D {
			return deviceSubTitle(gtx, th, device)
		}),
	}
	return Card(gtx, th, func(gtx C) D {
		return VerticalList(gtx, widgets...)
	})
}
