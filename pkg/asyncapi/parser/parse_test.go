package parser

import (
	"fmt"
	"testing"

	"github.com/lerenn/asyncapi-codegen/pkg/asyncapi"
	asyncapiv3 "github.com/lerenn/asyncapi-codegen/pkg/asyncapi/v3"
	"github.com/stretchr/testify/suite"
)

func (suite *ParseSuite) TestOpenAPIDependencyKeepsOnlyComponents() {
	// An OpenAPI document with an array-typed "servers" (invalid in AsyncAPI)
	// must still parse, keeping only its components/schemas.
	data := []byte(`{
		"openapi": "3.0.0",
		"servers": [{"url": "https://example.com"}],
		"paths": {},
		"components": {"schemas": {"Pet": {"type": "object"}}}
	}`)

	spec, err := FromJSON(FromJSONParams{Data: data})
	suite.Require().NoError(err)

	v3, ok := spec.(*asyncapiv3.Specification)
	suite.Require().True(ok)
	suite.Require().Contains(v3.Components.Schemas, "Pet")
	suite.Require().Empty(v3.Servers)
}

func TestParseSuite(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}

type ParseSuite struct {
	suite.Suite
}

func (suite *ParseSuite) TestCorrectVersions() {
	correctVersions := []string{
		"2.0.0", "2.1.0", "2.2.0", "2.3.0", "2.4.0", "2.5.0", "2.6.0",
		"3.0.0", "3.1.0",
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
