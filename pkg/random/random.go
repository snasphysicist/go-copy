package random

import "crypto/rand"

// Bytes returns n random bytes
func Bytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
