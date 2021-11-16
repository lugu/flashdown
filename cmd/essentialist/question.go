package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lugu/flashdown"
)

type maxWidthCenterLayout struct {
	width float32
}

// NewCenterWithMaxWidthLayout returns a layout which centers its content
// elements given them no more than width.
func NewMaxWidthCenterLayout(width float32) fyne.Layout {
	return &maxWidthCenterLayout{width}
}

func (c *maxWidthCenterLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	if size.Width > c.width {
		pos.X = (size.Width - c.width) / 2
		size.Width = c.width
	}
	for _, child := range objects {
		child.Resize(size)
		child.Move(pos)
	}
}

func (c *maxWidthCenterLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}
	return minSize
}

func space() fyne.CanvasObject {
	return layout.NewSpacer()
}

// TODO: make the test selectable
func card(md string) fyne.CanvasObject {
	richText := widget.NewRichTextFromMarkdown(md)
	width := richText.MinSize().Width
	richText.Wrapping = fyne.TextWrapWord
	return container.New(NewMaxWidthCenterLayout(width), richText)
}

type QuestionScreen struct {
	game *flashdown.Game
}

func NewQuestionScreen(game *flashdown.Game) Screen {
	return &QuestionScreen{game: game}
}

func (s *QuestionScreen) questionKeyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeySpace, fyne.KeyEnter:
			app.Display(NewAnswerScreen(s.game))
		}
	}
}

func (s *QuestionScreen) Show(app Application) {
	window := app.Window()

	topBar := newProgressTopBar(app, s.game)
	question := card("### " + s.game.Question())
	button := continueButton(app, s.game)

	vbox := container.New(layout.NewVBoxLayout(), topBar, space(), question,
		space(), button)
	window.SetContent(vbox)
	window.Canvas().SetOnTypedKey(s.questionKeyHandler(app))
}

func (s *QuestionScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}

type AnswerScreen struct {
	game *flashdown.Game
}

func NewAnswerScreen(game *flashdown.Game) Screen {
	return &AnswerScreen{game: game}
}

func (s *AnswerScreen) answersButton(app Application) *fyne.Container {
	bt := func(label string, score flashdown.Score) *widget.Button {
		return widget.NewButton(label,
			func() {
				s.reviewScore(app, score)
			})
	}
	buttons := []fyne.CanvasObject{
		bt("Total blackout", flashdown.TotalBlackout),
		bt("Perfect recall", flashdown.PerfectRecall),
		bt("Incorrect difficult", flashdown.IncorrectDifficult),
		bt("Correct difficult", flashdown.CorrectDifficult),
		bt("Incorrect easy", flashdown.IncorrectEasy),
		bt("Correct easy", flashdown.CorrectEasy),
	}
	return container.New(layout.NewGridLayout(2), buttons...)
}

func (s *AnswerScreen) reviewScore(app Application, score flashdown.Score) {
	s.game.Review(score)
	if s.game.IsFinished() {
		s.game.Save()
		app.Display(NewCongratsScreen(s.game))
	} else {
		app.Display(NewQuestionScreen(s.game))
	}
}

func (s *AnswerScreen) answerKeyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.Key0:
			s.reviewScore(app, flashdown.TotalBlackout)
		case fyne.Key1:
			s.reviewScore(app, flashdown.IncorrectDifficult)
		case fyne.Key2:
			s.reviewScore(app, flashdown.IncorrectEasy)
		case fyne.Key3:
			s.reviewScore(app, flashdown.CorrectDifficult)
		case fyne.Key4:
			s.reviewScore(app, flashdown.CorrectEasy)
		case fyne.Key5:
			s.reviewScore(app, flashdown.PerfectRecall)
		}
	}
}

func (s *AnswerScreen) Show(app Application) {
	window := app.Window()
	topBar := newProgressTopBar(app, s.game)

	question := card("### " + s.game.Question())
	line := canvas.NewLine(color.Gray16{0xaaaa})
	answer := card(s.game.Answer())

	buttons := s.answersButton(app)
	vbox := container.New(layout.NewVBoxLayout(), topBar, space(), question,
		space(), line, space(), answer, space(), buttons)
	window.SetContent(vbox)
	window.Canvas().SetOnTypedKey(s.answerKeyHandler(app))
}

func (s *AnswerScreen) Hide(app Application) {
	app.Window().Canvas().SetOnTypedKey(nil)
}
