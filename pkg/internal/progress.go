package internal

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

func NewProgressReporter(toTransfer uint64, shutdown <-chan struct{}) ProgressReporter {
	return ProgressReporter{read: 0, written: 0, toTransfer: toTransfer, shutdown: shutdown}
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

// printProgress is called to print the current progress, expecting
// start to be the start time of the transfer, suffix will be
// added to the end of the printed line (intended to allow a newline
// to printed on the final output, so that final output won't be
// overwritten by something on the command line)
func (pr *ProgressReporter) printProgress(start time.Time, suffix string) {
	elapsed := time.Since(start)
	bytesRead := pr.BytesRead()
	bytesWritten := pr.BytesWritten()
	rate := (float64(Minimum(bytesRead, bytesWritten)) / float64(elapsed.Microseconds()/1000000))
	remaining := (float64(pr.toTransfer) - float64(Minimum(bytesRead, bytesWritten))) / rate
	print("\r")
	print(
		"Read ", FormatSize(bytesRead),
		" Written ", FormatSize(bytesWritten),
		" Speed ", FormatSize(uint64(rate)), "/s",
		" Elapsed ", elapsed.Round(1*time.Second).String(),
		" Remaining ", (time.Duration(remaining) * time.Second).String(),
		"             ", suffix,
	)
}

// Report prints the progress reporter to the reporter
// out to the terminal in an infinte loop,
// designed to be run in a goroutine from a command.
// Prints about once per second, or when the application
// is being shut down.
func (pr *ProgressReporter) Report(start time.Time) {
	eachSecond := time.NewTicker(time.Second)
	for {
		select {
		case <-eachSecond.C:
			pr.printProgress(start, "")
		case <-pr.shutdown:
			pr.printProgress(start, "\n")
			return
		}
	}
}
