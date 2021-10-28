package flashdown

import "time"

type Game struct {
	cards    []Card
	decks    []Deck
	index    int
	success  int
	total    int
	finished bool
}

func NewGame(forceAllCards bool, decks []Deck) *Game {
	game := &Game{
		cards: make([]Card, 0),
		decks: make([]Deck, len(decks)),
	}
	for i, deck := range decks {
		if forceAllCards {
			game.cards = append(game.cards, deck.Cards...)
		} else {
			game.cards = append(game.cards, deck.SelectBefore(time.Now())...)
		}
		game.success += deck.DeckSuccessNb()
		game.decks[i] = deck
	}
	game.cards = ShuffleCards(game.cards)
	game.total = len(game.cards)
	return game
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

func (g *Game) Progress() (index, total int) {
	return g.index + 1, g.total
}

func (g *Game) Success() float32 {
	return (float32(g.success) / float32(g.total)) * 100
}

func (g *Game) IsFinished() bool {
	return g.finished
}

func (g *Game) Save() {
	for _, d := range g.decks {
		defer SaveDeckMeta(d)
	}
}
