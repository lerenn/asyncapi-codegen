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
		RequiredProp:         "test",
		ArrayProp:            []string{"test1", "test2"},
		IntegerProp:          Ptr[int64](2),
		IntegerExclusiveProp: Ptr[int64](3),
		FloatProp:            Ptr[float64](2.55),
		EnumProp:             Ptr("amber"),
		ConstProp:            Ptr("Canada"),
	}
}

func (suite *Suite) TestOmitEmpty() {
	testData := ValidTestSchema()
	err := validator.New().Struct(testData)

	assert.NoError(suite.T(), err)

	js, err := json.Marshal(testData)
	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"RequiredProp":"test","ArrayProp":["test1","test2"],"IntegerProp":2,"IntegerExclusiveProp":3,"FloatProp":2.55,"EnumProp":"amber","ConstProp":"Canada"}`, string(js))

}
