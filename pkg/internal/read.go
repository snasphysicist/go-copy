package internal

import (
	"io"
	"log"
	"os"
	"time"
)

// Writer implements reading of the input
// & reporting on the progress thereof
type Reader struct {
	path       string
	b          rbuffer
	done       chan struct{}
	pr         *ProgressReporter
	toTransfer uint64
}

func NewReader(path string, b rbuffer, done chan struct{}, pr *ProgressReporter, toTransfer uint64) Reader {
	return Reader{path: path, b: b, done: done, pr: pr, toTransfer: toTransfer}
}

type rbuffer interface {
	Offer([]byte) bool
}

// Start will uninterruptably start the reader
// reading the input and moving the contents to the buffer.
// It reports progress to the progress reporter as it reads,
// and it closes the done channel when it has read toTransfer bytes.
func (r *Reader) Start() {
	readF, err := os.Open(r.path)
	if err != nil {
		panic(err)
	}
	defer readF.Close()
	for {
		buf := make([]byte, 1000)
		n, err := readF.Read(buf)
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
