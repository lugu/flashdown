package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lugu/flashdown"
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

func (s *SettingsScreen) selectRepetition(app Application) *widget.Select {
	selections := []string{
		"10 cards",
		"20 cards",
		"30 cards",
		"40 cards",
		"50 cards",
		"Remaining cards to learn",
		"All cards",
	}
	values := []int{
		10,
		20,
		30,
		40,
		50,
		flashdown.CARDS_TO_REVIEW,
		flashdown.ALL_CARDS,
	}
	onChange := func(selected string) {
		for i, s := range selections {
			if s == selected {
				setRepetitionLenght(values[i])
				return
			}
		}
	}
	repetitions := widget.NewSelect(selections, onChange)
	repetitions.Alignment = fyne.TextAlignCenter
	cardsNb := getRepetitionLenght()
	for i, v := range values {
		if v == cardsNb {
			repetitions.SetSelected(selections[i])
			break
		}
	}
	return repetitions
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

func (s *SettingsScreen) Show(app Application) {
	window := app.Window()
	topBar := s.newSettingsTopBar(app)

	objects := make([]fyne.CanvasObject, 0)
	if fyne.CurrentDevice().IsMobile() {
		objects = append(objects, s.importDirectoryButton(app))
		objects = append(objects, s.cleanUpStorageButton(app))
	} else {
		objects = append(objects, s.changeDirectoryButton(app))
	}
	objects = append(objects, s.switchThemeButton(app))
	objects = append(objects, s.selectRepetition(app))
	center := container.NewVScroll(container.New(layout.NewVBoxLayout(),
		objects...))
	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
	window.Canvas().SetOnTypedKey(EscapeKeyHandler(app))
}

func (s *SettingsScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
