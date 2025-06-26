//go:generate go run ../../../../cmd/asyncapi-codegen -p issue283 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue283

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func (suite *Suite) ReadFile(filename string) (string, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(file), err
}

func NewSuite() *Suite {
	return &Suite{}
}

func (suite *Suite) TestAsyncApiGenGoContainsOnlyOneTestHeadersSchema() {
	// This test checks that the generated asyncapi.gen.go file contains only one
	content, err := suite.ReadFile("asyncapi.gen.go")
	suite.Require().NoError(err, "Failed to read asyncapi.gen.go file")
	suite.Require().Contains(content, "type TestHeadersSchema struct {", "asyncapi.gen.go should contain TestHeadersSchema schema")
}
