package utils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestMapToListWithString(t *testing.T) {
	// Parameters
	input := map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
	}
	expectedOutput := []string{"a", "b", "c", "d"}

	// Test by sorting slice
	less := func(a, b string) bool { return a < b }
	assert.Equal(t, cmp.Diff(expectedOutput, MapToList(input), cmpopts.SortSlices(less)), "")
}
