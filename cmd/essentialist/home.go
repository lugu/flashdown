package main

import (
	"fmt"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	flashdown "github.com/lugu/flashdown/internal"
)

type HomeScreen struct {
	decks   []flashdown.DeckAccessor
	games   []*flashdown.Game
	cardsNb int
}

func NewHomeScreen(decks []flashdown.DeckAccessor) Screen {
	return &HomeScreen{
		decks:   decks,
		games:   make([]*flashdown.Game, len(decks)),
		cardsNb: getRepetitionLenght(),
	}
}

func (s *HomeScreen) keyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		if key.Name != "" {
			switch key.Name {
			case fyne.KeyQ, fyne.KeyEscape:
				app.Window().Close()
			case fyne.KeyReturn:
				s.startQuickSession(app)
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
				s.startQuickSession(app)
			case fyne.HardwareKey{ScanCode: 39}: // S
				app.Display(NewSettingsScreen())
			case fyne.HardwareKey{ScanCode: 43}: // H
				app.Display(NewHelpScreen())
			}
		}
	}
}

func (s *HomeScreen) startQuickSession(app Application) {
	game, err := flashdown.NewGameFromAccessors("all", s.cardsNb, s.decks...)
	if err != nil {
		app.Display(NewErrorScreen(err))
		return
	}
	app.Display(NewQuestionScreen(game))
}

func (s *HomeScreen) updateDeckButton(app Application, label *widget.Label, deck flashdown.DeckAccessor) *flashdown.Game {
	deckName := deck.DeckName()
	game, err := flashdown.NewGameFromAccessors(deckName, s.cardsNb, deck)
	if err != nil {
		label.SetText(fmt.Sprintf("Failed to load %s: %s", deckName, err))
		return nil
	}
	name := path.Base(game.Name())
	current, total := game.Progress()
	success := game.Success()
	content := fmt.Sprintf("%s (%.0f%% - %d/%d)", name, success, current, total)
	label.SetText(content)
	return game
}

func (s *HomeScreen) deckList(app Application) fyne.CanvasObject {
	if len(s.decks) == 0 {
		info := fmt.Sprintf("No deck found in %s", getDirectory().String())
		label := widget.NewLabel(info)
		label.Wrapping = fyne.TextWrapBreak
		return label
	}
	list := widget.NewList(
		func() int {
			return len(s.decks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			s.games[i] = s.updateDeckButton(app, o.(*widget.Label), s.decks[i])
		})
	list.OnSelected = func(id widget.ListItemID) {
		app.Display(NewQuestionScreen(s.games[id]))
	}
	return list
}

func (s *HomeScreen) Show(app Application) {
	topBar := newHomeTopBar(app, s)
	list := s.deckList(app)
	app.Window().SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, list))
	app.Window().Canvas().SetOnTypedKey(s.keyHandler(app))
}

func (s *HomeScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
