//go:generate go run ../../../cmd/asyncapi-codegen -g types -p issue135 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue135

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

type Suite struct {
	suite.Suite
}

func (suite *Suite) TestCheckPaths() {
	// Check path contents
	suite.Require().Equal("group", GroupPath)
	suite.Require().Equal("info", InfoPath)
	suite.Require().Equal("project", ProjectPath)
	suite.Require().Equal("resource", ResourcePath)
	suite.Require().Equal("status", StatusPath)

	// Check path list
	for _, p := range ChannelsPaths {
		switch p {
		case GroupPath:
		case InfoPath:
		case ProjectPath:
		case ResourcePath:
		case StatusPath:
		default:
			suite.Require().Fail("unknown channel path", p)
		}
	}
}
