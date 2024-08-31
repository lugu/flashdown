package main

import (
	"image/color"
	"sort"

	"fyne.io/fyne/v2/canvas"
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

// Show a white screen until the games are loaded, then shows HomeScreen.
func (s *SplashScreen) Show(app Application) {
	// Display white content
	rect := canvas.NewRectangle(color.White)
	app.Window().SetContent(rect)
	// load the games in the background
	go s.load(app)
}

func (s *SplashScreen) Hide(app Application) {}

func NewSplashScreen() Screen {
	return &SplashScreen{}
}
