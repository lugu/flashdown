package main

import (
	"crypto/md5"
	"fmt"
	"image/color"
	"io"
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

func newSelectDirectory(dir fyne.URI, window fyne.Window) *fyne.Container {
	dirText := fmt.Sprintf("Directory: %s", forHuman(dir))
	dirLabel := widget.NewLabel(dirText)
	dirLabel.Wrapping = fyne.TextWrapBreak
	dirChange := widget.NewButton("Change directory", func() {
		dialog.NewFolderOpen(func(d fyne.ListableURI, err error) {
			if err != nil {
				ErrorScreen(window, err)
				return
			}
			if d != nil {
				if d.String() != dir.String() {
					setDirectoryURI(d)
				}
				WelcomeScreen(window)
			}
		}, window).Show()
	})
	return container.New(layout.NewBorderLayout(nil, nil, nil,
		dirChange), dirChange, dirLabel)
}

// getFiles returns the list of decks inside the given directory.
// getFiles returns an error if it does not any deck.
func getFiles(dir fyne.URI) ([]fyne.URI, error) {

	files := []fyne.URI{}
	childs, err := storage.List(dir)
	if err != nil {
		return nil, fmt.Errorf("List %s: %v", forHuman(dir), err)
	}
	filter := storage.NewExtensionFileFilter([]string{".md"})
	for _, child := range childs {
		if child == nil {
			continue
		}
		if filter.Matches(child) {
			files = append(files, child)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("No deck found in: %s", forHuman(dir))
	}
	return files, nil
}

func dbFile(file fyne.URI) (fyne.URI, error) {
	a := fyne.CurrentApp()

	md5 := md5.Sum([]byte(file.String()))
	name := fmt.Sprintf("%x\n", md5)

	db, err := a.Storage().Open(name)
	if err != nil {
		db, err := a.Storage().Create(name)
		if err != nil {
			return nil, fmt.Errorf("Failed to open db: %s", err)
		}
		defer db.Close()
		w, err := storage.Writer(db.URI())
		if err != nil {
			return nil, fmt.Errorf("Failed to open reader: %s", err)
		}
		defer w.Close()
		return db.URI(), nil
	}
	defer db.Close()
	return db.URI(), nil
}

func loadGames(dir fyne.URI) ([]*flashdown.Game, error) {
	files, err := getFiles(dir)
	if err != nil {
		return nil, err
	}

	// Create a list of games just to print the current progress.
	games := make([]*flashdown.Game, len(files))
	for i, f := range files {
		db, err := dbFile(f)
		if err != nil {
			return nil, err
		}
		acc := &uriDeckAccessor{
			deck: f,
			db:   db,
		}

		games[i], err = flashdown.NewGameFromAccessor(f.Name(), acc)
		if err != nil {
			return nil, err
		}
	}
	return games, nil
}

func WelcomeScreen(window fyne.Window) {
	directory := directoryURI()
	topBar := newSelectDirectory(directory, window)

	games, err := loadGames(directory)
	if err != nil {
		topBar := newSelectDirectory(directory, window)
		errLabel := widget.NewLabel(err.Error())
		errLabel.Wrapping = fyne.TextWrapBreak
		vbox := container.New(layout.NewVBoxLayout(),
			topBar, layout.NewSpacer(), errLabel)
		window.SetContent(vbox)
		return
	}

	gameList := widget.NewList(
		func() int {
			return len(games)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			game := games[i]
			name := path.Base(game.Name())
			label := fmt.Sprintf("%s (%.0f%%)", name,
				game.Success())
			o.(*widget.Button).SetText(label)
			o.(*widget.Button).OnTapped = func() {
				window.SetCloseIntercept(func() {
					game.Save()
					window.Close()
				})
				QuestionScreen(window, game)
			}
		})

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, gameList))
}

func main() {
	window := app.NewWithID("flashdown").NewWindow("Flashdown")
	window.Resize(fyne.NewSize(640, 480))
	WelcomeScreen(window)
	window.ShowAndRun()
}
