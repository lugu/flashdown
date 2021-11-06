package main

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"net/url"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
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

func getThemeName() string {
	prefs := fyne.CurrentApp().Preferences()
	return prefs.StringWithFallback("theme", "light")
}

func getTheme() fyne.Theme {
	prefs := fyne.CurrentApp().Preferences()
	dir := prefs.String("theme")
	switch dir {
	case "light":
		return theme.LightTheme()
	case "dark":
		return theme.DarkTheme()
	}
	return theme.LightTheme()
}

func setThemeName(name string) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("theme", name)
}

// getDirectory return the location where to look for decks.
func getDirectory() fyne.URI {
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

func setDirectory(dir fyne.URI) {
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

func newProgressTopBar(window fyne.Window, game *flashdown.Game) *fyne.Container {
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

func newWelcomeTopBar(window fyne.Window) *fyne.Container {
	label := widget.NewLabel("Welcome")
	back := widget.NewButton("Settings", func() {
		SettingsScreen(window)
	})
	return container.New(layout.NewBorderLayout(nil, nil, nil,
		back), back, label)
}

func newSettingsTopBar(window fyne.Window) *fyne.Container {
	label := widget.NewLabel("Settings")
	back := widget.NewButton("Home", func() {
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
	topBar := newProgressTopBar(window, game)
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

	topBar := newProgressTopBar(window, game)
	answer := continueButton(window, game)
	cards := newCards(game.Question(), "")

	vbox := container.New(layout.NewBorderLayout(topBar, answer, nil, nil),
		topBar, answer, cards)
	window.SetContent(vbox)
}

func CongratulationScreen(window fyne.Window, g *flashdown.Game) {
	topBar := newProgressTopBar(window, g)
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

func cleanUpStorage() error {
	files, err := storage.List(getDirectory())
	if err != nil {
		return err
	}
	for _, file := range files {
		// BUG: Fyne returns some empty entries.
		if file == nil {
			continue
		}
		err := storage.Delete(file)
		if err != nil {
			continue
		}
	}
	return nil
}

func importFile(source fyne.URI) error {
	err := storage.Copy(source, getDirectory())
	if err != nil {
		return err
	}
	return nil
}

func importDirectory(directory fyne.ListableURI) error {
	files, err := directory.List()
	if err != nil {
		return err
	}
	for _, file := range files {
		// BUG: Fyne returns some empty entries.
		if file == nil {
			continue
		}
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

func cleanUpStorageButton(window fyne.Window) *widget.Button {
	cb := func(yes bool) {
		if !yes {
			return
		}
		err := cleanUpStorage()
		if err != nil {
			ErrorScreen(window,
				fmt.Errorf("Failed to erase data: %s", err))
		}
	}
	label := "Do you want to delete all imported data?"
	return widget.NewButton("Erase storage", func() {
		dialog.ShowConfirm("Erase storage", label, cb, window)
	})
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
	return widget.NewButton("Import Directory", func() {
		dialog.NewFolderOpen(importCallback, window).Show()
	})
}

func changeDirectoryButton(window fyne.Window) *widget.Button {
	changeDirectoryCallback := func(d fyne.ListableURI, err error) {
		if err != nil {
			ErrorScreen(window, err)
			return
		}
		if d == nil {
			return
		}
		setDirectory(d)
		WelcomeScreen(window)
	}
	return widget.NewButton("Change Directory", func() {
		dialog.NewFolderOpen(changeDirectoryCallback, window).Show()
	})
}

func switchThemeButton(window fyne.Window) *widget.Button {
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
		SettingsScreen(window)
	})
}

func dbFile(file fyne.URI) (fyne.URI, error) {
	return storage.ParseURI(file.String() + ".db")
}

func loadGames() ([]*flashdown.Game, error) {
	games := make([]*flashdown.Game, 0)

	log.Printf("!!! directory %s", getDirectory().String())
	files, err := storage.List(getDirectory())
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file == nil {
			continue
		}
		log.Printf("!!! file %s (%s)", file.String(), file.Extension())
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

func SettingsScreen(window fyne.Window) {
	topBar := newSettingsTopBar(window)

	buttons := make([]fyne.CanvasObject, 0)
	if fyne.CurrentDevice().IsMobile() {
		buttons = append(buttons, importDirectoryButton(window))
		buttons = append(buttons, cleanUpStorageButton(window))
	} else {
		buttons = append(buttons, changeDirectoryButton(window))
	}
	buttons = append(buttons, switchThemeButton(window))
	center := container.NewVScroll(container.New(layout.NewVBoxLayout(),
		buttons...))
	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
}

func WelcomeScreen(window fyne.Window) {
	topBar := newWelcomeTopBar(window)

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
	if len(games) == 0 {
		info := fmt.Sprintf("No deck found in %s",
			getDirectory().String())
		label := widget.NewLabel(info)
		label.Wrapping = fyne.TextWrapBreak
		buttons = append(buttons, label)
	}

	vbox := container.New(layout.NewVBoxLayout(), buttons...)
	center := container.NewVScroll(vbox)

	window.SetContent(container.New(layout.NewBorderLayout(
		topBar, nil, nil, nil), topBar, center))
}

func main() {
	application := app.NewWithID("flashdown")
	application.Settings().SetTheme(getTheme())
	window := application.NewWindow("Flashdown")
	window.Resize(fyne.NewSize(640, 480))
	WelcomeScreen(window)
	window.ShowAndRun()
}
