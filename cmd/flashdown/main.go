package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	flashdown "github.com/lugu/flashdown/internal"
)

const (
	helpQuestion = `Press space to continue, 's' to skip or 'q' to quit`

	usageMsg = `Spaced repetition program for flashcards in Markdown.

Usage: %s [-a] [-n <number of cards>] <file or directory> [<file> ...]
Flags:
	-a | --all    : force all cards in the deck to be used.
	-h | --help   : show this message.
	-n | --number : set the number of cards used.
	-d | --debug  : debug logs are written to a temprorary file.

A deck is a plain text Markdown file where questions have heading level 1 like:

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
	if len(os.Args) < 2 {
		fmt.Printf(usageMsg, os.Args[0])
		os.Exit(1)
	}

	cardsNb := flashdown.CARDS_TO_REVIEW
	files := make([]string, 0, len(os.Args))

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a", "--all", "-all":
			cardsNb = flashdown.ALL_CARDS
			continue
		case "-d", "--debug", "-debug":
			file, err := os.CreateTemp(".", "log")
			if err != nil {
				log.Fatal(err)
			}
			log.SetOutput(file)
		case "-h", "--help", "-help":
			fmt.Printf(usageMsg, os.Args[0])
			os.Exit(0)
		case "-n", "--number", "-number":
			var err error
			if i < len(os.Args) {
				i++
				cardsNb, err = strconv.Atoi(os.Args[i])
			}
			if cardsNb <= 0 || err != nil {
				fmt.Print("Argument -n must be followed by a positive number.\n")
				os.Exit(1)
			}
			continue
		}

		file := os.Args[i]
		// If file is a directory. Add all markdown files in the directory.
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Cannot open %s: %s.\n", file, err)
			os.Exit(1)
		}
		info, err := f.Stat()
		if err != nil {
			fmt.Printf("Cannot access %s: %s.\n", file, err)
			os.Exit(1)
		}
		if info.IsDir() {
			entries, err := f.ReadDir(-1)
			if err != nil {
				fmt.Printf("Cannot list files inside %s.\n", file)
				os.Exit(1)
			}
			for _, entry := range entries {
				// skip diretories
				if entry.IsDir() || path.Ext(entry.Name()) != ".md" {
					continue
				}
				filename := path.Join(path.Base(file), entry.Name())
				files = append(files, filename)
			}
		} else {
			files = append(files, file)
		}
	}

	game, err := flashdown.NewGameFromFiles(cardsNb, files)
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
		a.Title = fmt.Sprintf(`Deck: %s`, game.DeckName())
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
			case "s", "n":
				game.Skip()
				if game.IsFinished() {
					return
				}
				ask()
			case "w":
				game.Save()
			case "p":
				game.Previous()
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
