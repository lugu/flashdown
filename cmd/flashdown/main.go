package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	forceAllCards := false
	var success, total int
	cards := make([]flashdown.Card, 0)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-a" {
			forceAllCards = true
			continue
		}
		deck, err := flashdown.OpenDeck(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}
		defer flashdown.SaveDeckMeta(deck)
		if forceAllCards {
			cards = append(cards, deck.Cards...)
		} else {
			cards = append(cards, deck.SelectBefore(time.Now())...)
		}
		success += flashdown.DeckSuccessNb(deck)
		total += len(deck.Cards)
	}

	if len(cards) == 0 {
		return
	}
	cards = flashdown.ShuffleCards(cards)
	index := 0

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
		percent := (float32(success) / float32(total)) * 100
		q.Title = fmt.Sprintf(`Card: %d/%d — Success %2.0f%%`,
			index+1, len(cards), percent)
	}

	ask := func(c flashdown.Card) {
		updateTitle()
		q.Text = c.Question
		a.Text = ""
		help.Text = helpQuestion
		ui.Clear()
		ui.Render(grid, help)
	}
	review := func(i int) {
		if i >= 3 {
			success++
		}
		cards[index].Review(flashdown.Score(i))
		index++
		if index >= len(cards) {
			return
		}
		ask(cards[index])
	}
	answer := func(c flashdown.Card) {
		q.Text = c.Question
		a.Text = c.Answer
		help.Text = helpAnswer
		ui.Clear()
		ui.Render(grid, help)
	}
	ask(cards[index])

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
				index++
				if index >= len(cards) {
					return
				}
				ask(cards[index])
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
				answer(cards[index])
			case "0":
				review(0)
			case "1":
				review(1)
			case "2":
				review(2)
			case "3":
				review(3)
			case "4":
				review(4)
			case "5":
				review(5)
			}
		}
		if index >= len(cards) {
			break
		}
	}
}
