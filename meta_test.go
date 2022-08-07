package flashdown

import (
	"bytes"
	"testing"
	"time"
)

var (
	metaInput = []Meta{
		Meta{
			Hash:       0,
			NextTime:   time.Unix(0, 0),
			Repetition: 0,
			Easiness:   0.0,
		},
		Meta{
			Hash:       1,
			NextTime:   time.Unix(2, 0),
			Repetition: 1,
			Easiness:   1.0,
		},
		Meta{
			Hash:       1,
			NextTime:   time.Unix(3, 0),
			Repetition: 1,
			Easiness:   2.0,
		},
	}
)

func TestStrip(t *testing.T) {
	input := []string{
		"\n",
		"\t",
		" ",
		"A",
		"A?",
		"This \n is a long test",
	}
	expected := []string{
		"",
		"",
		"",
		"a",
		"a",
		"thisisalongtest",
	}
	for i, in := range input {
		out := strip(in)
		if out != expected[i] {
			t.Errorf("%d: %s instead of %s", i, out, expected[i])
		}
	}
}

func TestHash(t *testing.T) {
	input := []string{
		"Question a",
		"question A",
		"question A???",
	}
	var card Card
	card.Question = "question a"
	expected := Hash(card)
	for _, in := range input {
		card.Question = in
		out := Hash(card)
		if out != expected {
			t.Errorf("%v instead of %v", out, expected)
		}
	}
}

func TestMetaReview(t *testing.T) {

	var card Card
	meta := NewMeta(card)
	if meta.Repetition != 0 {
		t.Errorf("Invalid repetition: %d", meta.Repetition)
	}
	if time.Now().Before(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	errors := []Score{
		0, 1, 2,
	}
	for _, score := range errors {
		meta.Review(score)
		if meta.Repetition != 0 {
			t.Errorf("%d: Invalid repetition: %d", score, meta.Repetition)
		}
		if time.Now().Before(meta.NextTime) {
			t.Errorf("%d: Invalid next time: %v", score, meta.NextTime)
		}
	}
	easiness := meta.Easiness
	meta.Review(3)
	if meta.Repetition != 1 {
		t.Errorf("Invalid repetition: %d", meta.Repetition)
	}
	if time.Now().AddDate(0, 0, 1).Before(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	if meta.Easiness > easiness {
		t.Errorf("Easiness should not raise: %f", meta.Easiness)
	}
	easiness = meta.Easiness
	meta.Review(4)
	if meta.Repetition != 2 {
		t.Errorf("Invalid repetition: %d", meta.Repetition)
	}
	if time.Now().AddDate(0, 0, 6).Before(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	if meta.Easiness != easiness {
		t.Errorf("Easiness not the same: %f", meta.Easiness)
	}
	easiness = meta.Easiness
	meta.Review(5)
	if meta.Repetition != 3 {
		t.Errorf("Invalid repetition: %d", meta.Repetition)
	}
	days := 6*2*float64(easiness) - 1
	if time.Now().AddDate(0, 0, int(days)).After(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	days = 6*2*float64(easiness) + 1
	if time.Now().AddDate(0, 0, int(days)).Before(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	if meta.Easiness <= easiness {
		t.Errorf("Easiness should raise: %f", meta.Easiness)
	}
	meta.Review(2)
	if meta.Repetition != 0 {
		t.Errorf("Invalid repetition: %d", meta.Repetition)
	}
	if time.Now().Before(meta.NextTime) {
		t.Errorf("Invalid next time: %v", meta.NextTime)
	}
	for i := 0; i < 100; i++ {
		meta.Review(3)
	}
	if meta.Easiness < 1.3 {
		t.Errorf("Easiness should plateau: %f", meta.Easiness)
	}
}

func TestWriteRead(t *testing.T) {
	var buf bytes.Buffer
	err := writeDB(&buf, metaInput)
	if err != nil {
		t.Fatal(err)
	}
	output, err := readDB(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(metaInput) != len(output) {
		t.Errorf("len %d, expected: %d", len(output), len(metaInput))
	}
	for i, _ := range output {
		if metaInput[i].Hash != output[i].Hash {
			t.Errorf("%d, Hash: %d / %d", i,
				output[i].Hash, metaInput[i].Hash)

		}
	}
}
