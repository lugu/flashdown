Unsorted list

- GUI: Erase directory does not work.
- GUI: Add a about/help/licence page (probably needed for mobile).
- GUI: Shortcut up and down arrow to select the deck.
- GUI: Consider a grid pattern instead of a long list.
- GUI: Option for a simplified Correct/Incorrect UI.
- GUI: Cards are no more in alphabetical order: sort them before display.
- GUI: Add option to load a Chinese / Japanese / Korean font (and create a bug
  to automate this choice when a deck/card is loaded based on the needed
  unicode).

- BUG: line matching '^#' is parsed like a question while it can be in a code
  section. Should use a parser to split the questions. Work around, use '>' in
  front.

- FEATURE: save the question in the db for debug purpose
- FEATURE: when argument is a directory, open all markdown files.

- FEATURE: --timeout <seconds> : how to long to wait before skipping a card.
- FEATURE: investigate https://github.com/slaypni/SM-15/blob/master/sm.js
