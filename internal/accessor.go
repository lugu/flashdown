package flashdown

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// DeckAccessor abstract IO operations around Deck handling.
type DeckAccessor interface {
	DeckName() string
	CardsReader() (io.ReadCloser, error)
	MetaReader() (io.ReadCloser, error)
	MetaWriter() (io.WriteCloser, error)
}

type fileAccessor struct {
	filename string
}

func newFileDeckAccessor(filename string) DeckAccessor {
	return &fileAccessor{filename}
}

func (f *fileAccessor) metaFile() string {
	base := filepath.Base(f.filename)
	base = "." + base + ".db"
	dir := filepath.Dir(f.filename)
	return filepath.Join(dir, base)
}

func (f *fileAccessor) CardsReader() (io.ReadCloser, error) {
	fmt.Println(f.filename)
	return os.Open(f.filename)
}

func (f *fileAccessor) MetaReader() (io.ReadCloser, error) {
	return os.Open(f.metaFile())
}

func (f *fileAccessor) MetaWriter() (io.WriteCloser, error) {
	return os.Create(f.metaFile())
}

func (f *fileAccessor) DeckName() string {
	return path.Base(f.filename)
}
