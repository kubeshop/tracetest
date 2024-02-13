package components

import "github.com/rivo/tview"

type RenderScheduler interface {
	Render(f func())
}

type appRenderScheduler struct {
	app *tview.Application
}

func (s *appRenderScheduler) Render(f func()) {
	s.app.QueueUpdateDraw(f)
}

func NewRenderScheduler(app *tview.Application) RenderScheduler {
	return &appRenderScheduler{app: app}
}
