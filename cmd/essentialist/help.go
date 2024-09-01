package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	helpContent = `
----
### Shortcuts during a session:

-   Escape or 'q' - show home menu
-   Space or Return - return the card to see the answer
-   '0' - Total blackout
-   '1' - Incorrect response, but upon seeing the answer it felt familiar
-   '2' - Incorrect response, but upon seeing the answer it seemed easy to remember
-   '3' - Correct response, with serious difficulty
-   '4' - Correct response, after some hesitation
-   '5' - Perfect response
-   's' or 'n' - skip the card and go to the next card
-   'p' - go to the previous card

----
### Home screen shortcuts:

-   Escape or 'q' - exit
-   Enter - start quick session
-   'h' - show help menu
-   's' - show settings menu
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
	window.Canvas().SetOnTypedKey(EscapeKeyHandler(app))
}

func (e *helpScreen) Hide(app Application) {}
