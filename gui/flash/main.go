package main

import (
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"net/url"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/lugu/flashdown"
)

// uriDeckAccessor must implement flashdown.DeckAccessor
var _ flashdown.DeckAccessor = (*uriDeckAccessor)(nil)

type uriDeckAccessor struct {
	deck fyne.URI
	db   fyne.URI
}

func (u *uriDeckAccessor) CardsReader() (io.ReadCloser, error) {
	return storage.Reader(u.deck)
}

func (u *uriDeckAccessor) MetaReader() (io.ReadCloser, error) {
	r, err := storage.Reader(u.db)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (u *uriDeckAccessor) MetaWriter() (io.WriteCloser, error) {
	w, err := storage.Writer(u.db)
	if err != nil {
		return nil, err
	}
	return w, err
}

// directoryURI return the location where to look for decks.
func directoryURI() fyne.URI {
	a := fyne.CurrentApp()
	prefs := a.Preferences()
	dir := prefs.StringWithFallback("directory",
		a.Storage().RootURI().String())
	uri, err := storage.ParseURI(dir)
	if err != nil {
		return a.Storage().RootURI()
	}
	return uri
}

func setDirectoryURI(dir fyne.URI) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("directory", dir.String())
}

func ErrorScreen(window fyne.Window, err error) {
	label := widget.NewLabel(err.Error())
	label.Wrapping = fyne.TextWrapBreak
	button := widget.NewButton("Return", func() {
		WelcomeScreen(window)
	})
	vbox := container.New(layout.NewVBoxLayout(),
		label, layout.NewSpacer(), button)
	window.SetContent(vbox)
}

func newTopBar(window fyne.Window, game *flashdown.Game) *fyne.Container {
	percent := game.Success()
	current, total := game.Progress()
	text := fmt.Sprintf("Session: %d/%d — Success: %.0f%%",
		current, total, percent)
	label := widget.NewLabel(text)
	back := widget.NewButton("Home", func() {
		game.Save()
		WelcomeScreen(window)
	})

	return container.New(layout.NewBorderLayout(nil, nil, nil,
		back), back, label)
}

// TODO: make the test selectable
func card(md string) *fyne.Container {
	o := NewMarkdownContainer(md)
	return container.New(layout.NewCenterLayout(), o)
}

func newCards(question, answer string) *fyne.Container {
	questionBox := card(question)
	answerBox := card(answer)

	var border float32 = 2.0
	t := canvas.NewRectangle(color.White)
	t.SetMinSize(fyne.NewSize(0, border))
	b := canvas.NewRectangle(color.White)
	b.SetMinSize(fyne.NewSize(0, border))

	questionCard := container.New(layout.NewBorderLayout(t, nil, nil, nil),
		t, questionBox)
	answerCard := container.New(layout.NewBorderLayout(t, b, nil, nil),
		t, b, answerBox)

	return container.New(layout.NewGridLayout(1), questionCard, answerCard)
}

// bottomButton return a large button.
func bottomButton(label string, cb func()) *fyne.Container {
	button := widget.NewButton("See Answer", cb)

	// Construct a invisible rectanble to force the height of the button to
	// be the same as answersButton (3 rows).
	height := container.New(layout.NewGridLayout(1), button, button,
		button).Size().Height
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(0, height))

	return container.New(layout.NewBorderLayout(nil, nil, rect, nil),
		rect, button)
}

func continueButton(window fyne.Window, game *flashdown.Game) *fyne.Container {
	return bottomButton("See Answer", func() {
		AnswerScreen(window, game)
	})
}

func answersButton(window fyne.Window, game *flashdown.Game) *fyne.Container {
	bt := func(label string, s flashdown.Score) *widget.Button {
		return widget.NewButton(label,
			func() {
				game.Review(s)
				QuestionScreen(window, game)
			})
	}
	buttons := []fyne.CanvasObject{
		bt("Total blackout", flashdown.TotalBlackout),
		bt("Correct difficult", flashdown.CorrectDifficult),
		bt("Incorrect difficult", flashdown.IncorrectDifficult),
		bt("Correct easy", flashdown.CorrectEasy),
		bt("Incorrect easy", flashdown.IncorrectEasy),
		bt("Perfect recall", flashdown.PerfectRecall),
	}
	return container.New(layout.NewGridLayout(2), buttons...)
}

func AnswerScreen(window fyne.Window, game *flashdown.Game) {
	topBar := newTopBar(window, game)
	cards := newCards(game.Question(), game.Answer())
	answers := answersButton(window, game)

	vbox := container.New(layout.NewBorderLayout(topBar, answers, nil, nil),
		topBar, answers, cards)

	window.SetContent(vbox)
}

func QuestionScreen(window fyne.Window, game *flashdown.Game) {
	if game.IsFinished() {
		game.Save()
		CongratulationScreen(window, game)
		return
	}

	topBar := newTopBar(window, game)
	answer := continueButton(window, game)
	cards := newCards(game.Question(), "")

	vbox := container.New(layout.NewBorderLayout(topBar, answer, nil, nil),
		topBar, answer, cards)
	window.SetContent(vbox)
}

func CongratulationScreen(window fyne.Window, g *flashdown.Game) {
	topBar := newTopBar(window, g)
	label := container.New(layout.NewCenterLayout(),
		widget.NewLabel("Congratulations!"))
	button := bottomButton("Press to continue", func() {
		WelcomeScreen(window)
	})

	box := container.New(layout.NewBorderLayout(topBar, button, nil, nil),
		topBar, button, label)
	window.SetContent(box)
}

func forHuman(f fyne.URI) string {
	file, err := url.PathUnescape(f.String())
	if err != nil {
		return f.String()
	}
	return file
}

func cleanUpStorage() {
	a := fyne.CurrentApp()
	root := a.Storage()
	names := root.List()
	for _, name := range names {
		child, err := root.Open(name)
		if err != nil {
			continue
		}
		defer child.Close()
		file := child.URI()
		if file.Extension() == ".md" {
			root.Remove(name)
		}
	}
}

func importFile(source fyne.URI) error {
	reader, err := storage.Reader(source)
	if err != nil {
		return nil
	}
	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("Failed to read file: %s, %s",
			forHuman(source), err)
	}

	a := fyne.CurrentApp()
	root := a.Storage()

	decoded, err := url.PathUnescape(source.Name())
	filename := path.Base(decoded)

	writer, err := root.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create file: %s, %s",
			source.Name(), err)
	}
	defer writer.Close()

	n, err := writer.Write(bytes)
	if n != len(bytes) {
		return fmt.Errorf("Partial copy: %d/%d bytes", n, len(bytes))
	}
	return nil
}

func importDirectory(directory fyne.ListableURI) error {
	files, err := directory.List()
	if err != nil {
		return err
	}
	cleanUpStorage()
	for _, file := range files {
		if file.Extension() != ".md" {
			continue
		}
		err = importFile(file)
		if err != nil {
			return fmt.Errorf("Cannot import %s: %s", file.String(), err)
		}
	}
	return nil
}

func importDirectoryButton(window fyne.Window) *widget.Button {
	importCallback := func(d fyne.ListableURI, err error) {
		if err != nil {
			ErrorScreen(window, err)
			return
		}
		if d == nil {
			return
		}
		if err = importDirectory(d); err != nil {
			ErrorScreen(window, err)
			return
		}
		WelcomeScreen(window)
	}
	button := widget.NewButton("Import Directory", func() {
		dialog.NewFolderOpen(importCallback, window).Show()
	})
	return button
}

func dbFile(file fyne.URI) (fyne.URI, error) {
	return storage.ParseURI(file.String() + ".db")
}

func loadGames() ([]*flashdown.Game, error) {
	games := make([]*flashdown.Game, 0)

	a := fyne.CurrentApp()
	root := a.Storage()

	for _, name := range root.List() {
		child, err := root.Open(name)
		if err != nil {
			continue
		}
		defer child.Close()
		file := child.URI()
		if file.Extension() != ".md" {
			continue
		}

		db, err := dbFile(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to create URI: %s", err)
		}

		game, err := flashdown.NewGameFromAccessor(file.Name(),
			&uriDeckAccessor{
				deck: file,
				db:   db,
			})
		if err != nil {
			return nil, fmt.Errorf("Failed to load %s: %s",
				forHuman(file), err)
		}
		games = append(games, game)
	}

	return games, nil
}

func WelcomeScreen(window fyne.Window) {
	topBar := importDirectoryButton(window)

	games, err := loadGames()
	if err != nil {
		errLabel := widget.NewLabel(err.Error())
		errLabel.Wrapping = fyne.TextWrapBreak
		vbox := container.New(layout.NewVBoxLayout(),
			topBar, layout.NewSpacer(), errLabel)
		window.SetContent(vbox)
		return
	}

	buttons := make([]fyne.CanvasObject, len(games))
	for i, g := range games {
		game := g
		name := path.Base(game.Name())
		label := fmt.Sprintf("%s (%.0f%%)", name, game.Success())
		button := widget.NewButton(label, func() {
			window.SetCloseIntercept(func() {
				game.Save()
				window.Close()
			})
			QuestionScreen(window, game)
		})
		buttons[i] = button
	}
	vbox := container.New(layout.NewVBoxLayout(), buttons...)
	center := container.NewVScroll(vbox)

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
}

func main() {
	window := app.NewWithID("flashdown").NewWindow("Flashdown")
	window.Resize(fyne.NewSize(640, 480))
	WelcomeScreen(window)
	window.ShowAndRun()
}
