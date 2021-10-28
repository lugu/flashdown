package internal

import (
	"bytes"
	"testing"
)

func TestSplitDeck(t *testing.T) {
	input := `
# L1 cache reference
0.5 ns

#Branch mispredict
5 ns
# L2 cache reference

7 ns,  14x L1 cache
`
	expected := []Card{
		Card{
			Question: "L1 cache reference",
			Answer:   "0.5 ns",
		},
		Card{
			Question: "Branch mispredict",
			Answer:   "5 ns",
		},
		Card{
			Question: "L2 cache reference",
			Answer:   "7 ns,  14x L1 cache",
		},
	}
	cards, err := readCards(bytes.NewBufferString(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(cards) != len(expected) {
		t.Fatalf("wrong length: %d", len(cards))
	}
	for i, card := range cards {
		if card.Question != expected[i].Question {
			t.Errorf("Question: %s, instead of: %s",
				card.Question, expected[i].Question)
		}
		if card.Answer != expected[i].Answer {
			t.Errorf("Answer: %s, instead of: %s",
				card.Answer, expected[i].Answer)
		}
	}
}
