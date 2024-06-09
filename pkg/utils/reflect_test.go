package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldValueExists(t *testing.T) {
	// Parameters
	type testStruct struct {
		Field1 string
		Field2 int
		Field3 bool
	}
	test := testStruct{
		Field1: "test",
		Field2: 1,
		Field3: true,
	}

	// Test with field that exists and correct values
	assert.True(t, FieldValueExists(test, "Field1", "test"))
	assert.True(t, FieldValueExists(test, "Field2", 1))
	assert.True(t, FieldValueExists(test, "Field3", true))

	// Test with field that exists and incorrect values
	assert.False(t, FieldValueExists(test, "Field1", "test2"))
	assert.False(t, FieldValueExists(test, "Field2", 2))
	assert.False(t, FieldValueExists(test, "Field3", false))

	// Test with field that does not exist
	assert.False(t, FieldValueExists(test, "Field4", "test"))
}
