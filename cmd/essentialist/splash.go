package main

import (
	"sort"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type SplashScreen struct{}

func (s *SplashScreen) load(app Application) {
	decks, err := loadDecks()
	if err != nil {
		app.Display(NewErrorScreen(err))
		return
	}
	sort.SliceStable(decks, func(i, j int) bool {
		return decks[i].DeckName() < decks[j].DeckName()
	})
	app.Display(NewHomeScreen(decks))
}

// Show an empty screen until the games are loaded, then shows HomeScreen.
func (s *SplashScreen) Show(app Application) {
	emptyContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer())
	app.Window().SetContent(emptyContainer)
	go s.load(app) // load the games in the background
}

func (s *SplashScreen) Hide(app Application) {}

func NewSplashScreen() Screen {
	return &SplashScreen{}
}
