package internal_test

import (
	"bytes"
	"testing"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

func TestBufferAcceptsOfferedBytesThatDoNotFillBufferWhenNotFull(t *testing.T) {
	b := internal.NewBuffer(100)
	if !b.Offer(make([]byte, 10)) {
		t.Error("Unfilled buffer did not accept byte array smaller than size available")
	}
}

func TestBufferAcceptsOfferedBytesThatDoFillBufferWhenNotFull(t *testing.T) {
	b := internal.NewBuffer(100)
	if !b.Offer(make([]byte, 10)) {
		t.Error("Failed to partially fill the buffer")
	}
	if !b.Offer(make([]byte, 90)) {
		t.Error("Unfilled buffer did not accept byte array smaller than size available")
	}
}

func TestBufferDoesNotAcceptOfferedBytesThatWouldOverfillBuffer(t *testing.T) {
	b := internal.NewBuffer(100)
	if b.Offer(make([]byte, 101)) {
		t.Error("Buffer accepted more bytes that its size")
	}
	if !b.Offer(make([]byte, 10)) {
		t.Error("Failed to partially fill the buffer")
	}
	if b.Offer(make([]byte, 91)) {
		t.Error("Buffer accepted more bytes than the space it has remaining")
	}
}

func TestBufferPopsNothingWhenEmpty(t *testing.T) {
	b := internal.NewBuffer(2048)
	byts, err := b.Pop()
	if len(byts) != 0 {
		t.Errorf("%s was popped from buffer, expected nothing", byts)
	}
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
}

func TestBufferPopsAllBytesWhenItContainsUnder1024Bytes(t *testing.T) {
	b := internal.NewBuffer(2048)
	storedBytes := bytes.NewBufferString("Contents of the buffer").Bytes()
	b.Offer(storedBytes)
	poppedBytes, err := b.Pop()
	if !bytes.Equal(poppedBytes, storedBytes) {
		t.Errorf("%v was popped, but put in %v", poppedBytes, storedBytes)
	}
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
}

func TestBufferPopsOnly1024BytesWhenItContainsMore(t *testing.T) {
	b := internal.NewBuffer(2048)
	storedBytes := bytes.NewBufferString("")
	for i := 0; i < 64; i++ {
		sixteenBytes := "Somesixteenbytes"
		_, _ = storedBytes.WriteString(sixteenBytes)
	}
	b.Offer(storedBytes.Bytes())
	b.Offer(storedBytes.Bytes())
	poppedBytes, err := b.Pop()
	if !bytes.Equal(poppedBytes, storedBytes.Bytes()) {
		t.Errorf("%v was popped, but first 1024 put in %v", poppedBytes, storedBytes)
	}
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
}

func TestBytesPoppedAreRemovedFromBufferTest(t *testing.T) {
	b := internal.NewBuffer(2048)
	ignoredBytes := bytes.NewBufferString("")
	for i := 0; i < 64; i++ {
		sixteenBytes := "Somesixteenbytes"
		_, _ = ignoredBytes.WriteString(sixteenBytes)
	}
	b.Offer(ignoredBytes.Bytes())
	bytesOfInterest := bytes.NewBufferString("Should come out at end").Bytes()
	b.Offer(bytesOfInterest)
	_, err := b.Pop()
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
	poppedBytes, err := b.Pop()
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
	if !bytes.Equal(poppedBytes, bytesOfInterest) {
		t.Errorf("%v was popped, expected %v left in the buffer", poppedBytes, bytesOfInterest)
	}
	if err != nil {
		t.Errorf("%v error returned from pop, expected none", err)
	}
}
