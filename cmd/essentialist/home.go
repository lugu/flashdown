package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	flashdown "github.com/lugu/flashdown/internal"
)

// Warning: decks are loaded on demand by the list widget
type HomeScreen struct {
	accessors []flashdown.DeckAccessor
	decks     []*flashdown.Deck
	cardsNb   int
}

func NewHomeScreen(accessors []flashdown.DeckAccessor) Screen {
	return &HomeScreen{
		accessors: accessors,
		decks:     make([]*flashdown.Deck, len(accessors)),
		cardsNb:   getRepetitionLenght(),
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
	// We need to load remaning decks not yet loaded by by list widget.
	decks := make([]*flashdown.Deck, 0, len(s.accessors))
	for i, a := range s.accessors {
		if s.decks[i] == nil {
			deck, err := flashdown.NewDeck(a)
			if err != nil {
				s.decks[i] = deck
				decks = append(decks, deck)
			}
		} else {
			decks = append(decks, s.decks[i])
		}
	}
	game := flashdown.NewGame(s.cardsNb, decks...)
	if game.IsFinished() {
		app.Display(NewCongratsScreen(game))
	} else {
		app.Display(NewQuestionScreen(game))
	}
}

func (s *HomeScreen) updateDeckButton(app Application, label *widget.Label, i int) {
	deck := s.decks[i]
	toReview, total := deck.Stats()
	success := 100.0
	if total != 0 {
		success = 100 * ((float64(total) - float64(toReview)) / float64(total))
	}
	content := fmt.Sprintf("%s (%.0f%% of %d)", deck.Name, success, total)
	label.SetText(content)
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
			return widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if s.decks[i] == nil { // lazy loading
				var err error
				s.decks[i], err = flashdown.NewDeck(s.accessors[i])
				if err != nil {
					s.decks[i] = flashdown.NewEmptyDeck(
						s.accessors[i].DeckName())
					label.SetText(fmt.Sprintf("Failed to load %s: %s",
						s.accessors[i].DeckName(), err))
					return
				}
			}
			s.updateDeckButton(app, label, i)
		})
	list.OnSelected = func(id widget.ListItemID) {
		game := flashdown.NewGame(getRepetitionLenght(), s.decks[id])
		if game.IsFinished() {
			app.Display(NewCongratsScreen(game))
		} else {
			app.Display(NewQuestionScreen(game))
		}
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
