package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lugu/flashdown"
)

type CongratsScreen struct {
	game *flashdown.Game
}

func NewCongratsScreen(game *flashdown.Game) Screen {
	return &CongratsScreen{game: game}
}

func (s *CongratsScreen) Show(app Application) {
	window := app.Window()
	topBar := newProgressTopBar(app, s.game)
	label := container.New(layout.NewCenterLayout(),
		widget.NewLabel("Congratulations!"))
	button := bottomButton("Press to continue", func() {
		app.Display(NewSplashScreen())
	})

	box := container.New(layout.NewBorderLayout(topBar, button, nil, nil),
		topBar, button, label)
	window.SetContent(box)
	window.Canvas().SetOnTypedKey(EscapeKeyHandler(app))
}

func (s *CongratsScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
