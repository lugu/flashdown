package main

import (
	"fmt"
	"image/color"
	"log"
	"net/url"
	"path"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	flashdown "github.com/lugu/flashdown/internal"
)

const (
	cardsNbEntry   = "number of cards per session"
	directoryEntry = "directory"
	themeEntry     = "theme"
)

// overrideDirectory is used to specify a directory as an argument.
var overrideDirectory = ""

func getThemeName() string {
	prefs := fyne.CurrentApp().Preferences()
	return prefs.StringWithFallback("theme", "light")
}

func getTheme() fyne.Theme {
	prefs := fyne.CurrentApp().Preferences()
	dir := prefs.String(themeEntry)
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
	prefs.SetString(themeEntry, name)
}

func makeDefaultDirectory() (fyne.URI, error) {
	root := fyne.CurrentApp().Storage().RootURI()
	child, err := storage.Child(root, "Cards")
	if err != nil {
		return nil, err
	}
	b, err := storage.Exists(child)
	if err != nil {
		return nil, err
	}
	if !b {
		err := storage.CreateListable(child)
		if err != nil {
			return nil, err
		}
	}
	return child, nil
}

func getRepetitionLenght() int {
	a := fyne.CurrentApp()
	prefs := a.Preferences()
	return prefs.IntWithFallback(cardsNbEntry, 20)
}

func setRepetitionLenght(nbCards int) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetInt(cardsNbEntry, nbCards)
}

// getDirectory return the location where to look for decks. Set
// overrideDirectory to select which directory is returned by getDirectory. If
// overrideDirectory is unset, getDirectory returns the dirextory from the
// settings.
func getDirectory() fyne.URI {
	if overrideDirectory != "" {
		return storage.NewFileURI(overrideDirectory)
	}

	app := fyne.CurrentApp()
	prefs := app.Preferences()
	dir := prefs.StringWithFallback(directoryEntry, "")
	if dir != "" {
		uri, err := storage.ParseURI(dir)
		if err == nil {
			return uri
		}
		log.Printf("Failed to parse %s: %v", dir, err)
	}
	directory, err := makeDefaultDirectory()
	if err != nil {
		return app.Storage().RootURI()
	}
	setDirectory(directory)
	return directory
}

func setDirectory(dir fyne.URI) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString(directoryEntry, dir.String())
}

func newTopBar(leftText string, buttons ...fyne.CanvasObject) *fyne.Container {
	label := widget.NewLabel(leftText)
	objects := append([]fyne.CanvasObject{label, layout.NewSpacer()}, buttons...)
	return container.New(layout.NewHBoxLayout(), objects...)
}

func newHelpTopBar(app Application) *fyne.Container {
	home := widget.NewButton("Home", func() {
		app.Display(NewSplashScreen())
	})
	return newTopBar("Help", home)
}

func newErrorTopBar(app Application) *fyne.Container {
	home := widget.NewButton("Settings", func() {
		app.Display(NewSettingsScreen())
	})
	return newTopBar("Error", home)
}

func newHomeTopBar(app Application, s *HomeScreen) *fyne.Container {
	settings := widget.NewButton("Settings", func() {
		app.Display(NewSettingsScreen())
	})
	start := widget.NewButton("Start", func() {
		s.startQuickSession(app)
	})
	help := widget.NewButton("Help", func() {
		app.Display(NewHelpScreen())
	})
	return newTopBar("Home", start, help, settings)
}

func newProgressTopBar(app Application, game *flashdown.Game) *fyne.Container {
	percent := game.Success()
	current, total := game.Progress()
	text := fmt.Sprintf("Session: %d/%d — Success: %.0f%% — %s",
		current, total, percent, game.DeckName())
	home := widget.NewButton("Home", func() {
		game.Save()
		app.Display(NewSplashScreen())
	})
	return newTopBar(text, home)
}

// bottomButton return a large button.
func bottomButton(label string, cb func()) *fyne.Container {
	button := widget.NewButton(label, cb)

	// Construct a invisible rectanble to force the height of the button to
	// match those of the answer screen (3 rows).
	height := container.New(layout.NewGridLayout(1), button, button,
		button).Size().Height
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(0, height))

	return container.New(layout.NewBorderLayout(nil, nil, rect, nil),
		button)
}

func continueButton(app Application, game *flashdown.Game) *fyne.Container {
	return bottomButton("See answer", func() {
		app.Display(NewAnswerScreen(game))
	})
}

func forHuman(f fyne.URI) string {
	file, err := url.PathUnescape(f.String())
	if err != nil {
		return f.String()
	}
	return file
}

func cleanDirectory() error {
	dir := getDirectory()
	files, err := storage.List(dir)
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
			return fmt.Errorf("Failed to delete %s: %s",
				file.Name(), err)
		}
	}
	return nil
}

func importFile(source fyne.URI) error {
	decoded, err := url.PathUnescape(source.Name())
	filename := path.Base(decoded)

	destination, err := storage.Child(getDirectory(), filename)
	if err != nil {
		return fmt.Errorf("Failed to create %s at %s, %s",
			filename, getDirectory().String(), err)
	}

	err = storage.Copy(source, destination)
	if err == nil {
		return nil
	}
	err = repository.GenericCopy(source, destination)
	if err == nil {
		return nil
	}
	return fmt.Errorf("Copy error\nSource: %s\nDestination: %s\n%s",
		source.String(), destination.String(), err)
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

func dbFile(file fyne.URI) (fyne.URI, error) {
	uri := file.String()
	base := path.Base(uri)
	uri = strings.Replace(uri, base, "."+base+".db", 1)
	return storage.ParseURI(uri)
}

func loadDecks() ([]flashdown.DeckAccessor, error) {
	return loadDir(getDirectory())
}

func loadDir(dir fyne.URI) ([]flashdown.DeckAccessor, error) {
	files, err := storage.List(dir)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	results := make(chan flashdown.DeckAccessor, len(files))
	errors := make(chan error, len(files))

	for _, file := range files {
		go func(file fyne.URI) {
			defer wg.Done()
			if file == nil {
				return
			}
			if file.Extension() != ".md" && file.Name() != dir.Name() {
				println("%s", file.String())
				accessors, err := loadDir(file)
				if err != nil {
					return
				}
				for _, a := range accessors {
					results <- a
				}
				return
			}

			db, err := dbFile(file)
			if err != nil {
				errors <- fmt.Errorf("Failed to create URI: %s", err)
				return
			}
			results <- NewDeckAccessor(file, db)
		}(file)
	}

	var wg2 sync.WaitGroup
	wg2.Add(1)
	accessors := make([]flashdown.DeckAccessor, 0)
	go func() {
		defer wg2.Done()
		for accessor := range results {
			accessors = append(accessors, accessor)
		}
	}()
	wg.Wait()
	close(results)
	wg2.Wait()

	select {
	case err := <-errors:
		return nil, err
	default:
		return accessors, nil
	}
}

func EscapeKeyHandler(app Application) func(*fyne.KeyEvent) {
	return func(key *fyne.KeyEvent) {
		if key.Name != "" {
			switch key.Name {
			case fyne.KeyQ, fyne.KeyEscape:
				app.Display(NewSplashScreen())
			}
		} else {
			switch key.Physical {
			case fyne.HardwareKey{ScanCode: 9}, fyne.HardwareKey{ScanCode: 24}: // Escape
				app.Display(NewSplashScreen())
			}
		}
	}
}
