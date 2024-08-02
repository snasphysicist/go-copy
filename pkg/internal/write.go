package internal

import (
	"os"
	"time"
)

// syncEachBytes specifies after approximately
// how many written bytes we will try to
// force to flush to disk
const syncEachBytes = uint64(1000000)

// Writer implements writing of the output taking it from the buffer
// & reporting on the progress thereof
type Writer struct {
	path       string
	b          wbuffer
	done       chan struct{}
	pr         *ProgressReporter
	toTransfer uint64
}

// NewWriter creates a new Writer, writing to the file at path from the buffer b,
// signalling when it's done on done, reporting progress to pr,
// and knowing when its done when it has transferred toTransfer bytes
func NewWriter(path string, b wbuffer, done chan struct{}, pr *ProgressReporter, toTransfer uint64) Writer {
	return Writer{path: path, b: b, done: done, pr: pr, toTransfer: toTransfer}
}

// wbuffer has the required method on the buffer that the Writer takes from
type wbuffer interface {
	Pop() ([]byte, error)
}

// Start starts the writer writing to the output
// uninterruptably. It first deletes the file
// before starting to pull from the buffer and
// write the buffer contents out to the file.
// It reports progress to the progress reporter
// as it goes, and will close done when it has
// written out toTransfer bytes.
func (w *Writer) Start() {
	os.Remove(w.path)
	writeF, err := os.Create(w.path)
	if err != nil {
		panic(err)
	}
	defer writeF.Close()
	syncIncrement := uint64(0)
	for {
		next, err := w.b.Pop()
		n := len(next)
		if n == 0 {
			time.Sleep(1 * time.Millisecond)
		}
		if err == nil {
			writeF.Write(next)
			w.pr.ReportBytesWritten(uint64(n))
		}
		newSyncIncrement := w.pr.BytesWritten() / syncEachBytes
		if syncIncrement != newSyncIncrement {
			writeF.Sync()
			syncIncrement = newSyncIncrement
		}
		if w.pr.BytesWritten() == w.toTransfer {
			close(w.done)
			return
		}
	}
}
