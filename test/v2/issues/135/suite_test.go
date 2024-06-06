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
	suite.Require().Equal("v2.issue135.group", V2Issue135GroupPath)
	suite.Require().Equal("v2.issue135.info", V2Issue135InfoPath)
	suite.Require().Equal("v2.issue135.project", V2Issue135ProjectPath)
	suite.Require().Equal("v2.issue135.resource", V2Issue135ResourcePath)
	suite.Require().Equal("v2.issue135.status", V2Issue135StatusPath)

	// Check path list
	for _, p := range ChannelsPaths {
		switch p {
		case V2Issue135GroupPath:
		case V2Issue135InfoPath:
		case V2Issue135ProjectPath:
		case V2Issue135ResourcePath:
		case V2Issue135StatusPath:
		default:
			suite.Require().Fail("unknown channel path", p)
		}
	}
}
