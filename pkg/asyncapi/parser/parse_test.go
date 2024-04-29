package parser

import (
	"fmt"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	"github.com/stretchr/testify/suite"
)

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

type ParseSuite struct {
	suite.Suite
}

func (suite *ParseSuite) TestCorrectVersions() {
	correctVersions := []string{
		"2.0.0", "2.1.0", "2.2.0", "2.3.0", "2.4.0", "2.5.0", "2.6.0",
		"3.0.0",
	}

	suite.Require().Equal(len(correctVersions), len(asyncapi.SupportedVersions))

	for _, v := range correctVersions {
		b := []byte(fmt.Sprintf("{\"asyncapi\":\"%s\"}", v))
		_, err := FromJSON(FromJSONParams{
			Data: b,
		})
		suite.Require().NoError(err)
	}
}

func (suite *ParseSuite) TestIncorrectVersions() {
	correctVersions := []string{
		"1.0.0",
		"abc",
	}

	for _, v := range correctVersions {
		b := []byte(fmt.Sprintf("{\"asyncapi\":\"%s\"}", v))
		_, err := FromJSON(FromJSONParams{
			Data: b,
		})
		suite.Require().Error(err)
		suite.Require().ErrorIs(err, ErrInvalidVersion)
	}
}
