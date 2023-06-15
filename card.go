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
	errCardEmpty       = errors.New("Empty note")
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
	cards := make([]string, 0)
	isCode := false // true when parsing "```"
	card := ""      // current card being parsed
	lines := strings.Split(md, "\n")

	for _, line := range lines {
		if splitQuestion.Match([]byte(line)) && !isCode {
			// 1. add previous card to the deck if any.
			// 2. start the card with the title.
			if card != "" {
				cards = append(cards, card)
			}
			card = line
		} else {
			if strings.HasPrefix(line, "```") {
				isCode = !isCode
			}
			// If this isn't a title, add it to the card.
			card = fmt.Sprintf("%s\n%s", card, line)
		}
	}
	if card != "" {
		cards = append(cards, card)
	}
	return cards
}

func parseCards(md string) ([]Card, error) {
	cards := make([]Card, 0)

	sheets := splitCards(md)
	for _, sheet := range sheets {
		card, err := loadCard(sheet)
		if err == errCardEmpty {
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
		return c, errCardEmpty
	}
	sheets := strings.SplitN(md, "\n", 2)
	if len(sheets) != 2 {
		return c, errInvalidCard
	}
	// Remove the '#' from the question.
	c.Question = trim(sheets[0][1:])
	c.Answer = trim(sheets[1])
	c.Meta = NewMeta(c)
	return c, nil
}
