package flashdown

import (
	"time"

	"github.com/lugu/flashdown/internal"
)

type Game struct {
	cards    []internal.Card
	decks    []internal.Deck
	index    int
	success  int
	total    int
	finished bool
}

// NewGame returns a game given a set of markdown files.
func NewGame(forceAllCards bool, files []string) (*Game, error) {
	game := &Game{
		cards: make([]internal.Card, 0),
		decks: make([]internal.Deck, len(files)),
	}
	for i, file := range files {
		deck, err := internal.OpenDeck(file)
		if err != nil {
			return nil, err
		}
		if forceAllCards {
			game.cards = append(game.cards, deck.Cards...)
		} else {
			game.cards = append(game.cards, deck.SelectBefore(time.Now())...)
		}
		game.success += deck.DeckSuccessNb()
		game.decks[i] = deck
	}
	game.cards = internal.ShuffleCards(game.cards)
	game.total = len(game.cards)
	return game, nil
}

func (g *Game) Question() string {
	if len(g.cards) == 0 {
		return "No cards"
	}
	return g.cards[g.index].Question
}

func (g *Game) Answer() string {
	if len(g.cards) == 0 {
		return "No cards"
	}
	return g.cards[g.index].Answer
}

// Score represents how easly one responded to a question.
type Score int

const (
	// 0: Total blackout, complete failure to recall the information.
	TotalBlackout Score = 0
	// 1: Incorrect response, but upon seeing the correct answer it felt familiar.
	IncorrectDifficult Score = iota
	// 2: Incorrect response, but upon seeing the correct answer it seemed easy to remember.
	IncorrectEasy Score = iota
	// 3: Correct response, but required significant difficulty to recall.
	CorrectDifficult Score = iota
	// 4: Correct response, after some hesitation.
	CorrectEasy Score = iota
	// 5: Correct response with perfect recall.
	PerfectRecall Score = iota
)

func (g *Game) Review(s Score) {
	if g.index < len(g.cards) {
		if s >= 3 {
			g.success++
		}
		g.cards[g.index].Meta.Review(internal.Score(s))
		g.index++
	}
	if g.index == len(g.cards) {
		g.index = 0
		g.finished = true
	}
}

func (g *Game) Skip() {
	if g.index < len(g.cards) {
		g.index++
	}
	if g.index == len(g.cards) {
		g.index = 0
		g.finished = true
	}
}

func (g *Game) Progress() (current, total int) {
	return g.index + 1, g.total
}

func (g *Game) Success() float32 {
	if g.total == 0 {
		return 100
	}
	return (float32(g.success) / float32(g.total)) * 100
}

func (g *Game) IsFinished() bool {
	return g.finished
}

func (g *Game) Save() {
	for _, d := range g.decks {
		defer d.SaveDeckMeta()
	}
}
