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
	window := app.Window()
	label := widget.NewLabel(e.err.Error())
	label.Wrapping = fyne.TextWrapBreak
	button := widget.NewButton("Return", func() {
		app.Display(NewSplashScreen())
	})
	vbox := container.New(layout.NewVBoxLayout(),
		label, layout.NewSpacer(), button)
	window.SetContent(vbox)
}

func (e *ErrorScreen) Hide(app Application) {}

type FatalScreen struct {
	err error
}

func NewFatalScreen(err error) Screen {
	return &FatalScreen{err: err}
}

func (e *FatalScreen) Show(app Application) {
	topBar := newHomeTopBar(app)
	errLabel := widget.NewLabel(e.err.Error())
	errLabel.Wrapping = fyne.TextWrapBreak
	vbox := container.New(layout.NewVBoxLayout(),
		topBar, layout.NewSpacer(), errLabel)
	app.Window().SetContent(vbox)
}

func (e *FatalScreen) Hide(app Application) {}
