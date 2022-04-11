package sdr

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/components"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/icon"
	page "github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Page struct {
	widget.List
	*page.Router
	deviceList layout.List
}

// New constructs a Page with the provided router.
func New(router *page.Router) *Page {
	var p = &Page{
		Router: router,
	}
	p.List.Axis = layout.Vertical
	return p
}

var _ page.Page = &Page{}

func (p *Page) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *Page) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{
		Name: "SDR",
		Icon: icon.RadioIcon,
	}
}

func (p *Page) Layout(gtx C, th *themes.Theme) D {
	return components.ListInset(gtx, th.Inset, func(gtx C) D {
		return material.List(th.Theme, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
			return components.VerticalList(gtx, p.State.DeviceCards()...)
		})
	})
}
