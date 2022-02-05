package flashdown

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	errNoteEmpty       = errors.New("Empty note")
	errQuestionMissing = errors.New("Missing question")
	errInvalidCard     = errors.New("Invalid card")
)

var (
	splitQuestion = regexp.MustCompile(`(?m)^#\s*`)
)

type Card struct {
	Question string
	Answer   string
	DeckName string
	Meta     *Meta
}

func (c Card) Review(s Score) {
	c.Meta.Review(s)
}

func splitCards(md string) []string {
	return splitQuestion.Split(md, -1)
}

func parseCards(md string) ([]Card, error) {
	cards := make([]Card, 0)

	sheets := splitCards(md)
	for _, sheet := range sheets {
		card, err := loadCard(sheet)
		if err == errNoteEmpty {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("%w (%s)", err, sheet)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func readCards(r io.Reader) ([]Card, error) {
	dat, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return parseCards(string(dat))
}

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

// loadCard parse a card description
func loadCard(md string) (c Card, err error) {

	md = trim(md)
	if md == "" {
		return c, errNoteEmpty
	}
	sheets := strings.SplitN(md, "\n", 2)
	if len(sheets) != 2 {
		return c, errInvalidCard
	}
	c.Question = trim(sheets[0])
	c.Answer = trim(sheets[1])
	c.Meta = NewMeta(Hash(c))
	return c, nil
}
