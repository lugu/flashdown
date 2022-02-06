# Flashdown & Essentialist

Programs for [spaced repetition][1] using flashcards in [Markdown][2].
- **Flashdown**: a terminal application
- **Essentialist**: a GUI for desktops and mobiles

The algorithm used is based on [SM-2][3].

Key features:

- **Privacy**: Your data never leave your device. On mobile they are encrypted.
- **Focus on edition**: It must be as easy as possible to create decks.
- **Productivity**: Minimalist interface with keyboard navigation.

Similar project: https://github.com/Yvee1/hascard.

[1]: https://en.wikipedia.org/wiki/Spaced_repetition
[2]: https://en.wikipedia.org/wiki/Markdown
[3]: https://en.wikipedia.org/wiki/SuperMemo.

## Deck syntax

- A deck is a simple Markdown file (like: `my_deck.md`).
- Questions are heading level 1 followed with their answers.

```markdown
# Question 1

Answer 1

# Second __question__

Answer with a **list**:
- one
- two
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

A GUI version for desktops and mobile **under development**.

> :warning: UTF-8 and Markdown tables are not yet supported.

### Desktop installation

```shell
go install ./cmd/essentialist
essentialist
```

### Mobile installation

```shell
cd cmd/essentialist
fyne package -os android
install Essentialist.apk
```

Decks are imported via local storage (ex: SD card).
