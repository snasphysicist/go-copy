package internal_test

import (
	"regexp"
	"testing"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

func withoutNumberPart(s string) string {
	re := regexp.MustCompile(`[\d\.]`)
	return re.ReplaceAllString(s, "")
}

func TestZeroBytesSuffixedAsBytes(t *testing.T) {
	s := internal.FormatSize(0)
	s = withoutNumberPart(s)
	if s != "b" {
		t.Errorf("0 bytes formatted as %s, expected b", s)
	}
}

func TestFewerThan1024BytesSuffixedAsBytes(t *testing.T) {
	s := internal.FormatSize(1023)
	s = withoutNumberPart(s)
	if s != "b" {
		t.Errorf("0 bytes formatted as %s, expected b", s)
	}
}

func Test1024BytesToSuffixedAsKilobytesBytes(t *testing.T) {
	s := internal.FormatSize(1024)
	s = withoutNumberPart(s)
	if s != "kb" {
		t.Errorf("1024 bytes formatted as %s, expected kb", s)
	}
}

func TestOneByteFewerThanMegabyteSuffixedAsKilobytes(t *testing.T) {
	byteCount := uint64((1024 * 1024) - 1)
	s := internal.FormatSize(byteCount)
	s = withoutNumberPart(s)
	if s != "kb" {
		t.Errorf("1 megabyte - 1 byte formatted as %s, expected kb", s)
	}
}

func TestExactlyOneMegabyteSuffixedAsMegabytes(t *testing.T) {
	byteCount := uint64(1024 * 1024)
	s := internal.FormatSize(byteCount)
	s = withoutNumberPart(s)
	if s != "mb" {
		t.Errorf("1 megabyte formatted as %s, expected mb", s)
	}
}

func TestOneByteFewerThanOneGigabyteSuffixedAsMegabytes(t *testing.T) {
	byteCount := uint64((1024 * 1024 * 1024) - 1)
	s := internal.FormatSize(byteCount)
	s = withoutNumberPart(s)
	if s != "mb" {
		t.Errorf("1 gigabyte - 1 byte formatted as %s, expected mb", s)
	}
}

func TestExactlyOneGigabyteSuffixedAsGigabytes(t *testing.T) {
	byteCount := uint64(1024 * 1024 * 1024)
	s := internal.FormatSize(byteCount)
	s = withoutNumberPart(s)
	if s != "gb" {
		t.Errorf("1 gigabyte formatted as %s, expected gb", s)
	}
}

func TestOneTerabyteMinusOneGigabyteSuffixedAsGigabytes(t *testing.T) {
	byteCount := uint64((1024 * 1024 * 1024 * 1024) - 1)
	s := internal.FormatSize(byteCount)
	s = withoutNumberPart(s)
	if s != "gb" {
		t.Errorf("1 terabyte - 1 gigabyte formatted as %s, expected gb", s)
	}
}

func TestAllSizesFormattedWithTwoDecimalPlacesTest(t *testing.T) {
	r := regexp.MustCompile(`\.\d\d[^\d]`)
	byteCounts := []uint64{
		1, 10, 100, 1024, 1025,
		2456, 10324, 128457,
		3 * 1024 * 1024, 79 * 1024 * 1024 * 2345,
	}
	for _, bc := range byteCounts {
		formatted := internal.FormatSize(bc)
		twoDecimalPlaces := r.FindAllString(formatted, -1)
		if len(twoDecimalPlaces) != 1 {
			t.Errorf("%d .dd strings found in %s, expected 1", len(twoDecimalPlaces), formatted)
		}
	}
}

func TestCorrectNumbersPrintedForEachByteValue(t *testing.T) {
	byteCounts := []uint64{
		1, 10, 100, 1024, 1333,
		2456, 10324, 128457,
		5678223, 38927583, 936573842,
		3295837586, 3295837586,
		29287536205,
	}
	allExpectedTwoDecimalPlaces := []string{
		"1.00", "10.00", "100.00", "1.00",
		"1.30", "2.40", "10.08", "125.45",
		"5.42", "37.12", "893.19", "3.07",
		"3.07", "27.28",
	}
	r := regexp.MustCompile(`[a-z]`)
	for i, bc := range byteCounts {
		formatted := internal.FormatSize(bc)
		twoDecimalPlaces := r.ReplaceAllString(formatted, "")
		expectedTwoDecimalPlaces := allExpectedTwoDecimalPlaces[i]
		if expectedTwoDecimalPlaces != twoDecimalPlaces {
			t.Errorf("Number printed as %s, expected %s", twoDecimalPlaces, expectedTwoDecimalPlaces)
		}
	}
}
