package copy

import (
	"time"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

// FileToFile copies a single file, from the from path to the to path,
// using a buffer of size bufferSizeBytes in bytes & forcing write
// to target durable storage after each approximately syncEachBytes
func FileToFile(from string, to string, bufferSizeBytes uint64, syncEachBytes uint64) {
	s := internal.SizeOf(from)

	crossBuffer := internal.NewBuffer(bufferSizeBytes)

	shutdown := make(chan struct{})

	readerDone := make(chan struct{})
	writerDone := make(chan struct{})

	pr := internal.NewProgressReporter(s, shutdown)
	readingFile := internal.NewSourceFile(from)
	reader := internal.NewReader(&readingFile, &crossBuffer, readerDone, &pr, s)
	writingFile := internal.NewWritingFile(to)
	writer := internal.NewWriter(&writingFile, &crossBuffer, writerDone, &pr, s, syncEachBytes)

	go pr.Report(time.Now())
	go reader.Start()
	go writer.Start()

	<-readerDone
	<-writerDone

	close(shutdown)
	time.Sleep(10 * time.Millisecond)
}
