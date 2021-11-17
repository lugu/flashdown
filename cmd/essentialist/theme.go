package main

import (
	fyne "fyne.io/fyne/v2"
)

type myTheme struct {
	fyne.Theme
}

func NewTheme(theme fyne.Theme) fyne.Theme {
	return myTheme{theme}
}

func (t myTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return t.Theme.Font(s)
	}
	if s.Bold {
		if s.Italic {
			return resourceNotoSansBoldItalicTtf
		}
		return resourceNotoSansBoldTtf
	}
	if s.Italic {
		return resourceNotoSansItalicTtf
	}
	return resourceNotoSansRegularTtf
}
