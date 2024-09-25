package copy_test

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/snasphysicist/go-copy/pkg/copy"
	"github.com/snasphysicist/go-copy/pkg/random"
)

// writeFile write content to the given path
// with some permissions (not important), panic on error
func writeFile(path string, content []byte) {
	err := os.WriteFile(path, content, 0666)
	if err != nil {
		panic(err)
	}
}

// randomFilePath returns a path to a file with random name in /tmp
func randomFilePath() string {
	return fmt.Sprintf("%s/%d", os.TempDir(), rand.Int63())
}

// deleteFile tries to delete the file at given path, panics on error
func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestCopyCopiesFileBetweenFromAndToPaths(t *testing.T) {
	content := random.Bytes(2780)
	from := randomFilePath()
	writeFile(from, content)
	defer deleteFile(from)
	to := randomFilePath()
	defer deleteFile(to)
	copy.FileToFile(from, to, 50, 250)
	written, err := os.ReadFile(to)
	if err != nil {
		t.Errorf("Failed to read target file with %v", err)
	}
	if len(content) != len(written) {
		t.Errorf("%d bytes were written, should be %d", len(written), len(content))
	}
	if !reflect.DeepEqual(content, written) {
		t.Error("Copied content did not match source content")
	}
}
