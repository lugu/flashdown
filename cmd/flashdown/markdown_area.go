package main

import (
	"image"

	"github.com/fatih/color"
	. "github.com/gizak/termui/v3"

	markdown "github.com/MichaelMure/go-term-markdown"
)

type MarkdownArea struct {
	Block
	Text      string
	TextStyle Style
}

func NewMarkdownArea() *MarkdownArea {
	return &MarkdownArea{
		Block:     *NewBlock(),
		TextStyle: StyleClear,
	}
}

// convertAnsi walks cells and replace terminal output sequence with
// the corresponding Style.
func convertAnsi(cells []Cell) []Cell {

	const (
		AttrUnderline Modifier = 1 << 13
	)
	modifierOn := func(mod Modifier) func(Style) Style {
		return func(s Style) Style {
			s.Modifier = mod
			return s
		}
	}
	modifierClear := func(s Style) Style {
		s.Modifier = ModifierClear
		return s
	}
	NormalText := modifierClear
	BoldText := modifierOn(ModifierBold)
	BoldTextOff := modifierClear
	ItalicText := modifierOn(AttrUnderline)
	ItalicTextOff := modifierClear

	styles := map[string]func(s Style) Style{
		"\x1b[0m":  NormalText, // Turn off all attributes
		"\x1b[1m":  BoldText,
		"\x1b[22m": BoldTextOff,
		"\x1b[3m":  ItalicText,
		"\x1b[23m": ItalicTextOff,
	}

	i := 0
	isEscaping := false
	escapeSeq := ""
	updateStyle := NormalText
	for _, x := range cells {
		if x.Rune == rune('\x1b') { // ESC
			isEscaping = true
			escapeSeq = ""
		}
		if isEscaping {
			escapeSeq = escapeSeq + string(x.Rune)
		} else {
			x.Style = updateStyle(x.Style)
			cells[i] = x
			i++
		}
		if isEscaping && x.Rune == rune('m') {
			isEscaping = false
			var ok bool
			updateStyle, ok = styles[escapeSeq]
			if !ok {
				updateStyle = NormalText
			}
		}
	}
	cells = cells[:i]
	return cells
}

func (self *MarkdownArea) Draw(buf *Buffer) {
	self.Block.Draw(buf)

	color.NoColor = true
	output := markdown.Render(self.Text, self.Inner.Dx(), 0)
	text := string(output)

	cells := ParseStyles(text, self.TextStyle)
	cells = convertAnsi(cells)
	cells = WrapCells(cells, uint(self.Inner.Dx()))

	rows := SplitCells(cells, '\n')

	width := 0
	height := len(rows)
	for y, _ := range rows {
		if width < len(rows[y]) {
			width = len(rows[y])
		}
	}
	y := (self.Inner.Dy() - height) / 2
	if y < 0 {
		y = 0
	}

	for _, row := range rows {
		dx := (self.Inner.Dx() - width) / 2
		if dx < 0 {
			dx = 0
		}
		if y+self.Inner.Min.Y >= self.Inner.Max.Y {
			break
		}
		row = TrimCells(row, self.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X+dx, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(self.Inner.Min))
		}
		y++
	}
}
