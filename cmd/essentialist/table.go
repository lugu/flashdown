package main

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
)

type fyneRenderer struct{}

// Render given AST node to given writer.
func (f *fyneRenderer) Render(w io.Writer, source []byte, n ast.Node) error {
	return fmt.Errorf("Not yet implemented.")
}

// AddOptions adds given option to this renderer.
func (f *fyneRenderer) AddOptions(...renderer.Option) {
}

// NewRichTextFromMarkdown configures a RichText widget by parsing the provided markdown content.
//
// Since: 2.1
func NewRichTextFromMarkdown(content string) *widget.RichText {
	return widget.NewRichText(parseMarkdown(content)...)
}

type markdownRenderer []widget.RichTextSegment

func (m *markdownRenderer) AddOptions(...renderer.Option) {}

func (m *markdownRenderer) Render(_ io.Writer, source []byte, n ast.Node) error {
	segs, err := renderNode(source, n, false)
	*m = segs
	return err
}

func renderNode(source []byte, n ast.Node, blockquote bool) ([]widget.RichTextSegment, error) {
	switch t := n.(type) {
	case *ast.Document:
		return renderChildren(source, n, blockquote)
	case *ast.Paragraph:
		children, err := renderChildren(source, n, blockquote)
		if !blockquote {
			linebreak := &widget.TextSegment{Style: widget.RichTextStyleParagraph}
			children = append(children, linebreak)
		}
		return children, err
	case *ast.List:
		items, err := renderChildren(source, n, blockquote)
		return []widget.RichTextSegment{
			&widget.ListSegment{Items: items, Ordered: t.Marker != '*' && t.Marker != '-' && t.Marker != '+'},
		}, err
	case *ast.ListItem:
		texts, err := renderChildren(source, n, blockquote)
		return []widget.RichTextSegment{&widget.ParagraphSegment{Texts: texts}}, err
	case *ast.TextBlock:
		return renderChildren(source, n, blockquote)
	case *ast.Heading:
		text := forceIntoHeadingText(source, n)
		switch t.Level {
		case 1:
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleHeading, Text: text}}, nil
		case 2:
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleSubHeading, Text: text}}, nil
		default:
			textSegment := widget.TextSegment{Style: widget.RichTextStyleParagraph, Text: text}
			textSegment.Style.TextStyle.Bold = true
			return []widget.RichTextSegment{&textSegment}, nil
		}
	case *ast.ThematicBreak:
		return []widget.RichTextSegment{&widget.SeparatorSegment{}}, nil
	case *ast.Link:
		link, _ := url.Parse(string(t.Destination))
		text := forceIntoText(source, n)
		return []widget.RichTextSegment{&widget.HyperlinkSegment{Alignment: fyne.TextAlignLeading, Text: text, URL: link}}, nil
	case *ast.CodeSpan:
		text := forceIntoText(source, n)
		return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleCodeInline, Text: text}}, nil
	case *ast.CodeBlock, *ast.FencedCodeBlock:
		var data []byte
		lines := n.Lines()
		for i := 0; i < lines.Len(); i++ {
			line := lines.At(i)
			data = append(data, line.Value(source)...)
		}
		if len(data) == 0 {
			return nil, nil
		}
		if data[len(data)-1] == '\n' {
			data = data[:len(data)-1]
		}
		return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleCodeBlock, Text: string(data)}}, nil
	case *ast.Emphasis:
		text := string(forceIntoText(source, n))
		switch t.Level {
		case 2:
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleStrong, Text: text}}, nil
		default:
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleEmphasis, Text: text}}, nil
		}
	case *ast.Text:
		text := string(t.Text(source))
		if text == "" {
			// These empty text elements indicate single line breaks after non-text elements in goldmark.
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleInline, Text: " "}}, nil
		}
		text = suffixSpaceIfAppropriate(text, n)
		if blockquote {
			return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleBlockquote, Text: text}}, nil
		}
		return []widget.RichTextSegment{&widget.TextSegment{Style: widget.RichTextStyleInline, Text: text}}, nil
	case *ast.Blockquote:
		return renderChildren(source, n, true)
	case *ast.Image:
		dest := string(t.Destination)
		u, err := storage.ParseURI(dest)
		if err != nil {
			u = storage.NewFileURI(dest)
		}
		return []widget.RichTextSegment{&widget.ImageSegment{Source: u, Title: string(t.Title), Alignment: fyne.TextAlignCenter}}, nil
	case *east.TableCell:
		segs, err := renderChildren(source, n, true)
		if err != nil {
			panic("Failed to parse cell children") // FIXME: delete
			return nil, err
		}
		if len(segs) != 1 {
			panic("More than one segment in a cell")
		}
		return []widget.RichTextSegment{NewTableCell(widget.NewRichText(segs[0]))}, nil

	case *east.TableHeader:
		segs, err := renderChildren(source, n, true)
		if err != nil {
			panic("Failed to parse row") // FIXME: delete
			return nil, err
		}
		cells := make([]*TableCell, len(segs))
		for i, seg := range segs {
			cell, ok := seg.(*TableCell)
			if !ok {
				return nil, fmt.Errorf("Unable to cast element %d to TableCell", i)
			}
			cells[i] = cell
		}
		return []widget.RichTextSegment{&TableRow{cells: cells}}, nil
	case *east.TableRow:
		segs, err := renderChildren(source, n, true)
		if err != nil {
			panic("Failed to parse row") // FIXME: delete
			return nil, err
		}
		cells := make([]*TableCell, len(segs))
		for i, seg := range segs {
			cell, ok := seg.(*TableCell)
			if !ok {
				return nil, fmt.Errorf("Unable to cast element %d to TableCell", i)
			}
			cells[i] = cell
		}
		return []widget.RichTextSegment{&TableRow{cells: cells}}, nil
	case *east.Table:
		segs, err := renderChildren(source, n, blockquote)
		if err != nil {
			panic("Failed to parse table") // FIXME: delete
			return nil, err
		}
		rows := make([]*TableRow, len(segs))
		for i, seg := range segs {
			row, ok := seg.(*TableRow)
			if !ok {
				return nil, fmt.Errorf("Unable to cast element %d to TableCell", i)
			}
			rows[i] = row
		}
		return []widget.RichTextSegment{NewTableSegment(rows)}, nil
	}
	fmt.Printf("===> n: %#v\n", n)
	if n == nil {
		fmt.Printf("===> n is nil!\n")
	}
	return nil, nil
}

func suffixSpaceIfAppropriate(text string, n ast.Node) string {
	next := n.NextSibling()
	if next != nil && next.Type() == ast.TypeInline && !strings.HasSuffix(text, " ") {
		return text + " "
	}
	return text
}

func renderChildren(source []byte, n ast.Node, blockquote bool) ([]widget.RichTextSegment, error) {
	children := make([]widget.RichTextSegment, 0, n.ChildCount())
	for childCount, child := n.ChildCount(), n.FirstChild(); childCount > 0; childCount-- {
		if child == nil {
			continue
		}
		segs, err := renderNode(source, child, blockquote)
		if err != nil {
			return children, err
		}
		children = append(children, segs...)
		child = child.NextSibling()
	}
	return children, nil
}

func forceIntoText(source []byte, n ast.Node) string {
	texts := make([]string, 0)
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				texts = append(texts, string(t.Text(source)))
			}
		}
		return ast.WalkContinue, nil
	})
	return strings.Join(texts, " ")
}

func forceIntoHeadingText(source []byte, n ast.Node) string {
	var text strings.Builder
	ast.Walk(n, func(n2 ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch t := n2.(type) {
			case *ast.Text:
				text.Write(t.Text(source))
			}
		}
		return ast.WalkContinue, nil
	})
	return text.String()
}

func parseMarkdown(content string) []widget.RichTextSegment {
	r := markdownRenderer{}
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(&r))
	err := md.Convert([]byte(content), nil)
	if err != nil {
		fyne.LogError("Failed to parse markdown", err)
	}
	return r
}

type (
	DummyRichTextSegment struct{}
	TableCell            struct {
		widget.BaseWidget
		DummyRichTextSegment
		provider *widget.RichText
		renderer *cellRenderer
	}
	TableRow struct {
		DummyRichTextSegment
		cells []*TableCell
	}
	TableSegment struct {
		DummyRichTextSegment
		rows  []*TableRow
		table *widget.Table
	}
)

// DummyRichTextSegment is used by TableRow and TableCell to conform with RichTextSegment.
func (c *DummyRichTextSegment) Inline() bool                    { return false }
func (c *DummyRichTextSegment) Textual() string                 { panic("not implemented") }
func (c *DummyRichTextSegment) Update(fyne.CanvasObject)        { panic("not implemented") }
func (c *DummyRichTextSegment) Visual() fyne.CanvasObject       { panic("not implemented") }
func (c *DummyRichTextSegment) Select(pos1, pos2 fyne.Position) { panic("not implemented") }
func (c *DummyRichTextSegment) SelectedText() string            { panic("not implemented") }
func (c *DummyRichTextSegment) Unselect()                       { panic("not implemented") }

// Cell implements CreateRenderer and draw the underlaying RichTextSegments using RichText.
func (c *TableCell) CreateRenderer() fyne.WidgetRenderer {
	return c.renderer
}

func NewTableCell(content *widget.RichText) *TableCell {
	cell := &TableCell{
		provider: content,
		renderer: NewCellRenderer(content),
	}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (c *TableCell) updateSegment(content *widget.RichText) {
	c.provider = content
	c.renderer.setObject(c.provider)
	// c.renderer.Refresh()
}

func NewTableSegment(rows []*TableRow) *TableSegment {
	// func NewTable(length func() (rows int, cols int), create func() fyne.CanvasObject, update func(TableCellID, fyne.CanvasObject)) *Table {
	table := widget.NewTable(
		func() (int, int) {
			if len(rows) > 0 {
				return len(rows), len(rows[0].cells)
			}
			return 0, 0
		},
		func() fyne.CanvasObject {
			return NewTableCell(
				widget.NewRichText(
					&widget.TextSegment{
						Style: widget.RichTextStyleCodeBlock,
						Text:  "content is unknown yet",
					}))
		},
		func(pos widget.TableCellID, o fyne.CanvasObject) {
			if pos.Row >= len(rows) {
				panic("Cannot create cell at row")
			}
			if pos.Col >= len(rows[pos.Row].cells) {
				panic("Cannot create cell at col")
			}
			cell := o.(*TableCell)
			cell.updateSegment(rows[pos.Row].cells[pos.Col].provider)
		},
	)
	// FIXME: Headers used for troubleshooting. Delete them.
	table.ShowHeaderColumn = true
	table.ShowHeaderRow = true
	table.HideSeparators = false
	return &TableSegment{
		rows:  rows,
		table: table,
	}
}

// Visual returns the graphical elements required to render this segment.
func (l *TableSegment) Visual() fyne.CanvasObject {
	return l.table
}

// Update applies the current state of this table segment to an existing visual.
func (l *TableSegment) Update(o fyne.CanvasObject) {}

// cellRenderer implements fyne.WidgetRenderer
type cellRenderer struct {
	// objects has exactly one object
	objects []fyne.CanvasObject
}

func NewCellRenderer(object fyne.CanvasObject) *cellRenderer {
	return &cellRenderer{[]fyne.CanvasObject{object}}
}

func (r *cellRenderer) setObject(object fyne.CanvasObject) {
	r.objects[0] = object
}

// Destroy does nothing in this implementation.
func (r *cellRenderer) Destroy() {
}

// Layout updates the contained object to be the requested size.
func (r *cellRenderer) Layout(s fyne.Size) {
	r.objects[0].Resize(s)
}

// MinSize returns the smallest size that this render can use, returned from the underlying object.
func (r *cellRenderer) MinSize() fyne.Size {
	return r.objects[0].MinSize()
}

// Objects returns the objects that should be rendered.
func (r *cellRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh requests the underlying object to redraw.
func (r *cellRenderer) Refresh() {
	r.objects[0].Refresh()
}
