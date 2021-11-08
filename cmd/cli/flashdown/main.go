package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/lugu/flashdown"
)

const (
	helpQuestion = `Press space to continue, 's' to skip or 'q' to quit`

	usageMsg = `Spaced repetition program for flashcards in markdown.

Usage: %s [-a] <deck_file> [<deck_file> ...]
Flags:
	-a : force all cards in the deck to be used.
	-h : show this message.

Deck files are plain text with heading level 1 being the question:

    # Question 1
    Answer 1
    # Question 2
    Answer 2
    [...]
`
)

var (
	helpAnswers = []string{
		` Press [0-5] to continue, 's' to skip, 'q' to quit, 'h' for help`,
		` Press [0-5] to continue, 's' to skip, 'q' to quit, 'h' for help

5: Perfect response
4: Correct response, after some hesitation
3: Correct response, with serious difficulty
2: Incorrect response, but upon seeing the answer it seemed easy to remember
1: Incorrect response, but upon seeing the answer it felt familiar
0: Total blackout`,
	}

	helpIndex  = 0
	helpAnswer = helpAnswers[helpIndex]

	game *flashdown.Game
)

func main() {
	// file, err := ioutil.TempFile(".", "log")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(file)

	if len(os.Args) < 2 {
		fmt.Printf(usageMsg, os.Args[0])
		os.Exit(1)
	}

	userAllCards := false
	files := make([]string, 0, len(os.Args))

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-a" {
			userAllCards = true
			continue
		}
		files = append(files, os.Args[i])
	}

	game, err := flashdown.NewGameFromFiles(userAllCards, files)
	if err != nil {
		log.Fatal(err)
	}
	if game.IsFinished() {
		return
	}
	defer game.Save()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	ui.Theme.Default = ui.StyleClear
	ui.Theme.Block.Title = ui.StyleClear
	ui.Theme.Block.Border = ui.StyleClear
	ui.Theme.Paragraph.Text = ui.StyleClear

	help := widgets.NewParagraph()
	help.Border = false
	help.PaddingLeft = -1
	help.PaddingRight = -1
	help.PaddingTop = -1
	help.PaddingBottom = -1
	help.Border = false

	q := NewMarkdownArea()
	a := NewMarkdownArea()

	grid := ui.NewGrid()
	grid.Set(
		ui.NewRow(1.0/2,
			q,
		),
		ui.NewRow(1.0/2,
			a,
		),
	)

	updateTitle := func() {
		percent := game.Success()
		current, total := game.Progress()
		q.Title = fmt.Sprintf(`Card: %d/%d — Success %2.0f%%`,
			current, total, percent)
	}

	ask := func() {
		updateTitle()
		q.Text = game.Question()
		a.Text = ""
		help.Text = helpQuestion
		ui.Clear()
		ui.Render(grid, help)
	}
	review := func(score flashdown.Score) {
		game.Review(score)
		ask()
	}
	answer := func() {
		q.Text = game.Question()
		a.Text = game.Answer()
		help.Text = helpAnswer
		ui.Clear()
		ui.Render(grid, help)
	}
	ask()

	resize := func(width, height int) {
		help.Text = helpAnswer
		helpHeigh := strings.Count(helpAnswer, "\n") + 1
		grid.SetRect(0, 0, width, height-helpHeigh)
		help.SetRect(0, height-helpHeigh, width, height)
		ui.Clear()
		ui.Render(grid, help)
	}

	termWidth, termHeight := ui.TerminalDimensions()
	resize(termWidth, termHeight)

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "s":
				game.Skip()
				if game.IsFinished() {
					return
				}
				ask()
			case "q", "<C-c>":
				return
			case "h":
				helpIndex = (helpIndex + 1) % 2
				helpAnswer = helpAnswers[helpIndex]
				termWidth, termHeight := ui.TerminalDimensions()
				resize(termWidth, termHeight)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				resize(payload.Width, payload.Height)
			case "<Space>", "<Enter>":
				answer()
			case "0":
				review(flashdown.TotalBlackout)
			case "1":
				review(flashdown.IncorrectDifficult)
			case "2":
				review(flashdown.IncorrectEasy)
			case "3":
				review(flashdown.CorrectDifficult)
			case "4":
				review(flashdown.CorrectEasy)
			case "5":
				review(flashdown.PerfectRecall)
			}
		}
		if game.IsFinished() {
			break
		}
	}
}