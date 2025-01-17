package internal

import "sync"

// buffer offers a synchronised byte buffer with
// and upper limit on its size
type buffer struct {
	b   []byte
	max uint64
	l   *sync.Mutex
}

// NewBuffer returns a new buffer with default size
func NewBuffer(size uint64) buffer {
	return buffer{b: make([]byte, 0), max: size, l: &sync.Mutex{}}
}

// Offer attempts to add the provided bytes to the buffer,
// returning true on success (there was room)
// and false on failure (there was not room)
func (b *buffer) Offer(byt []byte) bool {
	b.l.Lock()
	defer b.l.Unlock()
	if len(b.b)+len(byt) > int(b.max) {
		return false
	}
	b.b = append(b.b, byt...)
	return true
}

// Pop returns the first available up to 1024 bytes
// from the buffer, and removes them from the buffer
func (b *buffer) Pop() ([]byte, error) {
	b.l.Lock()
	defer b.l.Unlock()
	n := Minimum(1024, uint64(len(b.b)))
	toPop := b.b[:n]
	b.b = b.b[n:]
	return toPop, nil
}
