package internal_test

import (
	"testing"

	"github.com/snasphysicist/go-copy/pkg/internal"
)

func TestMinReturnsValueWhenOnlyOneValueSupplied(t *testing.T) {
	value := uint64(102)
	minimum := internal.Minimum(value)
	if value != minimum {
		t.Errorf("%d returned as minimum, even though %d was the only value supplied", minimum, value)
	}
}

func TestMinReturnsLowestValueOfAllThoseSupplied(t *testing.T) {
	lowest := uint64(102)
	values := []uint64{8749, 90943587, lowest, 4895, 934895}
	minimum := internal.Minimum(values[0], values...)
	if lowest != minimum {
		t.Errorf("%d returned as minimum, even though %d was the lowest value of %v", minimum, lowest, values)
	}
}
