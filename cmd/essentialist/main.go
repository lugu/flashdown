package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	application := app.NewWithID("essentialist")
	application.Settings().SetTheme(getTheme())
	window := application.NewWindow("Essentialist")
	window.Resize(fyne.NewSize(640, 480))
	NewApplication(application, window).Display(NewSplashScreen())
	window.ShowAndRun()
}
