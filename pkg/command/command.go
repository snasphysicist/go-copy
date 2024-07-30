package command

import (
	"flag"
	"time"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

// bufferSize is the default bufer size, around 100MB
const bufferSize = 100 * 1024 * 1024

func Copy() {
	arguments := parseFlags()
	from := arguments.from
	to := arguments.to

	s := internal.SizeOf(from)

	crossBuffer := internal.NewBuffer(bufferSize)

	shutdown := make(chan struct{})

	readerDone := make(chan struct{})
	writerDone := make(chan struct{})

	pr := internal.NewProgressReporter(s, shutdown)
	reader := internal.NewReader(from, &crossBuffer, readerDone, &pr, s)
	writer := internal.NewWriter(to, &crossBuffer, readerDone, &pr, s)

	go pr.Report(time.Now())
	go reader.Start()
	go writer.Start()

	<-readerDone
	<-writerDone

	close(shutdown)
	time.Sleep(10 * time.Millisecond)
}

type arguments struct {
	from string
	to   string
}

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
