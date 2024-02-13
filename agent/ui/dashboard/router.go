package dashboard

import (
	"github.com/rivo/tview"
)

type Router struct {
	*tview.Pages
}

func NewRouter() *Router {
	return &Router{
		Pages: tview.NewPages(),
	}
}

func (r *Router) AddPage(name string, page tview.Primitive) {
	r.Pages.AddPage(name, page, true, false)
}

func (r *Router) AddAndSwitchToPage(name string, page tview.Primitive) {
	r.Pages.AddAndSwitchToPage(name, page, true)
}
