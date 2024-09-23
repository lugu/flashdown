package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type ErrorScreen struct {
	err error
}

func NewErrorScreen(err error) Screen {
	return &ErrorScreen{err: err}
}

func ErrorKeyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		if key.Name != "" {
			switch key.Name {
			case fyne.KeyS:
				app.Display(NewSettingsScreen())
			case fyne.KeyQ, fyne.KeyEscape:
				app.Window().Close()
			}
		} else {
			switch key.Physical {
			case fyne.HardwareKey{ScanCode: 39}: // S
				app.Display(NewSettingsScreen())
			case fyne.HardwareKey{ScanCode: 9}, fyne.HardwareKey{ScanCode: 24}: // Escape
				app.Window().Close()
			}
		}
	}
}

func (e *ErrorScreen) Show(app Application) {
	topBar := newErrorTopBar(app)
	errLabel := widget.NewLabel(e.err.Error())
	errLabel.Wrapping = fyne.TextWrapBreak
	vbox := container.New(layout.NewVBoxLayout(),
		topBar, layout.NewSpacer(), errLabel)
	app.Window().SetContent(vbox)
	app.Window().Canvas().SetOnTypedKey(ErrorKeyHandler(app))
}

func (e *ErrorScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
