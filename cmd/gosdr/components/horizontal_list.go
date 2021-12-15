package components

import "gioui.org/layout"

func HorizontalList(gtx C, weightSum float32, widgets ...layout.FlexChild) D {
	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		WeightSum: weightSum,
	}.Layout(gtx, widgets...)
}
