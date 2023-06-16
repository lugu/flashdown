# Flashdown & Essentialist

Programs for [spaced repetition][1] using flashcards in [Markdown][2].
- **Flashdown**: a terminal application
- **Essentialist**: a GUI for desktops and mobiles (Android and iOS)

The space repetition algorithm used is based on [SM-2][3].

Key features:

- **No cloud**: Your data never leave your device. Privacy matters.
- **Markdown**: Flash cards are plain text Markdown files.
- **Keyboard shortcut**: Minimalist interface with keyboard navigation.

Similar project: https://github.com/Yvee1/hascard.

[1]: https://en.wikipedia.org/wiki/Spaced_repetition
[2]: https://en.wikipedia.org/wiki/Markdown
[3]: https://en.wikipedia.org/wiki/SuperMemo.

## Deck syntax

Save and edit your flash cards with in a dead simple Markdown file (like:
`my_deck.md`). Questions are heading level 1 followed with their answers, like:

```markdown
# Question 1

Answer in Markdown.

# Second __question__

- one
- **two**
- three

# Third question

Answer with a table.

|  A  |  B  |
| --- | --- |
| 124 | 456 |
```

## Flashdown

Flashdown is the terminal application.

To install it, clone this repo and run:

```shell
go install ./cmd/flashdown
```

Usage:

```shell
flashdown <deck_file> [<deck_file>]
```

This will automatically create an hidden file `.<deck file>.db` with
the recorded scores.

```
┌─Session: 1/2 — Success 33.33%────────────────────────────────────────────────┐
│                                                                              │
│                                                                              │
│                                Third question                                │
│                                                                              │
│                                                                              │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
┌─Deck: test.md────────────────────────────────────────────────────────────────┐
│                             Answer with a table                              │
│                             ┌───┬───┐                                        │
│                             │A  │B  │                                        │
│                             ╞═══╪═══╡                                        │
│                             │124│456│                                        │
│                             └───┴───┘                                        │
└──────────────────────────────────────────────────────────────────────────────┘
Press [0-5] to continue, 's' to skip or 'q' to quit

5: Perfect response
4: Correct response, after some hesitation
3: Correct response, with serious difficulty
2: Incorrect response, but upon seeing the answer it seemed easy to remember
1: Incorrect response, but upon seeing the answer it felt familiar
0: Total blackout
```

## Essentialist

A GUI version for desktops and mobile (Android, iOS support isn't tested).

> :warning: UTF-8 and Markdown tables are not yet supported.

![Screenshot](docs/essentialist-screenshot.png)

### Desktop installation

```shell
go install ./cmd/essentialist
essentialist
```

### Android installation

```shell
cd cmd/essentialist
fyne package -os android
adb install Essentialist.apk
```

Use the local storage (of your Android device) to import flash cards. For
example, you can put them in an SD card and import them from the Essentialist
application.
