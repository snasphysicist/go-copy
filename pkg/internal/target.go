package internal

import "os"

// writingFile provides deletion,
// creation, writing and flushing
// on a file being written to
type writingFile struct {
	path string
	f    *os.File
}

// NewWritingFile provides the writingFile's
// operations on file at given path
func NewWritingFile(path string) writingFile {
	return writingFile{path: path}
}

// Initialise deletes any exisiting file
// at wf.path and creates a fresh one there,
// making it ready for writing. An error
// is returned if deletion or creation fails.
func (wf *writingFile) Initialise() error {
	err := os.Remove(wf.path)
	becauseFileNotExists := os.IsNotExist(err)
	if err != nil && !becauseFileNotExists {
		return err
	}
	f, err := os.Create(wf.path)
	if err != nil {
		return err
	}
	wf.f = f
	return nil
}

// Sync calls os.File.Sync on the underlying os.File
func (wf *writingFile) Sync() error {
	return wf.f.Sync()
}

// Write exposes io.Writer on the underlying file
func (wf *writingFile) Write(b []byte) (int, error) {
	return wf.f.Write(b)
}

// Close exposes io.Closer on the underlying file
func (wf *writingFile) Close() error {
	return wf.f.Close()
}
