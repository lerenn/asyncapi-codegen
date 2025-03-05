//go:generate go run ../../../../cmd/asyncapi-codegen -p issue267 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue267

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
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

func (suite *Suite) TestValidate_Valid() {
	var res TestSchema
	err := json.Unmarshal([]byte(`{
		"EnumProp": "has a space"
	}`), &res)
	assert.NoError(suite.T(), err)
	err = validator.New().Struct(res)
	assert.NoError(suite.T(), err)
}

func (suite *Suite) TestValidate_Invalid() {
	var res TestSchema
	err := json.Unmarshal([]byte(`{
		"EnumProp": "nospace"
	}`), &res)
	assert.NoError(suite.T(), err)
	err = validator.New().Struct(res)
	assert.Error(suite.T(), err)
}
