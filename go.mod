module github.com/lugu/flashdown

go 1.15

require (
	fyne.io/fyne/v2 v2.1.2-0.20211103153407-1c0e05f066dd
	github.com/MichaelMure/go-term-markdown v0.1.3
	github.com/aarzilli/nucular v0.0.0-20210408133902-d3dd7b05a80a
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d
	github.com/fatih/color v1.9.0
	github.com/gizak/termui/v3 v3.1.0
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/nsf/termbox-go v0.0.0-20210114135735-d04385b850e8 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
)

// replace fyne.io/fyne/v2 => github.com/lugu/fyne/v2 develop
replace fyne.io/fyne/v2 => github.com/lugu/fyne/v2 v2.1.2-0.20211105180119-8b8e52e0b3a3
