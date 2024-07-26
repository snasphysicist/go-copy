package main

import (
	"flag"
	"time"
)

func cp() {
	arguments := parseFlags()
	from := arguments.from
	to := arguments.to

	s := sizeOf(from)

	crossBuffer := newBuffer()

	shutdown := make(chan struct{})

	readerDone := make(chan struct{})
	writerDone := make(chan struct{})

	pr := ProgressReporter{
		read:       0,
		written:    0,
		toTransfer: s,
		shutdown:   shutdown,
	}
	reader := Reader{path: from, b: &crossBuffer, done: readerDone, pr: &pr, toTransfer: s}
	writer := Writer{path: to, b: &crossBuffer, done: writerDone, pr: &pr, toTransfer: s}

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
