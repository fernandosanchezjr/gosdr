package components

import "gioui.org/layout"

func VerticalList(gtx C, widgets ...layout.FlexChild) D {
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Start,
	}.Layout(gtx, widgets...)
}
