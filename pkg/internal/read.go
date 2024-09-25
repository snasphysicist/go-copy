package internal

import (
	"io"
	"log"
	"time"
)

// Reader implements reading of the input into a buffer
// & reporting on the progress thereof
type Reader struct {
	source     rsource
	b          rbuffer
	done       chan struct{}
	pr         *ProgressReporter
	toTransfer uint64
}

// NewReader creates a new Reader, reading from the file at path into the buffer b,
// signalling when it's done on done, reporting progress to pr,
// and knowing when its done when it has transferred toTransfer bytes
func NewReader(source rsource, b rbuffer, done chan struct{}, pr *ProgressReporter, toTransfer uint64) Reader {
	return Reader{source: source, b: b, done: done, pr: pr, toTransfer: toTransfer}
}

// rbuffer has the required method on the buffer that the Reader reads into
type rbuffer interface {
	Offer([]byte) bool
}

// rsource represents a source of bytes to read from
type rsource interface {
	// Open prepares the source to start having its bytes read
	Open() error
	io.ReadCloser
}

// Start will uninterruptably start the reader
// reading the input and moving the contents to the buffer.
// It reports progress to the progress reporter as it reads,
// and it closes the done channel when it has read toTransfer bytes.
func (r *Reader) Start() {
	err := r.source.Open()
	if err != nil {
		panic(err)
	}
	defer r.source.Close()
	for {
		buf := make([]byte, 1000)
		n, err := r.source.Read(buf)
		if err == io.EOF {
			if r.pr.BytesRead() != r.toTransfer {
				log.Printf(
					"WARNING: transferred %d bytes, should have been %d",
					r.pr.BytesRead(), r.toTransfer,
				)
			}
			close(r.done)
			return
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n > 0 {
			success := r.b.Offer(buf[:n])
			for !success {
				success = r.b.Offer(buf[:n])
				time.Sleep(1 * time.Millisecond)
			}
		}
		r.pr.ReportBytesRead(uint64(n))
	}
}
