package main

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

type SplashScreen struct{}

func (s *SplashScreen) load(app Application) {
	games, err := loadGames()
	if err != nil {
		app.Display(NewFatalScreen(err))
		return
	}
	app.Display(NewWelcomeScreen(games))
}

// Show a white screen until the games are loaded, then shows WelcomeScreen.
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
