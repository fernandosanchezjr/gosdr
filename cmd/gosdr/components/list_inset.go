package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

func ListInset(gtx C, v unit.Value, w layout.Widget) D {
	return layout.Inset{
		Top:    v,
		Right:  unit.Dp(0),
		Bottom: v,
		Left:   v,
	}.Layout(gtx, w)
}
