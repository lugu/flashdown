Unsorted list

- GUI: Erase directory does not work.
- GUI: Add a about/help/licence page (probably needed for mobile).
- GUI: Shortcut up and down arrow to select the deck.
- GUI: Consider a grid pattern instead of a long list.
- GUI: Option for a simplified Correct/Incorrect UI.
- GUI: Cards are no more in alphabetical order: sort them before display.
- GUI: Add option to load a Chinese / Japanese / Korean font (and create a bug to automate this choice when a deck/card is loaded based on the needed unicode).

- BUG: More and 100% success; because repeat questions with correct answer after a long time.
- BUG: line matching '^#' is parsed like a question while it can be in a code section. Should use a parser to split the questions.

- FEATURE: shows the name of the deck to give context about the card being asked
        - DeckAccessor returns the name of the deck
        - loadCards adds the name of the deck to the card
        - readCards add the name of the deck to the card
        - updateTitle updates the name of the deck
- FEATURE: implement previous with 'p'
- FEATURE: save to file with 'w'
- FEATURE: save the question in the db for debug purpose
