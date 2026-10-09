//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue337 -i asyncapi.yaml -o ./asyncapi.gen.go
package issue337

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumValuesAreGenerated(t *testing.T) {
	assert.Equal(t, ColorSchema("red"), ColorSchemaRed)
	assert.Equal(t, ColorSchema("amber"), ColorSchemaAmber)
	assert.Equal(t, ColorSchema("green"), ColorSchemaGreen)

	color := ColorSchemaGreen
	foo := FooSchema{Color: &color}
	assert.Equal(t, ColorSchemaGreen, *foo.Color)
}

func TestInlineEnumValuesAreGenerated(t *testing.T) {
	assert.Equal(t, "not-available", string(InlinePropertyFromFooSchemaNotAvailable))

	foo := FooSchema{}
	foo.SetDefaults()
	assert.Equal(t, InlineDefaultPropertyFromFooSchema("low"), *foo.InlineDefault)
}
