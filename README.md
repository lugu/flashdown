# Flashcards in Markdown

![build](https://github.com/lugu/flashdown/workflows/Build/badge.svg)

A spaced repetition game (based on SM-2) for plain text flashcards.

Similar project: https://github.com/Yvee1/hascard

## Usage

```
flashdown <deck file>
```

This will automatically create an hidden file `.<deck file>.db` with
the recorded scores.

## Deck syntax

Questions and answer can be multiple lines of markdown.

```
# Question

Answer

# Second __question__

An answer with a **list**
- one
- two
- three

# Third question

Answer with a table

|  A  |  B  |
| --- | --- |
| 124 | 456 |
```

## Screenshoot

```
┌─Session: 1/2 — Success 33.33%────────────────────────────────────────────────┐
│                                                                              │
│                                                                              │
│                                Third question                                │
│                                                                              │
│                                                                              │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────────────────────────────┐
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
