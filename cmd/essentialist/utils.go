package main

import (
	"fmt"
	"image/color"
	"log"
	"net/url"
	"path"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/lugu/flashdown"
)

func getThemeName() string {
	prefs := fyne.CurrentApp().Preferences()
	return prefs.StringWithFallback("theme", "light")
}

func getTheme() fyne.Theme {
	prefs := fyne.CurrentApp().Preferences()
	dir := prefs.String("theme")
	switch dir {
	case "light":
		return NewTheme(theme.LightTheme())
	case "dark":
		return NewTheme(theme.DarkTheme())
	}
	return theme.LightTheme()
}

func setThemeName(name string) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("theme", name)
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

// getDirectory return the location where to look for decks.
func getDirectory() fyne.URI {
	a := fyne.CurrentApp()
	prefs := a.Preferences()
	dir := prefs.StringWithFallback("directory", "")
	if dir != "" {
		uri, err := storage.ParseURI(dir)
		if err == nil {
			return uri
		}
		log.Print("Failed to parse %s: %s", dir, err)
	}
	directory, err := makeDefaultDirectory()
	if err != nil {
		log.Printf("Failed to create %s", directory.String())
		return a.Storage().RootURI()
	}
	setDirectory(directory)
	return directory
}

func setDirectory(dir fyne.URI) {
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("directory", dir.String())
}

func newTopBar(text string, button *widget.Button) *fyne.Container {
	label := widget.NewLabel(text)
	return container.New(layout.NewHBoxLayout(), label, layout.NewSpacer(), button)

}
func newProgressTopBar(app Application, game *flashdown.Game) *fyne.Container {
	percent := game.Success()
	current, total := game.Progress()
	text := fmt.Sprintf("Session: %d/%d — Success: %.0f%%",
		current, total, percent)
	back := widget.NewButton("Home", func() {
		game.Save()
		app.Display(NewSplashScreen())
	})
	return newTopBar(text, back)
}

func newHomeTopBar(app Application) *fyne.Container {
	back := widget.NewButton("Settings", func() {
		app.Display(NewSettingsScreen())
	})
	return newTopBar("Home", back)
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
		rect, button)
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
			return fmt.Errorf("Failed to delete %s: %s",
				file.Name(), err)
		}
	}
	return nil
}

func importFile(source fyne.URI) error {

	decoded, err := url.PathUnescape(source.Name())
	filename := path.Base(decoded)

	directory := getDirectory()
	destination, err := storage.Child(directory, filename)
	if err != nil {
		return fmt.Errorf("Failed to create %s at %s, %s",
			filename, directory.String(), err)
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

func loadGames() ([]*flashdown.Game, error) {
	games := make([]*flashdown.Game, 0)

	files, err := storage.List(getDirectory())
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file == nil {
			continue
		}
		if file.Extension() != ".md" {
			continue
		}

		db, err := dbFile(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to create URI: %s", err)
		}

		game, err := flashdown.NewGameFromAccessor(file.Name(),
			NewDeckAccessor(file, db))
		if err != nil {
			return nil, fmt.Errorf("Failed to load %s: %s",
				forHuman(file), err)
		}
		games = append(games, game)
	}

	return games, nil
}
