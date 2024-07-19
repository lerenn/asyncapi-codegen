//go:generate go run ../../../../cmd/asyncapi-codegen -p issue245 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue245

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
func Ptr[T any](v T) *T {
	return &v
}

func ValidTestSchema() TestSchema {
	return TestSchema{
		RequiredProp: "test",
		ArrayProp:    []string{"test1", "test2"},
		IntegerProp:  Ptr[int64](2),
		FloatProp:    Ptr[float64](2.55),
		EnumProp:     Ptr("amber"),
		ConstProp:    Ptr("Canada"),
	}
}

func (suite *Suite) TestGenerateJsonOmitEmptyTag() {
	testData := ValidTestSchema()
	err := validator.New().Struct(testData)

	assert.NoError(suite.T(), err)

	testData.IntegerProp = nil
	testTable := []struct {
		name     string
		data     TestSchema
		expected string
	}{
		{
			name:     "ArrayProp is not nil",
			data:     TestSchema{RequiredProp: "test", ArrayProp: []string{"test1", "test2"}},
			expected: `{"RequiredProp":"test", "ArrayProp":["test1", "test2"]}`,
		},
		{
			name:     "IntegerProp is not nil",
			data:     TestSchema{RequiredProp: "test", IntegerProp: Ptr[int64](2)},
			expected: `{"RequiredProp":"test", "IntegerProp":2}`,
		},
		{
			name:     "FloatProp is not nil",
			data:     TestSchema{RequiredProp: "test", FloatProp: Ptr[float64](2.66)},
			expected: `{"RequiredProp":"test", "FloatProp":2.66}`,
		},
		{
			name:     "EnumProp is not nil",
			data:     TestSchema{RequiredProp: "test", EnumProp: Ptr("amber")},
			expected: `{"RequiredProp":"test", "EnumProp":"amber"}`,
		},
		{
			name:     "IntegerProp is not nil",
			data:     TestSchema{RequiredProp: "test", ConstProp: Ptr("Canada")},
			expected: `{"RequiredProp":"test", "ConstProp":"Canada"}`,
		},
	}
	for _, tt := range testTable {
		tt := tt
		suite.T().Run(tt.name, func(t *testing.T) {
			js, err := json.Marshal(tt.data)
			assert.NoError(suite.T(), err)
			assert.JSONEq(suite.T(), tt.expected, string(js))
		})
	}
}
