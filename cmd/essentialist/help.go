package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	helpContent = `Shortcuts:

-   Escape - quit
-   Space - return the card
-   5: Perfect response
-   4: Correct response, after some hesitation
-   3: Correct response, with serious difficulty
-   2: Incorrect response, but upon seeing the answer it seemed easy to remember
-   1: Incorrect response, but upon seeing the answer it felt familiar
-   0: Total blackout
`
)

type helpScreen struct{}

func NewHelpScreen() Screen {
	return &helpScreen{}
}

func helpMessage() *fyne.Container {
	richText := widget.NewRichTextFromMarkdown(helpContent)
	width := richText.MinSize().Width
	richText.Wrapping = fyne.TextWrapWord
	return container.New(NewMaxWidthCenterLayout(width), richText)
}

func (e *helpScreen) Show(app Application) {
	window := app.Window()
	topBar := newHelpTopBar(app)
	center := container.NewVScroll(helpMessage())
	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
}

func (e *helpScreen) Hide(app Application) {}
