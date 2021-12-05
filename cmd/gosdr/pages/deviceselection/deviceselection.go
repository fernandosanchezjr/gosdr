package deviceselection

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	alo "github.com/fernandosanchezjr/gosdr/cmd/gosdr/applayout"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/icon"
	page "github.com/fernandosanchezjr/gosdr/cmd/gosdr/pages"
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
	return &Page{
		Router: router,
	}
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
		Name: "Device Selection",
		Icon: icon.SettingsIcon,
	}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(
			gtx,
			layout.Rigid(func(gtx C) D {
				return alo.DetailRow{}.Layout(
					gtx,
					material.Body1(th, "Found devices").Layout,
					func(gtx C) D {
						return p.deviceList.Layout(gtx, len(p.Devices), func(gtx C, i int) D {
							return material.Body2(th, p.Devices[i].String()).Layout(gtx)
						})
					},
				)
			}),
		)
	})
}
