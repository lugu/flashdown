package flashdown

import (
	"time"
)

// Game represents a learning session.
type Game struct {
	cards    []Card
	decks    []*Deck
	name     string
	index    int
	success  int
	total    int
	finished bool
}

// NewGameFromFiles reads the markdown files to instantiate a Game.
func NewGameFromFiles(useAllCards bool, files []string) (*Game, error) {
	decks := make([]*Deck, len(files))
	name := ""
	for i, file := range files {
		deck, err := NewDeckFromFile(file)
		if err != nil {
			return nil, err
		}
		decks[i] = deck
		if i != 0 {
			name = name + " "
		}
		name = name + file
	}
	return NewGame(name, useAllCards, decks)
}

// NewGameFromFiles reads the markdown files to instantiate a Game.
func NewGameFromAccessor(name string, accessor DeckAccessor) (*Game, error) {
	deck, err := NewDeck(accessor)
	if err != nil {
		return nil, err
	}
	return NewGame(name, false, []*Deck{deck})
}

// NewGame returns a game given a set of markdown files.
func NewGame(name string, useAllCards bool, decks []*Deck) (*Game, error) {
	game := &Game{
		cards: make([]Card, 0),
		decks: decks,
		name:  name,
	}
	for i, deck := range decks {
		cards := deck.SelectBefore(time.Now())
		if useAllCards {
			cards = deck.Cards
		}
		game.cards = append(game.cards, cards...)
		game.success += len(deck.Cards) - len(cards)
		game.total += len(deck.Cards)
		game.decks[i] = deck
	}
	game.cards = ShuffleCards(game.cards)
	return game, nil
}

// Question returns the next question to answer. Idempotent.
func (g *Game) Question() string {
	if len(g.cards) == 0 {
		return "No cards"
	}
	return g.cards[g.index].Question
}

// Question returns the next question to answer. Idempotent.
func (g *Game) Answer() string {
	if len(g.cards) == 0 {
		return "No cards"
	}
	return g.cards[g.index].Answer
}

func (g *Game) DeckName() string {
	if len(g.cards) == 0 {
		return "zero"
	}
	return g.cards[g.index].DeckName
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
		g.cards[g.index].Meta.Review(s)
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
	if len(g.cards) == 0 {
		return true
	}
	return g.finished
}

func (g *Game) Save() {
	for _, d := range g.decks {
		defer d.SaveDeckMeta()
	}
}

func (g *Game) Name() string {
	return g.name
}
