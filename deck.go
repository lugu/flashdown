package flashdown

import (
	"math/rand"
	"os"
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

func (d *Deck) MetaFilename() string {
	base := filepath.Base(d.Filename)
	base = "." + base + ".db"
	dir := filepath.Dir(d.Filename)
	return filepath.Join(dir, base)
}

// DeckSuccess returns the nubmer of success in the deck.
func DeckSuccessNb(d Deck) int {
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

	file, err := os.Open(filename)
	if err != nil {
		return d, err
	}
	d.Filename = filename
	d.Cards, err = readCards(file)
	if err != nil {
		return d, err
	}

	metaFilename := d.MetaFilename()
	metaMap, err := OpenDB(metaFilename)
	if err != nil {
		// create an empty DB if it is missing
		_, errorMissing := os.Stat(metaFilename)
		if os.IsNotExist(errorMissing) {
			file, err := os.Create(metaFilename)
			if err != nil {
				return d, err
			}
			file.Close()
			metaMap = make(map[Digest]*Meta)
		} else {
			return d, err
		}
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

func SaveDeckMeta(d Deck) error {
	metas := make([]Meta, len(d.Cards))
	for i, _ := range d.Cards {
		metas[i] = *d.Cards[i].Meta
	}
	w, err := os.Create(d.MetaFilename())
	if err != nil {
		return err
	}
	return writeDB(w, metas)
}
