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

func (e *ErrorScreen) Show(app Application) {
	topBar := newErrorTopBar(app)
	errLabel := widget.NewLabel(e.err.Error())
	errLabel.Wrapping = fyne.TextWrapBreak
	vbox := container.New(layout.NewVBoxLayout(),
		topBar, layout.NewSpacer(), errLabel)
	app.Window().SetContent(vbox)
	app.Window().Canvas().SetOnTypedKey(EscapeKeyHandler(app))
}

func (e *ErrorScreen) Hide(app Application) {}
