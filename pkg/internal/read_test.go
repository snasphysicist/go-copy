package internal_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/snasphysicist/go-copy/pkg/internal"
	"github.com/snasphysicist/go-copy/pkg/random"
)

// mockSource is a source for the reader
// which allows us to inject a io.ReadWriter
// to imitate the actual source of bytes,
// also tracks whether it has been
// opened and/or closed
type mockSource struct {
	opened bool
	closed bool
	toRead io.ReadWriter
}

// mockReadWriter wraps an io.ReadWriter
// and allows us to inject and error for testing
type mockReadWriter struct {
	rw  io.ReadWriter
	err error
}

// Write implements io.ReadWriter on mockReadWriter,
// passes through to mockReadWriter's rw
func (rw *mockReadWriter) Write(b []byte) (int, error) {
	return rw.rw.Write(b)
}

// Read implements io.ReadWriter on mockReadWriter,
// passes through to mockReadWriter's rw
func (rw *mockReadWriter) Read(b []byte) (int, error) {
	n, _ := rw.rw.Read(b)
	return n, rw.err
}

// Open implements rsource on mockSource,
// never errors and records that it was called
func (ms *mockSource) Open() error {
	ms.opened = true
	return nil
}

// Close implements rsource on mockSource,
// never errors and records that it was called
func (ms *mockSource) Close() error {
	ms.closed = true
	return nil
}

// Read implements rsource on mockSource,
// passes through to mockSource's toRead
func (ms *mockSource) Read(b []byte) (int, error) {
	return ms.toRead.Read(b)
}

// ReadWriterAsAcceptor adapts an io.ReadWriter
// to be used as an Acceptor in ensureStopped
type ReadWriterAsAcceptor struct {
	rw io.ReadWriter
}

// Offer implements Acceptor on ReadWriterAsAcceptor
// passing through to the wrapped Write & always returning true
func (rwaa *ReadWriterAsAcceptor) Offer(b []byte) bool {
	_, _ = rwaa.rw.Write(b)
	return true
}

func TestReaderOpensSourceFirst(t *testing.T) {
	done := make(chan struct{})
	bb := bytes.NewBuffer(make([]byte, 0))
	rw := mockReadWriter{rw: bb}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10, 1)
	defer ensureStopped(&ReadWriterAsAcceptor{rw: bb}, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	await(func() bool { return ms.opened }, time.Second)
	if !ms.opened {
		t.Error("Reader did not open the source after a second")
	}
}

func TestReaderClosesSourceWhenTargetNumberOfBytesRead(t *testing.T) {
	done := make(chan struct{})
	rw := mockReadWriter{rw: bytes.NewBuffer(make([]byte, 0))}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10, 1)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	ms.toRead.Write(random.Bytes(9))
	await(func() bool { return pr.BytesRead() == 9 }, time.Second)
	if ms.closed {
		t.Error("Reader closed the source after only 9 bytes read, should wait until 10")
	}
	ms.toRead.Write(random.Bytes(1))
	rw.err = io.EOF
	await(func() bool { return ms.closed }, time.Second)
	if !ms.closed {
		t.Error("Reader did not close the source a second after all 10 bytes were read")
	}
}

func TestReaderStopsEarlyAndClosesSourceWhenEOFFromSourceReader(t *testing.T) {
	done := make(chan struct{})
	rw := mockReadWriter{rw: bytes.NewBuffer(make([]byte, 0))}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10, 1)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	ms.toRead.Write(random.Bytes(9))
	await(func() bool { return pr.BytesRead() == 9 }, time.Second)
	if ms.closed {
		t.Error("Reader closed the source after only 9 bytes read, should wait until 10")
	}
	rw.err = io.EOF
	await(func() bool { return ms.closed }, time.Second)
	if !ms.closed {
		t.Error("Reader did not close the source a second after EOF was returned by the source reader")
	}
}

func TestReaderOffersReadBytesToBuffer(t *testing.T) {
	done := make(chan struct{})
	rw := mockReadWriter{rw: bytes.NewBuffer(make([]byte, 0))}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10, 1)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	bytesIn := random.Bytes(9)
	_, _ = rw.rw.Write(bytesIn)
	await(func() bool { return pr.BytesRead() == 9 }, time.Second)
	bytesOut, _ := b.Pop()
	if len(bytesIn) != len(bytesOut) {
		t.Errorf(
			"%d bytes moved to the buffer, but %d should have been",
			len(bytesOut),
			len(bytesIn),
		)
	}
}

func TestReaderReportsNumberOfBytesReadToProgressReporter(t *testing.T) {
	done := make(chan struct{})
	rw := mockReadWriter{rw: bytes.NewBuffer(make([]byte, 0))}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10, 1)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	bytesIn := random.Bytes(4)
	_, _ = rw.rw.Write(bytesIn)
	await(func() bool { return pr.BytesRead() == uint64(len(bytesIn)) }, time.Second)
	if pr.BytesRead() != uint64(len(bytesIn)) {
		t.Errorf("%d read bytes progress was reported, expected %d", pr.BytesRead(), len(bytesIn))
	}
}
