package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SettingsScreen struct{}

func NewSettingsScreen() Screen {
	return &SettingsScreen{}
}

func (s *SettingsScreen) importDirectoryButton(app Application) *widget.Button {
	importCallback := func(d fyne.ListableURI, err error) {
		if err != nil {
			app.Display(NewErrorScreen(err))
			return
		}
		if d == nil {
			return
		}
		if err = importDirectory(d); err != nil {
			app.Display(NewErrorScreen(err))
			return
		}
		app.Display(NewSplashScreen())
	}
	return widget.NewButton("Import Directory", func() {
		dialog.NewFolderOpen(importCallback, app.Window()).Show()
	})
}

func (s *SettingsScreen) changeDirectoryButton(app Application) *widget.Button {
	changeDirectoryCallback := func(d fyne.ListableURI, err error) {
		if err != nil {
			app.Display(NewErrorScreen(err))
			return
		}
		if d == nil {
			return
		}
		setDirectory(d)
		app.Display(NewSplashScreen())
	}
	return widget.NewButton("Change Directory", func() {
		dialog.NewFolderOpen(changeDirectoryCallback, app.Window()).Show()
	})
}

func (s *SettingsScreen) switchThemeButton(app Application) *widget.Button {
	currentTheme := getThemeName()
	var newTheme string
	switch currentTheme {
	case "light":
		newTheme = "dark"
	case "dark":
		newTheme = "light"
	default:
		newTheme = "light"
	}
	buttonLabel := fmt.Sprintf("Theme: %s", newTheme)
	return widget.NewButton(buttonLabel, func() {
		setThemeName(newTheme)
		fyne.CurrentApp().Settings().SetTheme(getTheme())
		app.Display(NewSettingsScreen())
	})
}

func (s *SettingsScreen) cleanUpStorageButton(app Application) *widget.Button {
	cb := func(yes bool) {
		if !yes {
			return
		}
		err := cleanDirectory()
		if err != nil {
			app.Display(NewErrorScreen(err))
		}
	}
	label := fmt.Sprintf("Delete cards in %s/ ?", getDirectory().Name())
	return widget.NewButton("Erase storage", func() {
		dialog.ShowConfirm("Erase storage", label, cb, app.Window())
	})
}

func (s *SettingsScreen) newSettingsTopBar(app Application) *fyne.Container {
	back := widget.NewButton("Home", func() {
		app.Display(NewSplashScreen())
	})
	return newTopBar("Settings", back)
}

func (s *SettingsScreen) keyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyQ, fyne.KeyEscape:
			app.Display(NewSplashScreen())
		}
	}
}

func (s *SettingsScreen) Show(app Application) {
	window := app.Window()
	topBar := s.newSettingsTopBar(app)

	buttons := make([]fyne.CanvasObject, 0)
	if fyne.CurrentDevice().IsMobile() {
		buttons = append(buttons, s.importDirectoryButton(app))
		buttons = append(buttons, s.cleanUpStorageButton(app))
	} else {
		buttons = append(buttons, s.changeDirectoryButton(app))
	}
	buttons = append(buttons, s.switchThemeButton(app))
	center := container.NewVScroll(container.New(layout.NewVBoxLayout(),
		buttons...))
	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
	window.Canvas().SetOnTypedKey(s.keyHandler(app))
}

func (s *SettingsScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
