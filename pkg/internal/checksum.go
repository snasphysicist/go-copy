package internal

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"log"
)

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
func MD5Sum(r io.Reader, bufferSizeBytes int) string {
	h := md5.New()
	for {
		b := make([]byte, bufferSizeBytes)
		n, err := r.Read(b)
		if err == io.EOF {
			return hex.EncodeToString(appendToHash(h, b[:n], true))
		}
		if err != nil {
			panic(err)
		}
		appendToHash(h, b[:n], false)
	}
}
