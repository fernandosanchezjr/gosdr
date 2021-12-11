package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/x/component"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
)

func Card(gtx C, th *themes.Theme, w layout.Widget) D {
	return component.Surface(th.Theme).Layout(gtx, func(gtx C) D {
		return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx C) D {
			return w(gtx)
		})
	})
}
