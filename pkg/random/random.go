package random

import (
	"crypto/rand"

	"github.com/snasphysicist/go-copy/pkg/panicing"
)

// Bytes returns n random bytes
func Bytes(n int) []byte {
	b := make([]byte, n)
	panicing.OnWriteError(rand.Read(b))
	return b
}
