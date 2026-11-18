//go:generate go run ../../../../cmd/asyncapi-codegen -g types -p issue344 -i asyncapi.yaml -o ./asyncapi.gen.go
package issue344

import (
	"encoding/json"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableFieldsArePointers(t *testing.T) {
	var u UserSchema
	require.NoError(t, json.Unmarshal([]byte(`{"apiKey":null,"name":"a","age":null}`), &u))
	assert.Nil(t, u.ApiKey)
	assert.Nil(t, u.Age)

	out, err := json.Marshal(u)
	require.NoError(t, err)
	assert.JSONEq(t, `{"apiKey":null,"name":"a"}`, string(out))
}

func TestRequiredNullableFieldAcceptsNull(t *testing.T) {
	validate := validator.New()
	assert.NoError(t, validate.Struct(UserSchema{Name: "a"}))

	short := "abc"
	assert.Error(t, validate.Struct(UserSchema{ApiKey: &short, Name: "a"}))
}
