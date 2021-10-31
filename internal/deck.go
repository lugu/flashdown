package internal

import (
	"math/rand"
	"path/filepath"
	"time"
)

// MetaMap associate the index if the cards with its meta data.
type MetaMap map[uint64]Meta

// Deck reprensents a list of cards read from disk.
type Deck struct {
	Filename string
	Cards    []Card
}

func ShuffleCards(cards []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	return cards
}

// SelectBefore returns the cards to be review before a given date.
func (d *Deck) SelectBefore(now time.Time) []Card {
	cards := []Card{}
	for _, card := range d.Cards {
		if card.Meta.NextTime.Before(now) {
			cards = append(cards, card)
		}
	}
	return cards
}

// MetaFilename returns the location of the meta file.
// BUG: on Android the URI is not correctly constructed for content: URI and if
// it was it would points to a non writable path.
func (d *Deck) MetaFilename() string {
	base := filepath.Base(d.Filename)
	base = "." + base + ".db"
	dir := filepath.Dir(d.Filename)
	return filepath.Join(dir, base)
}

// DeckSuccess returns the nubmer of success in the deck.
func (d *Deck) DeckSuccessNb() int {
	var success int
	for _, card := range d.Cards {
		if card.Meta.Repetition > 0 {
			success++
		}
	}
	return success
}

// OpenDeck reads a Deck from disk
func OpenDeck(filename string) (d Deck, err error) {

	file, err := OpenReader(filename)
	if err != nil {
		return d, err
	}
	defer file.Close()
	d.Filename = filename
	d.Cards, err = readCards(file)
	if err != nil {
		return d, err
	}

	metaFilename := d.MetaFilename()
	metaMap, err := OpenDB(metaFilename)
	if err != nil {
		// create an empty DB if it is missing
		file, err := CreateWriter(metaFilename)
		if err != nil {
			return d, err
		}
		file.Close()
		metaMap = make(map[Digest]*Meta)
	}
	for i, _ := range d.Cards {
		hash := Hash(d.Cards[i])
		meta, ok := metaMap[hash]
		if ok {
			d.Cards[i].Meta = meta
		} else {
			d.Cards[i].Meta = NewMeta(hash)
		}
	}
	return d, nil
}

func (d *Deck) SaveDeckMeta() error {
	metas := make([]Meta, len(d.Cards))
	for i, _ := range d.Cards {
		metas[i] = *d.Cards[i].Meta
	}
	w, err := CreateWriter(d.MetaFilename())
	if err != nil {
		return err
	}
	defer w.Close()
	return writeDB(w, metas)
}
