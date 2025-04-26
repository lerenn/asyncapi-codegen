//go:generate go run ../../../../cmd/asyncapi-codegen -p issue275 -g types -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue275

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

func (suite *Suite) TestMarshal_OmitEmpty() {
	var res TestSchema
	data, err := json.Marshal(res)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), "{\"withoutOmitEmpty\":null}", string(data))
}
