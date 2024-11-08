package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
)

// bufferSizeBytes is the default size of the buffer
// used when reading through the file to calculate the checksum
const bufferSizeBytes = 1024

// appendToHash append the complete content of b to h,
// if sum then return the md5 sum, else return empty bytes.
// panics on any kind of error.
func appendToHash(h hash.Hash, b []byte, sum bool) []byte {
	nOut, err := h.Write(b)
	if len(b) != nOut {
		log.Panicf("Somehow read %d but wrote %d bytes", len(b), nOut)
	}
	if err != nil {
		panic(err)
	}
	if sum {
		s := make([]byte, 0)
		return h.Sum(s)
	}
	return []byte{}
}

// MD5Sum calculates the MD5 sum of the content of r,
// reading into an internal buffer with given size in bytes,
// and returns the sum as a hex encoded string.
// panics on any error.
func MD5Sum(r io.Reader, bufferSizeBytes int) (string, error) {
	h := md5.New()
	for {
		b := make([]byte, bufferSizeBytes)
		n, err := r.Read(b)
		if err == io.EOF {
			return hex.EncodeToString(appendToHash(h, b[:n], true)), nil
		}
		if err != nil {
			return "", fmt.Errorf("failed to read for md5 checksum %w", err)
		}
		appendToHash(h, b[:n], false)
	}
}

// byteSource represents a source of bytes from
// which a checksum is calculated. The Open
// method is provided for compatibility
// with file-like objects, but can be a noop.
type byteSource interface {
	io.ReadCloser
	Open() error
}

// verifier allows to check the checksums
// of two sources of bytes are equal
type verifier struct {
	source byteSource
	target byteSource
}

// NewVerifier creates a new verifier, reading bytes in
// from the source and target byteSources
func NewVerifier(source byteSource, target byteSource) verifier {
	return verifier{source: source, target: target}
}

// Compare calculates the checksums of both source and target,
// returning an error if they differ (or if anything else goes wrong)
func (v verifier) Compare() error {
	err := v.source.Open()
	if err != nil {
		return fmt.Errorf("failed to open checksum source: %w", err)
	}
	sourceHash, err := MD5Sum(v.source, bufferSizeBytes)
	if err != nil {
		return fmt.Errorf("failed to read checksum source: %w", err)
	}
	err = v.source.Close()
	if err != nil {
		return fmt.Errorf("failed to close checksum source: %w", err)
	}
	err = v.target.Open()
	if err != nil {
		return fmt.Errorf("failed to open checksum target: %w", err)
	}
	targetHash, err := MD5Sum(v.target, bufferSizeBytes)
	if err != nil {
		return fmt.Errorf("failed to read checksum target: %w", err)
	}
	err = v.target.Close()
	if err != nil {
		return fmt.Errorf("failed to close checksum target: %w", err)
	}
	if sourceHash != targetHash {
		return fmt.Errorf(
			"checksum of read file %s "+
				"does not match checksum of written file %s",
			sourceHash,
			targetHash,
		)
	}
	return nil
}
