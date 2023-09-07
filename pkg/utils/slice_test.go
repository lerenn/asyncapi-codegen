package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(UtilsSuite))
}

type UtilsSuite struct {
	suite.Suite
}

func (suite *UtilsSuite) TestIsInSlice() {
	cases := []struct {
		Slice  []string
		Match  string
		Result bool
	}{
		// In slice
		{Slice: []string{"1", "2", "3"}, Match: "2", Result: true},
		// Not in slice
		{Slice: []string{"1", "2", "3"}, Match: "4", Result: false},
	}

	for i, c := range cases {
		suite.Require().Equal(c.Result, IsInSlice(c.Slice, c.Match), i)
	}
}
