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
		if key.Name != "" {
			switch key.Name {
			case fyne.KeyQ, fyne.KeyEscape:
				app.Window().Close()
			case fyne.KeyReturn:
				s.StartQuickSession(app)
			case fyne.KeyH:
				app.Display(NewHelpScreen())
			case fyne.KeyS:
				app.Display(NewSettingsScreen())
			}
		} else {
			switch key.Physical {
			case fyne.HardwareKey{ScanCode: 9}, fyne.HardwareKey{ScanCode: 24}: // Escape
				app.Window().Close()
			case fyne.HardwareKey{ScanCode: 36}: // Enter
				s.StartQuickSession(app)
			case fyne.HardwareKey{ScanCode: 39}: // S
				app.Display(NewSettingsScreen())
			case fyne.HardwareKey{ScanCode: 43}: // H
				app.Display(NewHelpScreen())
			}
		}
	}
}

func (s *HomeScreen) StartQuickSession(app Application) {
	cardsNb := getRepetitionLenght()
	game, err := flashdown.NewGameFromAccessors("all", cardsNb, s.decks...)
	if err != nil {
		app.Display(NewErrorScreen(err))
		return
	}
	app.Display(NewQuestionScreen(game))
}

func (s *HomeScreen) Show(app Application) {
	window := app.Window()
	buttons := make([]fyne.CanvasObject, len(s.decks))
	cardsNb := getRepetitionLenght()
	if len(s.decks) == 0 {
		info := fmt.Sprintf("No deck found in %s",
			getDirectory().String())
		label := widget.NewLabel(info)
		label.Wrapping = fyne.TextWrapBreak
		buttons = append(buttons, label)
	}

	list := widget.NewList(
		func() int {
			return len(buttons)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			deck := s.decks[i]
			deckName := deck.DeckName()
			game, err := flashdown.NewGameFromAccessors(deckName, cardsNb, deck)
			if err != nil {
				// instead of unwrap, show line number
				o.(*widget.Button).SetText(fmt.Sprintf("Failed to load %s: %s", deckName, err))
				return
			}
			name := path.Base(game.Name())
			current, total := game.Progress()
			success := game.Success()
			label := fmt.Sprintf("%s (%.0f%% - %d/%d)", name, success, current, total)
			o.(*widget.Button).SetText(label)
			o.(*widget.Button).OnTapped = func() {
				window.SetCloseIntercept(func() {
					game.Save()
					window.Close()
				})
				app.Display(NewQuestionScreen(game))
			}
		})
	topBar := newHomeTopBar(app, s)
	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, list))
	window.Canvas().SetOnTypedKey(s.keyHandler(app))
}

func (s *HomeScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
