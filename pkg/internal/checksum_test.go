package internal_test

import (
	"os"
	"testing"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

func TestMD5SumTestFileMatchesCommandLineMD5Sum(t *testing.T) {
	f, err := os.Open("checksumme.txt")
	if err != nil {
		t.Fatalf("Failed to open file %v", err)
	}
	sum := internal.MD5Sum(f, 10)
	if sum != "7915fab42d254ffc3fbd14174217775f" {
		t.Errorf("Calculated md5sum %s, expected 7915fab42d254ffc3fbd14174217775f", sum)
	}
}
