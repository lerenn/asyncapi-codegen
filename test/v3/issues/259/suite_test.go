//go:generate go run ../../../../cmd/asyncapi-codegen -p issue259default -i ./asyncapi.yaml -o ./default/asyncapi.gen.go
//go:generate go run ../../../../cmd/asyncapi-codegen -p issue259forcepointers --force-pointers -i ./asyncapi.yaml -o ./forcepointers/asyncapi.gen.go

package issue259

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	issue259default "github.com/lerenn/asyncapi-codegen/test/v3/issues/259/default"
	issue259forcepointers "github.com/lerenn/asyncapi-codegen/test/v3/issues/259/forcepointers"
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

func (suite *Suite) TestDefault() {
	validate := validator.New()

	// Required absent - no error due to empty/nil fields
	assert.NoError(suite.T(), validate.Struct(
		issue259default.TestMessagePayload{
			NonReqArray: nil,
			NonReqField: nil,
			ReqArray:    []string{},
			ReqField:    "",
		},
	))

	var out issue259default.TestMessagePayload
	err := json.Unmarshal([]byte(`{}`), &out)
	assert.NoError(suite.T(), err)
}

func (suite *Suite) TestForcePointers() {
	validate := validator.New()

	cases := []struct {
		name   string
		inJSON string

		isErrExpected bool
	}{
		{
			name:          "empty input",
			inJSON:        `{}`,
			isErrExpected: true,
		},
		{
			name: "full input",
			inJSON: `{"reqField": "something","nonReqField": "something","reqArray": ["something"],
"nonReqArray": ["something"]}`,
			isErrExpected: false,
		},
		{
			name:          "only required fields",
			inJSON:        `{"reqField": "something","reqArray": ["something"]}`,
			isErrExpected: false,
		},
		{
			name:          "empty array",
			inJSON:        `{"reqField": "something","reqArray": []}`,
			isErrExpected: false,
		},
		{
			name:          "empty string",
			inJSON:        `{"reqField": "","reqArray": ["something"]}`,
			isErrExpected: false,
		},
		{
			name:          "missing string",
			inJSON:        `{"reqArray": ["something"]}`,
			isErrExpected: true,
		},
		{
			name:          "missing array",
			inJSON:        `{"reqField": "something"}`,
			isErrExpected: true,
		},
	}

	for _, tc := range cases {
		suite.T().Run(tc.name, func(t *testing.T) {
			var out issue259forcepointers.TestMessagePayload
			err := json.Unmarshal([]byte(tc.inJSON), &out)
			assert.NoError(t, err)

			if tc.isErrExpected {
				assert.Error(t, validate.Struct(out))
			} else {
				assert.NoError(t, validate.Struct(out))
			}
		})
	}
}
