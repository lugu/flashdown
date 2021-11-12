package main

import (
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lugu/flashdown"
)

type WelcomeScreen struct {
	games []*flashdown.Game
}

func NewWelcomeScreen(games []*flashdown.Game) Screen {
	return &WelcomeScreen{games: games}
}

func (s *WelcomeScreen) Show(app Application) {
	window := app.Window()
	buttons := make([]fyne.CanvasObject, len(s.games))
	for i, g := range s.games {
		game := g
		name := path.Base(game.Name())
		label := fmt.Sprintf("%s (%.0f%%)", name, game.Success())
		button := widget.NewButton(label, func() {
			window.SetCloseIntercept(func() {
				game.Save()
				window.Close()
			})
			app.Display(NewQuestionScreen(game))
		})
		buttons[i] = button
	}
	if len(s.games) == 0 {
		info := fmt.Sprintf("No deck found in %s",
			getDirectory().String())
		label := widget.NewLabel(info)
		label.Wrapping = fyne.TextWrapBreak
		buttons = append(buttons, label)
	}

	vbox := container.New(layout.NewVBoxLayout(), buttons...)
	center := container.NewVScroll(vbox)
	topBar := newWelcomeTopBar(app)

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
}

func (s *WelcomeScreen) Hide(app Application) {}
