package main

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/acarl005/stripansi"
)

func renderMarkdown(md string, charWidth int) string {
	output := markdown.Render(md, charWidth, 0)
	return stripansi.Strip(string(output))
}

func NewMarkdownContainer(md string) *fyne.Container {
	txt := renderMarkdown(md, 80)
	objects := make([]fyne.CanvasObject, 0)
	txtSize := fyne.MeasureText("", 0, fyne.TextStyle{Monospace: true})
	// Empirically adjust the height so that each lines touch each other
	// and tables are drawn correctly.
	height := txtSize.Height + 1.6
	for i, line := range strings.Split(txt, "\n") {
		text := canvas.NewText(line, color.White)
		text.TextStyle.Monospace = true
		objects = append(objects, text)
		text.Move(fyne.NewPos(0, float32(i)*(height)))
	}
	rect := canvas.NewRectangle(color.White)
	rect.SetMinSize(fyne.NewSize(0, float32(len(objects))*height))
	objects = append(objects, rect)
	c := container.NewWithoutLayout(objects...)
	return c
}
