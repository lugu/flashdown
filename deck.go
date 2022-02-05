package flashdown

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

// MetaMap associate the index if the cards with its meta data.
type MetaMap map[Digest]*Meta

// DeckAccessor abstract IO operations around Deck handling.
type DeckAccessor interface {
	DeckName() string
	CardsReader() (io.ReadCloser, error)
	MetaReader() (io.ReadCloser, error)
	MetaWriter() (io.WriteCloser, error)
}

type fileAccessor struct {
	filename string
}

func newFileDeckAccessor(filename string) DeckAccessor {
	return &fileAccessor{filename}
}

func (f *fileAccessor) metaFile() string {
	base := filepath.Base(f.filename)
	base = "." + base + ".db"
	dir := filepath.Dir(f.filename)
	return filepath.Join(dir, base)
}

func (f *fileAccessor) CardsReader() (io.ReadCloser, error) {
	return os.Open(f.filename)
}

func (f *fileAccessor) MetaReader() (io.ReadCloser, error) {
	return os.Open(f.metaFile())
}

func (f *fileAccessor) MetaWriter() (io.WriteCloser, error) {
	return os.Create(f.metaFile())
}

func (f *fileAccessor) DeckName() string {
	return path.Base(f.filename)
}

type Deck struct {
	Cards      []Card
	MetaWriter func() (io.WriteCloser, error)
}

func loadCards(accessor DeckAccessor) ([]Card, error) {
	cardReader, err := accessor.CardsReader()
	if err != nil {
		return nil, err
	}
	defer cardReader.Close()

	cards, err := readCards(cardReader)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func newMetaCards(accessor DeckAccessor) error {
	metas := make([]Meta, 0)
	metaWriter, err := accessor.MetaWriter()
	if err != nil {
		return fmt.Errorf("Cannot create DB: %v", err)
	}
	defer metaWriter.Close()

	if err = writeDB(metaWriter, metas); err != nil {
		return fmt.Errorf("Cannot write DB: %v", err)
	}
	return nil
}

func loadMetaMap(accessor DeckAccessor) (MetaMap, error) {
	var metaMap MetaMap = make(map[Digest]*Meta)

	metaReader, err := accessor.MetaReader()
	if err != nil {
		err = newMetaCards(accessor)
		if err != nil {
			return nil, err
		}
		return metaMap, nil
	}
	defer metaReader.Close()

	metas, err := readDB(metaReader)
	if err != nil {
		err = newMetaCards(accessor)
		if err != nil {
			return nil, err
		}
		return metaMap, nil
	}

	for i, meta := range metas {
		metaMap[meta.Hash] = &metas[i]
	}
	return metaMap, nil
}

// NewDeckFromFile reads a Deck from a file.
func NewDeckFromFile(filename string) (*Deck, error) {
	return NewDeck(newFileDeckAccessor(filename))
}

// NewDeck reads a Deck from DeckAccessor
func NewDeck(accessor DeckAccessor) (*Deck, error) {
	cards, err := loadCards(accessor)
	if err != nil {
		return nil, err
	}
	metaMap, err := loadMetaMap(accessor)
	if err != nil {
		return nil, err
	}

	for i, _ := range cards {
		hash := Hash(cards[i])
		meta, ok := metaMap[hash]
		if ok {
			cards[i].Meta = meta
		} else {
			cards[i].Meta = NewMeta(hash)
		}
	}

	return &Deck{
		Cards:      cards,
		MetaWriter: accessor.MetaWriter,
	}, nil
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

func (d *Deck) SaveDeckMeta() error {

	metas := make([]Meta, len(d.Cards))
	for i, _ := range d.Cards {
		metas[i] = *d.Cards[i].Meta
	}
	metaWriter, err := d.MetaWriter()
	if err != nil {
		return err
	}
	defer metaWriter.Close()
	return writeDB(metaWriter, metas)
}
