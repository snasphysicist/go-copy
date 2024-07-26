package main

import (
	"sync/atomic"
	"time"
)

// ProgressReporter allows progress to be reported to it
// and reports progress to the user
// by printing it to the terminal.
// The application should report when it is being
// shut down by closing the shutdown channel.
type ProgressReporter struct {
	read       uint64
	written    uint64
	toTransfer uint64
	shutdown   <-chan struct{}
}

// ReportBytesRead tells the reporter that
// an additional n bytes has been read
func (pr *ProgressReporter) ReportBytesRead(n uint64) {
	atomic.AddUint64(&pr.read, n)
}

// ReportBytesRead tells the reporter that
// an additional n bytes has been written
func (pr *ProgressReporter) ReportBytesWritten(n uint64) {
	atomic.AddUint64(&pr.written, n)
}

// BytesRead returns the number of bytes reported to be read
func (pr *ProgressReporter) BytesRead() uint64 {
	return atomic.LoadUint64(&pr.read)
}

// BytesWritten returns the number of bytes reported to be written
func (pr *ProgressReporter) BytesWritten() uint64 {
	return atomic.LoadUint64(&pr.written)
}

// min returns the minimum value of all uint64s provided
func min(n1 uint64, n2 ...uint64) uint64 {
	m := n1
	for _, n := range n2 {
		if n < n1 {
			m = n
		}
	}
	return m
}

// Report prints the progress reporter to the reporter
// out to the terminal in an infinte loop,
// designed to be run in a goroutine from a command.
// Prints about once per second, or when the application
// is being shut down.
func (pr *ProgressReporter) Report(start time.Time) {
	eachSecond := time.Tick(time.Second)
	for {
		select {
		case <-pr.shutdown:
		case <-eachSecond:
		}
		elapsed := time.Since(start)
		bytesRead := pr.BytesRead()
		bytesWritten := pr.BytesWritten()
		rate := (float64(min(bytesRead, bytesWritten)) / float64(elapsed.Microseconds()/1000000))
		remaining := (float64(pr.toTransfer) - float64(min(bytesRead, bytesWritten))) / rate
		print("\r")
		print(
			"Read ", formatSize(bytesRead),
			" Written ", formatSize(bytesWritten),
			" Speed ", formatSize(uint64(rate)), "/s",
			" Elapsed ", elapsed.Round(1*time.Second).String(),
			" Remaining ", (time.Duration(remaining) * time.Second).String(),
			"             ",
		)
	}
}
