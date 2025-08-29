# Essentialist

Programs for [spaced repetition][1] using flashcards in [Markdown][2].

- **Essentialist**: Application for desktops and mobiles
- **Flashdown**: Console application

The space repetition algorithm used is based on [SM-2][3].

Key features:

- **Privacy**: No cloud, your data never leave your device.
- **Easy card creation**: Flash cards are plain text Markdown files.
- **Cross-platform**: Runs on Linux, MacOS, Windows and android.

[1]: https://en.wikipedia.org/wiki/Spaced_repetition
[2]: https://en.wikipedia.org/wiki/Markdown
[3]: https://en.wikipedia.org/wiki/SuperMemo#Description_of_SM-2_algorithm

See the [CONTRIBUTING.md](/.github/CONTRIBUTING.md) for how to report bugs and
submit pull request.

## Flash cards syntax

Each deck of cards is a plain text Markdown files with the extension `.md` (ex:
`sample.md`). You can put all your decks in the same directory.

Each card starts with a heading level 2 (line starting with `##`) defining the
question. The answer is the content following (until the next heading level 2).

You progress is stored in a hidden file `.<deck file>.db` (ex: `.sample.md.db`).

Example of a deck with 3 cards:

```markdown
## Question: what format is used?

Questions and answers are in **Markdown**.

## Are lists supported?

Yes, here is an example:

- one
- **two**
- three

## How to include a table in the answer?

Answer with a table.

|  A  |  B  |
| --- | --- |
| 124 | 456 |
```

## Essentialist (GUI)

A GUI version for desktops and mobile (Android, iOS support isn't tested).

![Screenshot](docs/essentialist-screenshot.png)

### Installation

Download the latest version of Essentialist (available
[here](https://github.com/lugu/flashdown/releases)) or compile it with the
following instructions:

<details><summary>Linux</summary>
<p>

```shell
go install ./cmd/essentialist
```

</p>
</details>

<details><summary>MacOS</summary>
<p>

```shell
CGO_ENABLED=1 go build ./cmd/essentialist
./essentialist
```

</p>
</details>

<details><summary>Windows</summary>
<p>

```shell
go build -x -o essentialist.exe ./cmd/essentialist
```

</p>
</details>

<details><summary>Android</summary>
<p>

1. Install the Android NDK from <https://developer.android.com/ndk/downloads>.
   Set the `ANDROID_NDK_HOME` variable to the directory where the NDK is located.

1. Build the Android APK with:

  ```shell
  cd cmd/essentialist
  fyne package -os android
  ```

1. Plug your phone over USB and install the APK with:

  ```shell
  adb install Essentialist.apk
  ```

Use the local storage (of your Android device) to import flash cards. For
example, you can put them in an SD card and import them from the Essentialist
application.

</p>
</details>

## Flashdown (terminal version)

Flashdown is the terminal application.

To install it, clone this repo and run:

```shell
go install ./cmd/flashdown
```

Usage:

```shell
flashdown <deck_file> [<deck_file>]
```

![Screenshot](docs/flashdown-screenshot.png)

Similar project: <https://github.com/Yvee1/hascard>.

## Maintenance

### How to update dependencies

```shell
go get -u all
go mod tidy
go run github.com/dennwc/flatpak-go-mod@latest .
mv modules.txt cmd/essentialist/flatpak/
cat go.mod.yml >> cmd/essentialist/flatpak/io.github.lugu.essentialist.yml
```
