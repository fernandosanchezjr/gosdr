package pages

import (
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/themes"
	"github.com/fernandosanchezjr/gosdr/devices/sdr"
	"log"
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/x/component"
	"github.com/fernandosanchezjr/gosdr/cmd/gosdr/icon"
)

type Page interface {
	Actions() []component.AppBarAction
	Overflow() []component.OverflowAction
	Layout(gtx layout.Context, th *themes.Theme) layout.Dimensions
	NavItem() component.NavItem
}

type Router struct {
	pages   map[interface{}]Page
	current interface{}
	*component.ModalNavDrawer
	NavAnim component.VisibilityAnimation
	*component.AppBar
	*component.ModalLayer
	SDRManager *sdr.Manager
	State      *State
}

func NewRouter(th *themes.Theme, sdrManager *sdr.Manager) Router {
	modal := component.NewModal()
	nav := component.NewNav("GOSDR", "v0.0.1")
	modalNav := component.ModalNavFrom(&nav, modal)
	bar := component.NewAppBar(modal)
	bar.NavigationIcon = icon.MenuIcon
	na := component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 250,
	}

	return Router{
		pages:          make(map[interface{}]Page),
		ModalLayer:     modal,
		ModalNavDrawer: modalNav,
		AppBar:         bar,
		NavAnim:        na,
		SDRManager:     sdrManager,
		State:          NewState(th, sdrManager),
	}
}

func (r *Router) Register(tag interface{}, p Page) {
	r.pages[tag] = p
	navItem := p.NavItem()
	navItem.Tag = tag
	if r.current == interface{}(nil) {
		r.current = tag
		r.AppBar.Title = navItem.Name
		r.AppBar.SetActions(p.Actions(), p.Overflow())
	}
	r.ModalNavDrawer.AddNavItem(navItem)
}

func (r *Router) SwitchTo(tag interface{}) {
	p, ok := r.pages[tag]
	if !ok {
		return
	}
	navItem := p.NavItem()
	r.current = tag
	r.AppBar.Title = navItem.Name
	r.AppBar.SetActions(p.Actions(), p.Overflow())
}

func (r *Router) layoutContent(th *themes.Theme) layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X /= 3
				return r.NavDrawer.Layout(gtx, th.Theme, &r.NavAnim)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return r.pages[r.current].Layout(gtx, th)
			}),
		)
	})
}

func (r *Router) layoutBar(th *themes.Theme) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return r.AppBar.Layout(gtx, th.Theme, "Menu", "Actions")
	})
}

func (r *Router) layoutPage(gtx layout.Context, th *themes.Theme) layout.Dimensions {
	paint.Fill(gtx.Ops, th.Background.Dark.Bg)
	content := r.layoutContent(th)
	bar := r.layoutBar(th)
	flex := layout.Flex{Axis: layout.Vertical}
	flex.Layout(gtx, bar, content)
	r.ModalLayer.Layout(gtx, th.Theme)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

func (r *Router) Layout(gtx layout.Context, th *themes.Theme) layout.Dimensions {
	for _, event := range r.AppBar.Events(gtx) {
		switch event := event.(type) {
		case component.AppBarNavigationClicked:
			r.ModalNavDrawer.Appear(gtx.Now)
			r.NavAnim.Disappear(gtx.Now)
		case component.AppBarContextMenuDismissed:
			log.Printf("Context menu dismissed: %v", event)
		case component.AppBarOverflowActionClicked:
			log.Printf("Overflow action selected: %v", event)
		}
	}
	if r.ModalNavDrawer.NavDestinationChanged() {
		r.SwitchTo(r.ModalNavDrawer.CurrentNavDestination())
	}
	return r.layoutPage(gtx, th)
}
