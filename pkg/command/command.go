package command

import (
	"flag"

	"github.com/snasphysicist/go-copy/pkg/copy"
)

// bufferSizeBytes is the default bufer size, around 100MB
const bufferSizeBytes = 100 * 1024 * 1024

// syncEachBytes specifies after approximately
// how many written bytes we will try to
// force to flush to disk
const syncEachBytes = uint64(1000000)

// Copy implements the copy command, to copy a single source file to a single destination
func Copy() {
	arguments := parseFlags()
	from := arguments.from
	to := arguments.to
	copy.FileToFile(from, to, bufferSizeBytes, syncEachBytes)
}

// arguments contains the parsed and validated arguments to the Copy command
type arguments struct {
	from string
	to   string
}

// parseFlags extracts the flags/arguments for the Copy command
// panicing if anything is invalid or missing
func parseFlags() arguments {
	var a arguments
	flag.StringVar(&a.from, "from", "", "source file to be copied")
	flag.StringVar(&a.to, "to", "", "destination file to copy to")
	flag.Parse()
	if a.from == "" {
		panic("Must have from argument")
	}
	if a.to == "" {
		panic("Must have to argument")
	}
	return a
}
