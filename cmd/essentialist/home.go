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
	for i, d := range s.decks {
		game, err := flashdown.NewGameFromAccessors(d.DeckName(), cardsNb, d)
		if err != nil {
			// TODO: Show a message and continue
			app.Display(NewErrorScreen(
				fmt.Errorf("Failed to load %s: %s",
					d.DeckName(), err)))
			return
		}
		name := path.Base(game.Name())
		current, total := game.Progress()
		success := game.Success()
		label := fmt.Sprintf("%s (%.0f%% - %d/%d)", name, success, current, total)
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
	topBar := newHomeTopBar(app, s)

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
	window.Canvas().SetOnTypedKey(s.keyHandler(app))
}

func (s *HomeScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
