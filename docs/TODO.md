## Public Repo launch milestone

- [x] Github domain is registered
- [x] DNS name Essentialist.app is registered
- [x] CI generate builds for Windows, Mac, Linux and Android
- [x] App runs on Linux, Mac and Windows
- [x] Build instructions tested on Windows, Mac and Linux
- [x] Automate release publication using tags
- [x] Decision on the format: `#` vs `##`
- [ ] Declare project at Google

## Distro launch milestone

- [x] Document how to report bugs and features
- [x] GUI supports CJK option
- [ ] No usability bugs
- [ ] GUI supports tables

- [ ] Logo has been decided
- [ ] App tested on Android, Windows and Mac
- [ ] CI to create release from tag
- [ ] Release published with binaries
- [ ] Packages for flatpack or snap exists
- [ ] Arch AUR package

## App store launch milestone

- [ ] GUI supports binary mode (correct/incorrect)
- [ ] Basic website at <https://essentialist.app>
- [ ] Ask feedback from 3 users at r/GetStudying and r/spacedrepetition/
- [ ] App tested on low end Android and different emulators
- [ ] Instructions to load and save cards from micro sd
- [ ] Decision on gamification: star hunter?
- [ ] License page is setup
- [ ] F-droid version

## Unsorted TODO list

- BUG: GUI doesn't show the real number of cards in the deck
- BUG: GUI when starting a 100% deck, the question "no cards" is shown!?
- BUG: When an element from the home list is focussed, shortcuts don't work
- BUG: when argument is a relative directory directory (`~` or `..`), it fails
- FEATURE: GUI: Recursive search on the directory
- FEATURE: GUI: Text selectable
- FEATURE: investigate <https://github.com/slaypni/SM-15/blob/master/sm.js>
- FEATURE: Complete help with about & licence

## How rendering works in Fyne?

### What is a CanvasObject?

Something that can be drawn on the canvas.

```go
type CanvasObject interface {
  MinSize() Size
  Move(Position)
  Position() Position
  Resize(Size)
  Size() Size
  Hide()
  Visible() bool
  Show()
  Refresh()
}
```

For a primitive `CanvasObject` like `Line`, the `Refresh()` method is:

```go
// Refresh instructs the containing canvas to refresh the specified obj.
func Refresh(obj fyne.CanvasObject) {
 app := fyne.CurrentApp()
 if app == nil || app.Driver() == nil {
  return
 }

 c := app.Driver().CanvasForObject(obj)
 if c != nil {
  c.Refresh(obj)
 }
}
```

### What is a Canvas?

The area within which the application is drawn. Each window has a canvas (`Window.Canvas()`).

The Driver knows on which canvas is the object supposed to be drawn. Then the
canvas has some primitive operations like: draw a rectangle, draw a line, draw
a raster or some text. All the elements on screen a composed from those primitives.

The `Canvas.Refresh(CanvasObject)` method uses a refresh queue to store objects
that need to be redrawn. A `Painter` is used to draw on the canvas.

```go
func (p *painter) Paint(obj fyne.CanvasObject,
                        pos fyne.Position,
                        frame fyne.Size) {
 if obj.Visible() {
  p.drawObject(obj, pos, frame)
 }
}
```

### Widgets

A widget is a `CanvasObject` that is interactive. The `WidgetRenderer` isolate
the widget logic from the way it is represented.

### What is a `RichTextSegment`?

Element that can be displayed by a `widget.RichText`.

```go
type RichTextSegment interface {
 Inline() bool
 Textual() string
 Update(fyne.CanvasObject)
 Visual() fyne.CanvasObject

 Select(pos1, pos2 fyne.Position)
 SelectedText() string
 Unselect()
}
```

### How container works?
