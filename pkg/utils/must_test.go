package utils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleMust_withNoError() {
	// Following line will not panic as there is no error
	result := Must("abc", nil)

	// Should print "abc"
	fmt.Println(result)
}

func ExampleMust_withError() {
	// Following line will panic as there is an error
	Must("abc", errors.New("an error"))
}

func TestMust(t *testing.T) {
	// Check with successful
	s := "test"
	s2 := Must(s, nil)
	assert.Equal(t, s, s2)

	// Check with unsuccessful
	assert.Panics(t, func() {
		_ = Must(s, errors.New("error"))
	})
}
