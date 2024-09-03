
- [ ] FEATURE: add a timer (count down), per flashcard. After X seconds, the
   card is blank for 5s, then it shows the response.
- [ ] BUG: Essentialist doesn't show the real number of cards in the deck
- [ ] BUG: Essentialist when starting a 100% deck, the question "no cards" is shown!?
- [ ] BUG: When an element from the home list is focussed, shortcuts don't work
- [ ] BUG: Essentialist: Recursive search on the directory
- [ ] BUG: DNS and Github: book the name Essentialist and flashdown

  - When pressing tab to circle the focussed element, the canvas TypeKey
  handler isn't called anymore and it isn't active anymore (pressing 'h'
  doesn't works).

  - List's TypeKey method is only called when focus is given to the whole list.
  In this case, the buttons aren't focussed. This is done by pressing tab 4
  times (before the first element get high lighted). Strangely, when pressing
  up and down on the list focussed, it scrolls the list without highlighting
  any element.

  - The focussed element get its TypedKey method called instead of the canvas
  TypedKey method.

  - Requirements:

    - Menu short cuts work ('s', 'h', 'quick-session')
    - List navigation is works with up/down
    - List entry re-act to 'enter' and 'space' to start a session.

  - Options

    - Option 1: Cannot override all the widgets TypedKey handler (ex: button)
    - Option 2: Prevent all the widgets TypedKey handler from being called when focussed.
    - Option 3: Replace list buttons with non focusable element (ex: labels)

From the definition of `processKeyPressed`  at `v2@v2.5.1/internal/driver/glfw/window.go`:

```go
 // No shortcut detected, pass down to TypedKey
 focused := w.canvas.Focused()
 if focused != nil {
  w.QueueEvent(func() { focused.TypedKey(keyEvent) })
 } else if w.canvas.onTypedKey != nil {
  w.QueueEvent(func() { w.canvas.onTypedKey(keyEvent) })
 }
```
