package device_card

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/components"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices/rtlsdr"
	"strconv"
)

func rtlSDRCardBody(th *themes.Theme, enum *widget.Enum, device *rtlsdr.Connection) (contents []layout.FlexChild) {
	enum.Value = strconv.Itoa(int(device.Mode))
	contents = append(
		contents,
		layout.Rigid(
			func(gtx components.C) components.D {
				return samplingMode(gtx, th, enum)
			}),
	)
	return
}

func samplingMode(gtx components.C, th *themes.Theme, enum *widget.Enum) components.D {
	return components.HorizontalList(gtx, 4,
		layout.Rigid(material.Body1(th.Theme, "Direct Sampling Mode").Layout),
		layout.Rigid(material.RadioButton(th.Theme, enum, "0", "Off").Layout),
		layout.Rigid(material.RadioButton(th.Theme, enum, "1", "I Branch").Layout),
		layout.Rigid(material.RadioButton(th.Theme, enum, "2", "Q Branch").Layout),
	)
}
