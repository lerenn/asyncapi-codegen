//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue135 -i ./asyncapi.yaml -o ./asyncapi.gen.go

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
	suite.Require().Equal("v3.issue135.group", GroupChannelPath)
	suite.Require().Equal("v3.issue135.info", InfoChannelPath)
	suite.Require().Equal("v3.issue135.project", ProjectChannelPath)
	suite.Require().Equal("v3.issue135.resource", ResourceChannelPath)
	suite.Require().Equal("v3.issue135.status", StatusChannelPath)

	// Check path list
	for _, p := range ChannelsPaths {
		switch p {
		case GroupChannelPath:
		case InfoChannelPath:
		case ProjectChannelPath:
		case ResourceChannelPath:
		case StatusChannelPath:
		default:
			suite.Require().Fail("unknown channel path", p)
		}
	}
}
