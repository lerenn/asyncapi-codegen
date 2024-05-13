//go:generate go run ../../../../cmd/asyncapi-codegen -p issue131 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue131

import (
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
		StringProp:           Ptr("test"),
		ArrayProp:            []string{"test1", "test2"},
		IntegerProp:          Ptr[int64](2),
		IntegerExclusiveProp: Ptr[int64](3),
		FloatProp:            Ptr[float64](2.55),
		EnumProp:             Ptr("amber"),
		ConstProp:            Ptr("Canada"),
	}
}

func (suite *Suite) TestValid() {
	err := validator.New().Struct(ValidTestSchema())

	assert.NoError(suite.T(), err)
}

func (suite *Suite) TestString() {
	tooSmall := ValidTestSchema()
	tooSmall.StringProp = Ptr("n")

	assert.Error(suite.T(), validator.New().Struct(tooSmall))

	tooLong := ValidTestSchema()
	tooLong.StringProp = Ptr("WayTooLong")

	assert.Error(suite.T(), validator.New().Struct(tooLong))
}

func (suite *Suite) TestInteger() {
	tooSmall := ValidTestSchema()
	tooSmall.IntegerProp = Ptr[int64](1)

	assert.Error(suite.T(), validator.New().Struct(tooSmall))

	tooLarge := ValidTestSchema()
	tooLarge.IntegerProp = Ptr[int64](6)

	assert.Error(suite.T(), validator.New().Struct(tooLarge))
}

func (suite *Suite) TestExclusiveInteger() {
	tooSmall := ValidTestSchema()
	tooSmall.IntegerProp = Ptr[int64](2)

	assert.NoError(suite.T(), validator.New().Struct(tooSmall))

	tooLarge := ValidTestSchema()
	tooLarge.IntegerProp = Ptr[int64](5)

	assert.NoError(suite.T(), validator.New().Struct(tooLarge))
}

func (suite *Suite) TestFloat() {
	tooSmall := ValidTestSchema()
	tooSmall.FloatProp = Ptr[float64](2.49)

	assert.Error(suite.T(), validator.New().Struct(tooSmall))

	tooLarge := ValidTestSchema()
	tooLarge.FloatProp = Ptr[float64](5.51)

	assert.Error(suite.T(), validator.New().Struct(tooLarge))
}

func (suite *Suite) TestRequired() {
	invalidAbsent := ValidTestSchema()
	invalidAbsent.RequiredProp = ""

	assert.Error(suite.T(), validator.New().Struct(invalidAbsent))
}

func (suite *Suite) TestArray() {
	empty := ValidTestSchema()
	empty.ArrayProp = []string{}

	assert.Error(suite.T(), validator.New().Struct(empty))

	tooManyElt := ValidTestSchema()
	tooManyElt.ArrayProp = []string{"1", "2", "3", "4", "5", "6"}

	assert.Error(suite.T(), validator.New().Struct(tooManyElt))

	tooFewElt := ValidTestSchema()
	tooFewElt.ArrayProp = []string{"1"}

	assert.Error(suite.T(), validator.New().Struct(tooFewElt))

	notUnique := ValidTestSchema()
	notUnique.ArrayProp = []string{"1", "1"}
	assert.Error(suite.T(), validator.New().Struct(notUnique))
}

func (suite *Suite) TestEnum() {
	wrong := ValidTestSchema()
	wrong.EnumProp = Ptr("Wrong")

	assert.Error(suite.T(), validator.New().Struct(wrong))
}

func (suite *Suite) TestConst() {
	wrong := ValidTestSchema()
	wrong.EnumProp = Ptr("Wrong")

	assert.Error(suite.T(), validator.New().Struct(wrong))
}
