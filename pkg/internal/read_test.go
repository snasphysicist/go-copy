package internal_test

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

type mockSource struct {
	opened bool
	closed bool
	toRead io.ReadWriter
}

type mockReadWriter struct {
	rw  io.ReadWriter
	err error
}

func (rw *mockReadWriter) Write(b []byte) (int, error) {
	return rw.rw.Write(b)
}

func (rw *mockReadWriter) Read(b []byte) (int, error) {
	n, _ := rw.rw.Read(b)
	return n, rw.err
}

func (ms *mockSource) Open() error {
	ms.opened = true
	return nil
}

func (ms *mockSource) Close() error {
	ms.closed = true
	return nil
}

func (ms *mockSource) Read(b []byte) (int, error) {
	return ms.toRead.Read(b)
}

type ReadWriterAsAcceptor struct {
	rw io.ReadWriter
}

func (rwaa *ReadWriterAsAcceptor) Offer(b []byte) bool {
	rwaa.rw.Write(b)
	return true
}

func TestReaderOpensSourceFirst(t *testing.T) {
	done := make(chan struct{})
	bb := bytes.NewBuffer(make([]byte, 0))
	rw := mockReadWriter{rw: bb}
	ms := mockSource{toRead: &rw}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(10, done)
	r := internal.NewReader(&ms, &b, done, &pr, 10)
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
	r := internal.NewReader(&ms, &b, done, &pr, 10)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	ms.toRead.Write(randomBytes(9))
	await(func() bool { return pr.BytesRead() == 9 }, time.Second)
	if ms.closed {
		t.Error("Reader closed the source after only 9 bytes read, should wait until 10")
	}
	ms.toRead.Write(randomBytes(1))
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
	r := internal.NewReader(&ms, &b, done, &pr, 10)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	ms.toRead.Write(randomBytes(9))
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
	r := internal.NewReader(&ms, &b, done, &pr, 10)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	bytesIn := randomBytes(9)
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
	r := internal.NewReader(&ms, &b, done, &pr, 10)
	defer ensureStopped(&b, done)
	defer func() { (&rw).err = io.EOF }()
	go r.Start()
	if ms.closed {
		t.Error("Reader closed the source without anything being read")
	}
	bytesIn := randomBytes(4)
	_, _ = rw.rw.Write(bytesIn)
	await(func() bool { return pr.BytesRead() == uint64(len(bytesIn)) }, time.Second)
	if pr.BytesRead() != uint64(len(bytesIn)) {
		t.Errorf("%d read bytes progress was reported, expected %d", pr.BytesRead(), len(bytesIn))
	}
}
