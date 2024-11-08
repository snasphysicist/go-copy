package internal_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/snasphysicist/go-copy/pkg/internal"
	"github.com/snasphysicist/go-copy/pkg/random"
)

func TestMD5SumTestFileMatchesCommandLineMD5Sum(t *testing.T) {
	f, err := os.Open("checksumme.txt")
	if err != nil {
		t.Fatalf("Failed to open file %v", err)
	}
	sum, _ := internal.MD5Sum(f, 10)
	if sum != "7915fab42d254ffc3fbd14174217775f" {
		t.Errorf("Calculated md5sum %s, expected 7915fab42d254ffc3fbd14174217775f", sum)
	}
}

// mockByteSource is a byteSource for testing with controllable
// error behaviour and which can read from a wrapped reader
type mockByteSource struct {
	openError  error
	readData   io.Reader
	readError  error
	closeError error
}

// Open implements byteSource on mockByteSource,
// returns the configured openError
func (m mockByteSource) Open() error {
	return m.openError
}

// Read implements byteSource on mockByteSource,
// returns the readError if configured,
// else passes through to readData
func (m mockByteSource) Read(b []byte) (int, error) {
	if m.readError != nil {
		return 0, m.readError
	}
	return m.readData.Read(b)
}

// Close implements byteSource on mockByteSource,
// returns the configured closeError
func (m mockByteSource) Close() error {
	return m.closeError
}

func TestErrorWhenSourceCannotBeOpened(t *testing.T) {
	source := mockByteSource{openError: errors.New("some error")}
	target := mockByteSource{}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though source.Open errored")
	}
}

func TestErrorWhenSourceCannotBeRead(t *testing.T) {
	source := mockByteSource{readError: errors.New("some error")}
	target := mockByteSource{}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though source.Read errored")
	}
}

func TestErrorWhenSourceCannotBeClosed(t *testing.T) {
	source := mockByteSource{
		readData:   bytes.NewReader(random.Bytes(1024)),
		closeError: errors.New("some error"),
	}
	target := mockByteSource{}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though source.Close errored")
	}
}

func TestErrorWhenTargetCannotBeOpened(t *testing.T) {
	source := mockByteSource{readData: bytes.NewReader(random.Bytes(1024))}
	target := mockByteSource{openError: errors.New("some error")}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though target.Open errored")
	}
}

func TestErrorWhenTargetCannotBeRead(t *testing.T) {
	source := mockByteSource{readData: bytes.NewReader(random.Bytes(1024))}
	target := mockByteSource{readError: errors.New("some error")}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though target.Read errored")
	}
}

func TestErrorWhenTargetCannotBeClosed(t *testing.T) {
	source := mockByteSource{readData: bytes.NewReader(random.Bytes(1024))}
	target := mockByteSource{
		readData:   bytes.NewReader(random.Bytes(1024)),
		closeError: errors.New("some error"),
	}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err == nil {
		t.Error("No error returned even though target.Close errored")
	}
}

func TestNoErrorWhenSourceAndTargetContainSameBytes(t *testing.T) {
	b := random.Bytes(1024)
	source := mockByteSource{readData: bytes.NewReader(b)}
	target := mockByteSource{readData: bytes.NewReader(b)}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err != nil {
		t.Errorf("Error '%s' even though checksums should match", err)
	}
}

func TestErrorWhenSourceAndTargetDontContainSameBytes(t *testing.T) {
	b1 := random.Bytes(1024)
	b2 := append(b1[:1023], b1[1023]+1)
	if len(b1) != len(b2) {
		panic(fmt.Sprintf("len(b1) %d != len(b2) %d", len(b1), len(b2)))
	}
	source := mockByteSource{readData: bytes.NewReader(b1)}
	target := mockByteSource{readData: bytes.NewReader(b2)}
	v := internal.NewVerifier(source, target)
	err := v.Compare()
	if err != nil {
		t.Errorf("Error '%s' even though checksums should match", err)
	}
}
