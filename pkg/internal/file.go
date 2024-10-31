package internal

import (
	"fmt"
	"os"

	"github.com/snasphysicist/go-copy/pkg/panicing"
)

// FormatSize takes a size in bytes and
// returns a human readable representation
// like 165mb (with the number never exceeding 1024)
func FormatSize(s uint64) string {
	sf := float64(s)
	unitEvolution := map[string]string{
		"b":  "kb",
		"kb": "mb",
		"mb": "gb",
		"gb": "tb",
		"tb": "pb",
		"pb": "eb",
		"eb": "zb",
		"zb": "yb",
	}
	unit := "b"
	for {
		if sf < 1024.0 {
			return fmt.Sprintf("%.2f%s", sf, unit)
		}
		unit = unitEvolution[unit]
		assert(func() bool { return unit != "" }, fmt.Sprintf("no unit found at size %d", s))
		sf = sf / 1024.0
	}
}

// SizeOf returns the size of the file at given path
// in bytes as reported by the os
func SizeOf(path string) uint64 {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() { panicing.OnError(f.Close()) }()
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	return uint64(fi.Size())
}
