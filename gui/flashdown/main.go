package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	ui "github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/font"
	"github.com/aarzilli/nucular/style"
	"github.com/acarl005/stripansi"
	"github.com/lugu/flashdown"
)

var (
	DefaultDirectory = "/home/ludovicguegan/src/flashdown/samples"
	screen           ui.UpdateFn
	game             *flashdown.Game
)

const (
	TextScale       = 1.0
	TopBarHeight    = 40
	BottomBarHeight = 150
)

// newSelectDeckScreen searches for decks in directory and return a screen with
// a scroll bar to display them.
func newSelectDeckScreen() ui.UpdateFn {
	files, err := getFiles()
	if err != nil {
		return newErrorScreen(fmt.Errorf("Failed to list files: %s", err))
	}
	if len(files) == 0 {
		log.Fatal("no files")
	}
	// Create a list of games just to print the current progress.
	games := make([]*flashdown.Game, len(files))
	for i, f := range files {
		games[i], err = flashdown.NewGame(false, []string{f})
		if err != nil {
			return newErrorScreen(fmt.Errorf("Failed to open %s: %s", f, err))
		}
	}
	return func(w *ui.Window) {
		if len(games) == 0 {
			log.Fatal("no games")
		}
		w.Row(50).Dynamic(1)
		for i, g := range games {
			filename := files[i]
			label := fmt.Sprintf("%s (%.0f)", filename, g.Success())
			if w.ButtonText(label) {
				game, err = flashdown.NewGame(false, []string{
					filename,
				})
				if err != nil {
					screen = newErrorScreen(fmt.Errorf("Cannot start: %s", err))
				}
				screen = questionScreen
			}
		}
	}
}

func renderMarkdown(md string, width int) string {
	output := markdown.Render(md, width, 0)
	return stripansi.Strip(string(output))
}

func drawTopBar(w *ui.Window) {
	w.RowScaled(TopBarHeight).Dynamic(1)
	percent := game.Success()
	current, total := game.Progress()
	w.Label(fmt.Sprintf("Session: %d/%d — Success: %.0f", current, total, percent), "LT")
}

func drawQuestion(w *ui.Window, question string) {
	height := (w.LayoutAvailableHeight() - BottomBarHeight) / 2
	w.RowScaled(height).Dynamic(1)
	w.Label(renderMarkdown(question, 100), "CC")
}

func drawAnswer(w *ui.Window, answer string) {
	height := w.LayoutAvailableHeight() - BottomBarHeight
	w.RowScaled(height).Dynamic(1)
	w.Label(renderMarkdown(answer, 100), "CC")
}

func drawSeeAnswerButton(w *ui.Window) {
	w.RowScaled(BottomBarHeight).Dynamic(1)
	if w.ButtonText("See Answer") {
		screen = answerScreen
	}
}

func drawAnswerButtons(w *ui.Window) {
	w.RowScaled((BottomBarHeight-10)/3).Static(0, 0)

	if w.ButtonText("Total blackout") {
		game.Review(flashdown.TotalBlackout)
		screen = questionScreen
	}
	if w.ButtonText("Perfect recall") {
		game.Review(flashdown.PerfectRecall)
		screen = questionScreen
	}
	if w.ButtonText("Incorrect difficult") {
		game.Review(flashdown.IncorrectDifficult)
		screen = questionScreen
	}
	if w.ButtonText("Correct easy") {
		game.Review(flashdown.CorrectEasy)
		screen = questionScreen
	}
	if w.ButtonText("Incorrect easy") {
		game.Review(flashdown.IncorrectEasy)
		screen = questionScreen
	}
	if w.ButtonText("Correct difficult") {
		game.Review(flashdown.CorrectDifficult)
		screen = questionScreen
	}
}

// newQuestionScreen creates a screen which displays a question.
func questionScreen(w *ui.Window) {
	if game.IsFinished() {
		game.Save()
		screen = newSelectDeckScreen()
		return
	}
	drawTopBar(w)
	drawQuestion(w, game.Question())
	drawAnswer(w, "")
	drawSeeAnswerButton(w)
}

// newAnswerScreen creates a screen which display the response.
func answerScreen(w *ui.Window) {
	drawTopBar(w)
	drawQuestion(w, game.Question())
	drawAnswer(w, game.Answer())
	drawAnswerButtons(w)
}

// newErrorScreen returns a new screen which displays an error.
func newErrorScreen(err error) ui.UpdateFn {
	return func(w *ui.Window) {
		w.Row(0).Dynamic(1)
		w.Label(err.Error(), "CC")
	}
}

func updatefn(w *ui.Window) {
	screen(w)
}

func save() {
	if game != nil {
		game.Save()
	}
}

// getFiles returns the list of decks. Order of preferences:
// 1. Files from command line
// 2. Default directory
func getFiles() ([]string, error) {
	files := []string{}
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			files = append(files, os.Args[i])
		}
		return files, nil
	}

	fileInfos, err := ioutil.ReadDir(DefaultDirectory)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfos {
		if strings.HasSuffix(info.Name(), ".md") {
			filename := path.Join(DefaultDirectory, info.Name())
			files = append(files, filename)
		}
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("No deck found in: %s", DefaultDirectory)
	}
	return files, nil
}

func main() {
	screen = newSelectDeckScreen()
	wnd := ui.NewMasterWindow(0, "Flashdown", updatefn)
	style := style.FromTheme(style.DarkTheme, TextScale)
	var err error
	style.Font, err = font.NewFace(resourceNotoSansMonoRegularTtf.StaticContent, 18)
	if err != nil {
		log.Fatal(err)
	}
	style.Text.Padding = image.Point{0, 0}
	wnd.SetStyle(style)
	defer save()
	wnd.Main()
}
