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

type HomeScreen struct {
	decks []flashdown.DeckAccessor
}

func NewHomeScreen(decks []flashdown.DeckAccessor) Screen {
	return &HomeScreen{decks: decks}
}

func (s *HomeScreen) keyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyQ, fyne.KeyEscape:
			app.Window().Close()
		}
	}
}

func (s *HomeScreen) Show(app Application) {
	window := app.Window()
	buttons := make([]fyne.CanvasObject, len(s.decks))
	cardsNb := getRepetitionLenght()
	for i, d := range s.decks {
		game, err := flashdown.NewGameFromAccessors(d.DeckName(), cardsNb, d)
		if err != nil {
			// TODO: Show a message and continue
			app.Display(NewFatalScreen(
				fmt.Errorf("Failed to load %s: %s",
					d.DeckName(), err)))
			return
		}
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
	if len(s.decks) == 0 {
		info := fmt.Sprintf("No deck found in %s",
			getDirectory().String())
		label := widget.NewLabel(info)
		label.Wrapping = fyne.TextWrapBreak
		buttons = append(buttons, label)
	}

	vbox := container.New(layout.NewVBoxLayout(), buttons...)
	center := container.NewVScroll(vbox)
	topBar := newHomeTopBar(app, s.decks)

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
	window.Canvas().SetOnTypedKey(s.keyHandler(app))
}

func (s *HomeScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
