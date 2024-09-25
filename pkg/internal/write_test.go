package internal_test

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

type mockTarget struct {
	buffer         []byte
	destination    []byte
	wasInitialised bool
	wasClosed      bool
}

func (t *mockTarget) Initialise() error {
	t.buffer = make([]byte, 0)
	t.destination = make([]byte, 0)
	t.wasInitialised = true
	return nil
}

func (t *mockTarget) Sync() error {
	t.destination = append(t.destination, t.buffer...)
	t.buffer = make([]byte, 0)
	return nil
}

func (t *mockTarget) Close() error {
	t.wasClosed = true
	return nil
}

func (t *mockTarget) Write(b []byte) (int, error) {
	t.buffer = append(t.buffer, b...)
	return len(b), nil
}

func await(p func() bool, timeout time.Duration) {
	a := time.After(timeout)
	for {
		select {
		case <-a:
			return
		case <-time.After(time.Millisecond):
			if p() {
				return
			}
		}
	}
}

type Acceptor interface {
	Offer([]byte) bool
}

func ensureStopped(b Acceptor, done chan struct{}) {
	attempts := 0
	for {
		select {
		case _, ok := <-done:
			if !ok {
				return
			}
		default:
			b.Offer([]byte{0})
			attempts++
			if attempts > 1000000 {
				panic("Failed to shut down")
			}
		}
	}
}

func TestWriterInitialisesTarget(t *testing.T) {
	done := make(chan struct{})
	mt := mockTarget{}
	b := internal.NewBuffer(100)
	w := internal.NewWriter(
		&mt,
		&b,
		done,
		internal.From(internal.NewProgressReporter(100, done)),
		100,
		1000,
	)
	defer ensureStopped(&b, done)
	go w.Start()
	await(func() bool { return (&mt).wasInitialised }, 2*time.Second)
	if !(&mt).wasInitialised {
		t.Error("The target was not initialised after 2 seconds")
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func TestWriterTakesFromBufferPutsToTargetWhenAvailable(t *testing.T) {
	done := make(chan struct{})
	mt := mockTarget{}
	b := internal.NewBuffer(100)
	w := internal.NewWriter(
		&mt,
		&b,
		done,
		internal.From(internal.NewProgressReporter(100, done)),
		100,
		1000,
	)
	defer ensureStopped(&b, done)
	go w.Start()
	firstData := randomBytes(50)
	b.Offer(firstData)
	await(func() bool { return len(mt.buffer) == 50 }, 2*time.Second)
	if !(len(mt.buffer) == 50) {
		t.Errorf("Expected 50 bytes taken by writer, actually %d", len(mt.buffer))
	}
	secondData := randomBytes(22)
	b.Offer(secondData)
	await(func() bool { return len(mt.buffer) == 72 }, 2*time.Second)
	if !(len(mt.buffer) == 72) {
		t.Errorf("Expected 72 bytes taken by writer, actually %d", len(mt.buffer))
	}
	if !reflect.DeepEqual(append(firstData, secondData...), mt.buffer) {
		t.Errorf(
			"%v is actual buffer content, expected %v followed by %v",
			mt.buffer, firstData, secondData,
		)
	}
}

func TestWriterReportsWrittenBytesToProgressReporter(t *testing.T) {
	done := make(chan struct{})
	mt := mockTarget{}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(100, done)
	w := internal.NewWriter(
		&mt,
		&b,
		done,
		&pr,
		100,
		1000,
	)
	defer ensureStopped(&b, done)
	go w.Start()
	firstData := randomBytes(50)
	b.Offer(firstData)
	await(func() bool { return pr.BytesWritten() == 50 }, 2*time.Second)
	if !(pr.BytesWritten() == 50) {
		t.Errorf("Expected 50 bytes reported written, actually %d", pr.BytesWritten())
	}
}

func TestWriterSyncsWhenEnoughBytesTakenFromBuffer(t *testing.T) {
	done := make(chan struct{})
	mt := mockTarget{}
	b := internal.NewBuffer(100)
	pr := internal.NewProgressReporter(100, done)
	w := internal.NewWriter(
		&mt,
		&b,
		done,
		&pr,
		100,
		15,
	)
	defer ensureStopped(&b, done)
	go w.Start()
	b.Offer(randomBytes(15))
	await(func() bool { return len(mt.destination) == 15 }, 2*time.Second)
	if !(len(mt.destination) == 15) || !(len(mt.buffer) == 0) {
		t.Errorf("%d bytes in destination, %d in buffer, should be 15 and 0 after sync",
			len(mt.destination), len(mt.buffer),
		)
	}
	b.Offer(randomBytes(14))
	await(func() bool { return len(mt.buffer) != 0 }, 2*time.Second)
	if !(len(mt.destination) == 15) || !(len(mt.buffer) == 14) {
		t.Errorf("%d bytes in destination, %d in buffer, should be 15 and 14 after just one sync",
			len(mt.destination), len(mt.buffer),
		)
	}
	b.Offer(randomBytes(2))
	await(func() bool { return len(mt.destination) == 31 }, 2*time.Second)
	if !(len(mt.destination) == 31) || !(len(mt.buffer) == 0) {
		t.Errorf("%d bytes in destination, %d in buffer, should be 31 and 0 after just second sync",
			len(mt.destination), len(mt.buffer),
		)
	}
}

func TestWriterClosesTargetWhenTargetBytesWritten(t *testing.T) {
	done := make(chan struct{})
	mt := mockTarget{}
	b := internal.NewBuffer(2)
	pr := internal.NewProgressReporter(20, done)
	w := internal.NewWriter(
		&mt,
		&b,
		done,
		&pr,
		20,
		100,
	)
	defer ensureStopped(&b, done)
	go w.Start()
	await(func() bool {
		b.Offer(randomBytes(1))
		time.Sleep(10 * time.Millisecond)
		return mt.wasClosed
	}, 2*time.Second)
	if !mt.wasClosed {
		t.Errorf("The target has not been closed despite %d bytes transfered", pr.BytesWritten())
	}
	if !(pr.BytesWritten() == 20) {
		t.Errorf("%d bytes written at close, expected 20", pr.BytesWritten())
	}
}
