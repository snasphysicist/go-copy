package internal

import (
	"io"
	"os"
)

// SourceFile represents a file as a source
// to be read in to this program
type SourceFile struct {
	path string
	rc   io.ReadCloser
}

// NewSourceFile creates a new SourceFile,
// opening and reading the file from the given path
func NewSourceFile(path string) SourceFile {
	return SourceFile{path: path}
}

// Open attempts to open the file for reading
// at sf.path, returning an error when this fails
func (sf *SourceFile) Open() error {
	rc, err := os.Open(sf.path)
	sf.rc = rc
	return err
}

// Read implements io.ReadCloser on SourceFile,
// passes through to Read on the underlying file
func (sf *SourceFile) Read(b []byte) (int, error) {
	return sf.rc.Read(b)
}

// Close implements io.ReadCloser on SourceFile,
// passes through to Close on the underlying file
func (sf *SourceFile) Close() error {
	return sf.rc.Close()
}
